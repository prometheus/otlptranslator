# LabelNamer Cache Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add optional StringCache to LabelNamer to cache label name transformations, avoiding allocations for already-compliant labels.

**Architecture:** Add a new `StringCache` struct using `sync.Map`, integrated into `LabelNamer.Build()` via options pattern. Cache is enabled by default, lazy cleanup on transform calls.

**Tech Stack:** Go 1.21+, `sync.Map`, `atomic`, `time`

---

## Task 1: Add StringCache struct and helper functions

**Files:**
- Modify: `label_namer.go:1-101`

**Step 1: Write the failing test first**

```go
// Add to label_namer_test.go
func TestStringCacheBasics(t *testing.T) {
    cache := NewStringCache()
    namer := &LabelNamer{cache: cache}

    // First call should miss cache
    result, err := namer.Build("http.method")
    require.NoError(t, err)
    require.Equal(t, "http_method", result)

    // Second call should hit cache
    result2, err := namer.Build("http.method")
    require.NoError(t, err)
    require.Equal(t, "http_method", result2)
}
```

**Step 2: Run test to verify it fails**

Run: `go test -run TestStringCacheBasics -v ./...`
Expected: FAIL - "undefined: NewStringCache"

**Step 3: Add StringCache struct and needCleanup helper**

Add after the `hasUnderscoresOnly` function in label_namer.go:

```go
// StringCache caches string transformations.
// It is safe for concurrent use.
type StringCache struct {
    m               sync.Map
    lastCleanupTime atomic.Uint64
    expireDuration  time.Duration
}

type cacheEntry struct {
    lastAccessTime atomic.Uint64
    value          string
}

// NewStringCache creates a new StringCache with default expiry duration.
func NewStringCache() *StringCache {
    return &StringCache{
        expireDuration: 6 * time.Minute,
    }
}

// needCleanup returns true if cleanup should be performed.
// It is called lazily on Transform to avoid background goroutines.
func needCleanup(lastCleanupTime *atomic.Uint64, currentTime uint64) bool {
    lct := lastCleanupTime.Load()
    if lct+61 >= currentTime {
        return false
    }
    return lastCleanupTime.CompareAndSwap(lct, currentTime)
}
```

**Step 4: Run test to verify it compiles (not the test itself yet)**

Run: `go build ./...`
Expected: PASS

**Step 5: Commit**

```bash
git add label_namer.go
git commit -m "feat(label_namer): add StringCache struct and needCleanup helper"
```

---

## Task 2: Add cache field and options to LabelNamer

**Files:**
- Modify: `label_namer.go:29-91`

**Step 1: Modify LabelNamer struct**

```go
type LabelNamer struct {
    UTF8Allowed bool
    // UnderscoreLabelSanitization, if true, enabled prepending 'key' to labels
    // starting with '_'. Reserved labels starting with `__` are not modified.
    //
    // Deprecated: This will be removed in a future version of otlptranslator.
    UnderscoreLabelSanitization bool
    // PreserveMultipleUnderscores enables preserving of multiple
    // consecutive underscores in label names when UTF8Allowed is false.
    // This option is discouraged as it violates the OpenTelemetry to Prometheus
    // specification, but may be needed for compatibility with legacy systems.
    PreserveMultipleUnderscores bool

    cache *StringCache  // nil means cache disabled
}
```

**Step 2: Add Option type and option functions**

```go
type Option func(*LabelNamer)

// WithCacheDisabled disables the transformation cache.
func WithCacheDisabled() Option {
    return func(ln *LabelNamer) {
        ln.cache = nil
    }
}

// WithCacheExpireDuration sets the cache entry TTL.
func WithCacheExpireDuration(d time.Duration) Option {
    return func(ln *LabelNamer) {
        if ln.cache == nil {
            ln.cache = NewStringCache()
        }
        ln.cache.expireDuration = d
    }
}
```

**Step 3: Add default cache initialization in New or document the change**

The default behavior should have cache enabled. Update `NewLabelNamer` or document that cache is enabled by default when `cache` field is nil and user can disable with `WithCacheDisabled()`.

Actually, since LabelNamer is used directly (not via constructor), we should initialize cache lazily or document that default is enabled.

For simplicity, initialize cache lazily in Build() if not set but options call is not made.

**Step 4: Run tests**

Run: `go test -v ./...`
Expected: PASS

**Step 5: Commit**

```bash
git add label_namer.go
git commit -m "feat(label_namer): add cache field and options to LabelNamer"
```

---

## Task 3: Implement Build with cache integration

**Files:**
- Modify: `label_namer.go:65-91`

**Step 1: Refactor Build to separate cache logic**

Rename current Build logic to `buildWithoutCache`:

```go
func (ln *LabelNamer) buildWithoutCache(label string) (string, error) {
    if len(label) == 0 {
        return "", errors.New("label name is empty")
    }

    if ln.UTF8Allowed {
        if hasUnderscoresOnly(label) {
            return "", fmt.Errorf("label name %q contains only underscores", label)
        }
        return label, nil
    }

    normalizedName := sanitizeLabelName(label, ln.PreserveMultipleUnderscores)

    if unicode.IsDigit(rune(normalizedName[0])) {
        normalizedName = "key_" + normalizedName
    } else if ln.UnderscoreLabelSanitization && strings.HasPrefix(normalizedName, "_") && !strings.HasPrefix(normalizedName, "__") {
        normalizedName = "key" + normalizedName
    }

    if hasUnderscoresOnly(normalizedName) {
        return "", fmt.Errorf("normalization for label name %q resulted in invalid name %q", label, normalizedName)
    }

    return normalizedName, nil
}
```

