# LabelNamer Cache Design

## Problem

`LabelNamer.Build()` is called for every label name translation. Even when the label name is already compliant (e.g., `http_method`), the method still:
1. Calls `sanitizeLabelName()` which creates a `strings.Builder`
2. Iterates over every character to validate and replace invalid chars
3. Returns a new string

This creates unnecessary allocations for the common case where labels are already valid.

## Solution

Add an optional `*StringCache` to `LabelNamer` that caches label name transformations. Cache is enabled by default, with an option to disable.

Design reference: `FastStringTransformer` in VictoriaMetrics (`lib/bytesutil/fast_string_transformer.go`).

## Design

### 1. StringCache Structure

```go
type StringCache struct {
    m                sync.Map
    lastCleanupTime  atomic.Uint64
    expireDuration   time.Duration
}

type cacheEntry struct {
    lastAccessTime atomic.Uint64
    value          string
}
```

- `sync.Map` stores cached transformations
- `lastCleanupTime` controls cleanup frequency (lazy cleanup)
- `expireDuration` is the TTL for cache entries (default: 6 minutes)

### 2. LabelNamer Structure

```go
type LabelNamer struct {
    UTF8Allowed                    bool
    UnderscoreLabelSanitization    bool
    PreserveMultipleUnderscores    bool

    cache       *StringCache  // nil means cache disabled
}
```

### 3. Options Pattern

```go
type Option func(*LabelNamer)

func WithCacheDisabled() Option
func WithCacheExpireDuration(d time.Duration) Option
```

### 4. Build Flow

```
1. Check empty label → error
2. Check UTF8Allowed → special handling
3. Try cache.Load(label) → hit: return cached value
4. Miss: call buildWithoutCache(label) → get result
5. label = strings.Clone(label)           // Safe key memory
6. if result == label: result = label    // Point to safe memory
7. cache.Store(label, result)
8. if needCleanup(): delete expired entries
9. return result, nil
```

### 5. Memory Safety

Reference: VictoriaMetrics/#3227

When `sTransformed == s` (no transformation needed), the result may point to unsafe memory (temporary byte slice). Solution:
- Always `strings.Clone(label)` for the cache key
- If `result == label`, make result point to the cloned safe memory

### 6. Lazy Cleanup

Triggered on Transform call, not in background goroutine.

```go
func needCleanup(lastCleanupTime *atomic.Uint64, currentTime uint64) bool {
    lct := lastCleanupTime.Load()
    if lct+61 >= currentTime {
        return false  // Less than 61 seconds since last cleanup
    }
    return lastCleanupTime.CompareAndSwap(lct, currentTime)  // Only one goroutine cleans
}
```

- Cleanup frequency: max once per 61 seconds
- Entry TTL: configurable, default 6 minutes
- Access time: updated every 10 seconds (not every access)

### 7. Backward Compatibility

- Default `LabelNamer{}` has cache enabled
- `WithCacheDisabled()` returns a new namer with cache disabled
- Existing code behavior unchanged

## Files to Modify

- `label_namer.go` — add `StringCache`, update `LabelNamer`, add options
- `label_namer_test.go` — add cache tests
- `label_namer_bench_test.go` — add cache benchmark

## Testing Considerations

- Cache hit/miss scenarios
- Memory safety when result == label
- Cleanup triggered correctly
- Concurrent access safety
- Benchmark with/without cache