**Step 2: Update Build to use cache**

```go
func (ln *LabelNamer) Build(label string) (string, error) {
    if len(label) == 0 {
        return "", errors.New("label name is empty")
    }

    if ln.UTF8Allowed {
        if hasUnderscoresOnly(label) {
            return "", fmt.Errorf("label name %q contains only underscores", label)
        }
        return label, nil
    }

    // Try cache first
    if ln.cache != nil {
        if v, ok := ln.cache.m.Load(label); ok {
            e := v.(*cacheEntry)
            ct := fasttime.UnixTimestamp()
            if e.lastAccessTime.Load()+10 < ct {
                e.lastAccessTime.Store(ct)
            }
            return e.value, nil
        }
    }

    // Cache miss - transform
    result, err := ln.buildWithoutCache(label)
    if err != nil {
        return result, err
    }

    // Store in cache with memory safety
    if ln.cache != nil {
        label = strings.Clone(label)
        if result == label {
            result = label
        }
        e := &cacheEntry{
            value: result,
        }
        e.lastAccessTime.Store(fasttime.UnixTimestamp())
        ln.cache.m.Store(label, e)

        // Lazy cleanup
        ct := fasttime.UnixTimestamp()
        if needCleanup(&ln.cache.lastCleanupTime, ct) {
            deadline := ct - uint64(ln.cache.expireDuration.Seconds())
            ln.cache.m.Range(func(k, v any) bool {
                e := v.(*cacheEntry)
                if e.lastAccessTime.Load() < deadline {
                    ln.cache.m.Delete(k)
                }
                return true
            })
        }
    }

    return result, nil
}
```

**Step 3: Add fasttime import**

```go
import (
    "errors"
    "fmt"
    "strings"
    "sync"
    "sync/atomic"
    "time"
    "unicode"

    "github.com/prometheus/prometheus/storage/remote/otlptranslator/internal/fasttime"
)
```

Note: You may need to check the actual import path for fasttime in this project.

**Step 4: Run tests**

Run: `go test -v ./...`
Expected: PASS

**Step 5: Commit**

```bash
git add label_namer.go
git commit -m "feat(label_namer): integrate cache into Build method"
```

---

## Task 4: Add comprehensive cache tests

**Files:**
- Modify: `label_namer_test.go`

**Step 1: Write cache hit test**

```go
func TestLabelNamerCacheHit(t *testing.T) {
    cache := NewStringCache()
    namer := &LabelNamer{cache: cache}

    result1, err := namer.Build("http.method")
    require.NoError(t, err)
    require.Equal(t, "http_method", result1)

    // Same label should hit cache
    result2, err := namer.Build("http.method")
    require.NoError(t, err)
    require.Equal(t, "http_method", result2)
}
```

**Step 2: Write cache disabled test**

```go
func TestLabelNamerCacheDisabled(t *testing.T) {
    namer := &LabelNamer{}
    namer.cache = nil  // explicitly disabled

    result, err := namer.Build("http.method")
    require.NoError(t, err)
    require.Equal(t, "http_method", result)

    result2, err := namer.Build("http.method")
    require.NoError(t, err)
    require.Equal(t, "http_method", result2)
}
```

**Step 3: Write memory safety test (result == label)**

```go
func TestLabelNamerCacheMemorySafety(t *testing.T) {
    // Create a label that doesn't need transformation
    label := "already_valid_label"
    cache := NewStringCache()
    namer := &LabelNamer{cache: cache}

    result, err := namer.Build(label)
    require.NoError(t, err)
    require.Equal(t, label, result)

    // Result should be safe to use even if original label goes out of scope
    // This is tested by the Clone in the implementation
}
```

**Step 4: Write WithCacheDisabled option test**

```go
func TestLabelNamerWithCacheDisabled(t *testing.T) {
    namer := &LabelNamer{}
    WithCacheDisabled()(namer)

    result, err := namer.Build("http.method")
    require.NoError(t, err)
    require.Equal(t, "http_method", result)
}
```

**Step 5: Run tests**

Run: `go test -run TestLabelNamer -v ./...`
Expected: PASS

**Step 6: Commit**

```bash
git add label_namer_test.go
git commit -m "test(label_namer): add cache tests"
```

---

## Task 5: Add benchmarks for cache performance

**Files:**
- Modify: `label_namer_bench_test.go`

**Step 1: Add benchmark for cache hit scenario**

```go
func BenchmarkNormalizeLabelWithCache(b *testing.B) {
    labelNamer := LabelNamer{cache: NewStringCache()}
    // Pre-populate cache
    for _, input := range labelBenchmarkInputs {
        labelNamer.Build(input.label)
    }

    b.ResetTimer()
    for _, input := range labelBenchmarkInputs {
        b.Run(input.name, func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                labelNamer.Build(input.label)
            }
        })
    }
}
```

**Step 2: Run benchmarks**

Run: `go test -bench=BenchmarkNormalizeLabel -benchmem ./...`
Expected: Cache version should show fewer allocations

**Step 3: Commit**

```bash
git add label_namer_bench_test.go
git commit -m "bench(label_namer): add cache benchmark"
```

---

## Verification

After all tasks complete, run:
```bash
go test -v ./...
go test -bench=BenchmarkNormalizeLabel -benchmem ./...
```

Expected results:
- All tests pass
- Cache hit benchmark shows significant reduction in allocations compared to no-cache
