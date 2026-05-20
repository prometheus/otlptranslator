// Copyright 2025 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// Provenance-includes-location: https://github.com/prometheus/prometheus/blob/93e991ef7ed19cc997a9360c8016cac3767b8057/storage/remote/otlptranslator/prometheus/normalize_label_test.go
// Provenance-includes-license: Apache-2.0
// Provenance-includes-copyright: Copyright The Prometheus Authors

package otlptranslator

import (
	"errors"
	"testing"
)

var labelTestCases = []struct {
	label                                                      string
	sanitized                                                  string
	sanitizedMultipleUnderscores                               string
	sanitizedUnderscoreLabelSanitization                       string
	sanitizedMultipleUnderscoresAndUnderscoreLabelSanitization string
	wantSanitizationError                                      error
	wantSanitizationMultipleUnderscoresError                   error
	wantUTF8Error                                              error
}{
	{label: "label:with:colons", sanitized: "label_with_colons"},
	{label: "LabelWithCapitalLetters", sanitized: "LabelWithCapitalLetters"},
	{label: "label!with&special$chars)", sanitized: "label_with_special_chars_"},
	{
		label:                        "label_with_foreign_characters_字符",
		sanitized:                    "label_with_foreign_characters_",
		sanitizedMultipleUnderscores: "label_with_foreign_characters___",
	},
	{label: "label.with.dots", sanitized: "label_with_dots"},
	{
		label:                                ".foo",
		sanitized:                            "_foo",
		sanitizedUnderscoreLabelSanitization: "key_foo",
	},
	{label: "123label", sanitized: "key_123label"},
	{
		label:                                "_label_starting_with_underscore",
		sanitized:                            "_label_starting_with_underscore",
		sanitizedUnderscoreLabelSanitization: "key_label_starting_with_underscore",
	},
	{
		label:                                "__label_starting_with_2underscores",
		sanitized:                            "_label_starting_with_2underscores",
		sanitizedMultipleUnderscores:         "__label_starting_with_2underscores",
		sanitizedUnderscoreLabelSanitization: "key_label_starting_with_2underscores",
		sanitizedMultipleUnderscoresAndUnderscoreLabelSanitization: "__label_starting_with_2underscores",
	},
	{
		label:                                    "ようこそ",
		sanitized:                                "",
		wantSanitizationError:                    errors.New(`normalization for label name "ようこそ" resulted in invalid name "_"`),
		wantSanitizationMultipleUnderscoresError: errors.New(`normalization for label name "ようこそ" resulted in invalid name "____"`),
	},
	{
		label:                        "label__with__double__underscores",
		sanitized:                    "label_with_double_underscores",
		sanitizedMultipleUnderscores: "label__with__double__underscores",
	},
	{
		label:                        "label.name__with&&special##chars",
		sanitized:                    "label_name_with_special_chars",
		sanitizedMultipleUnderscores: "label_name__with__special__chars",
	},
	{
		label:                        "__reserved__label__name__",
		sanitized:                    "__reserved_label_name__",
		sanitizedMultipleUnderscores: "__reserved__label__name__",
	},
	{
		label:                        "trailing_underscores___",
		sanitized:                    "trailing_underscores_",
		sanitizedMultipleUnderscores: "trailing_underscores___",
	},
	{
		label:                 "",
		wantSanitizationError: errors.New("label name is empty"),
		wantUTF8Error:         errors.New("label name is empty"),
	},
	{
		label:                                    "__",
		wantSanitizationError:                    errors.New(`normalization for label name "__" resulted in invalid name "_"`),
		wantSanitizationMultipleUnderscoresError: errors.New(`normalization for label name "__" resulted in invalid name "__"`),
		wantUTF8Error:                            errors.New(`label name "__" contains only underscores`),
	},
}

func TestBuildLabel(t *testing.T) {
	for _, tt := range labelTestCases {
		t.Run(tt.label, func(t *testing.T) {
			t.Run("Not preserving multiple underscores", func(t *testing.T) {
				labelNamer := LabelNamer{}
				got, err := labelNamer.Build(tt.label)
				if tt.wantSanitizationError != nil {
					if err == nil {
						t.Fatalf("LabelNamer.Build(%q) should have returned an error %q", tt.label, tt.wantSanitizationError)
					}
					if err.Error() != tt.wantSanitizationError.Error() {
						t.Fatalf("LabelNamer.Build(%q) should have returned an error %q, but returned: %q", tt.label, tt.wantSanitizationError, err)
					}
					return
				}
				if err != nil {
					t.Fatalf("LabelNamer.Build(%q) returned an error: %s", tt.label, err)
				}

				if got != tt.sanitized {
					t.Errorf("LabelNamer.Build(%q) = %q, want %q", tt.label, got, tt.sanitized)
				}

				t.Run("Sanitizing leading underscore", func(t *testing.T) {
					labelNamer := LabelNamer{UnderscoreLabelSanitization: true}
					got, err := labelNamer.Build(tt.label)
					if err != nil {
						t.Fatalf("LabelNamer.Build(%q) returned an error: %s", tt.label, err)
					}

					want := tt.sanitized
					if tt.sanitizedUnderscoreLabelSanitization != "" {
						want = tt.sanitizedUnderscoreLabelSanitization
					}
					if got != want {
						t.Errorf("LabelNamer.Build(%q) = %q, want %q", tt.label, got, want)
					}
				})
			})
			t.Run("Preserving multiple underscores", func(t *testing.T) {
				labelNamer := LabelNamer{PreserveMultipleUnderscores: true}
				got, err := labelNamer.Build(tt.label)
				var wantError error
				switch {
				case tt.wantSanitizationMultipleUnderscoresError != nil:
					wantError = tt.wantSanitizationMultipleUnderscoresError
				case tt.wantSanitizationError != nil:
					wantError = tt.wantSanitizationError
				}
				if wantError != nil {
					if err == nil {
						t.Fatalf("LabelNamer.Build(%q) should have returned an error %q", tt.label, wantError)
					}
					if err.Error() != wantError.Error() {
						t.Fatalf("LabelNamer.Build(%q) should have returned an error %q, but returned: %q", tt.label, wantError, err)
					}
					return
				}
				if err != nil {
					t.Fatalf("LabelNamer.Build(%q) returned an error: %s", tt.label, err)
				}

				want := tt.sanitized
				if tt.sanitizedMultipleUnderscores != "" {
					want = tt.sanitizedMultipleUnderscores
				}
				if got != want {
					t.Errorf("LabelNamer.Build(%q) = %q, want %q", tt.label, got, want)
				}

				t.Run("Sanitizing leading underscore", func(t *testing.T) {
					labelNamer := LabelNamer{
						PreserveMultipleUnderscores: true,
						UnderscoreLabelSanitization: true,
					}
					got, err := labelNamer.Build(tt.label)
					if err != nil {
						t.Fatalf("LabelNamer.Build(%q) returned an error: %s", tt.label, err)
					}

					var want string
					switch {
					case tt.sanitizedMultipleUnderscoresAndUnderscoreLabelSanitization != "":
						want = tt.sanitizedMultipleUnderscoresAndUnderscoreLabelSanitization
					case tt.sanitizedUnderscoreLabelSanitization != "":
						want = tt.sanitizedUnderscoreLabelSanitization
					case tt.sanitizedMultipleUnderscores != "":
						want = tt.sanitizedMultipleUnderscores
					default:
						want = tt.sanitized
					}
					if got != want {
						t.Errorf("LabelNamer.Build(%q) = %q, want %q", tt.label, got, want)
					}
				})
			})
		})
	}
}

func TestBuildLabel_UTF8Allowed(t *testing.T) {
	for _, tt := range labelTestCases {
		t.Run(tt.label, func(t *testing.T) {
			labelNamer := LabelNamer{UTF8Allowed: true}
			got, err := labelNamer.Build(tt.label)
			if tt.wantUTF8Error != nil {
				if err == nil {
					t.Fatalf("LabelNamer.Build(%q) should have returned an error %q", tt.label, tt.wantUTF8Error)
				}
				if err.Error() != tt.wantUTF8Error.Error() {
					t.Fatalf("LabelNamer.Build(%q) should have returned an error %q, but returned: %q", tt.label, tt.wantUTF8Error, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("LabelNamer.Build(%q) returned an error: %s", tt.label, err)
			}
			if got != tt.label {
				t.Errorf("LabelNamer.Build(%q) = %q, want %q", tt.label, got, tt.label)
			}
		})
	}
}

// TestCanFastPathLabel verifies that the fast-path predicate agrees with
// LabelNamer.Build across every entry in labelTestCases × the four configs.
// Whenever canFastPathLabel returns true, Build must succeed and return the
// input unchanged; whenever Build returns the input unchanged for a valid
// input, canFastPathLabel must return true.
func TestCanFastPathLabel(t *testing.T) {
	configs := []struct {
		name                        string
		preserveMultipleUnderscores bool
		underscoreLabelSanitization bool
	}{
		{"default", false, false},
		{"preserve multiple underscores", true, false},
		{"underscore label sanitization", false, true},
		{"both", true, true},
	}
	for _, cfg := range configs {
		t.Run(cfg.name, func(t *testing.T) {
			for _, tt := range labelTestCases {
				t.Run(tt.label, func(t *testing.T) {
					if tt.label == "" {
						return
					}

					namer := LabelNamer{
						PreserveMultipleUnderscores: cfg.preserveMultipleUnderscores,
						UnderscoreLabelSanitization: cfg.underscoreLabelSanitization,
					}
					got, err := namer.Build(tt.label)
					fast := canFastPathLabel(tt.label, cfg.preserveMultipleUnderscores, cfg.underscoreLabelSanitization)

					if fast {
						if err != nil {
							t.Fatalf("canFastPathLabel=true but Build returned error: %s", err)
						}
						if got != tt.label {
							t.Fatalf("canFastPathLabel=true but Build returned %q, want input %q", got, tt.label)
						}
						return
					}

					// fast == false: nothing to assert about Build, except that if
					// it succeeded with output equal to input, the predicate missed a
					// no-op case (a false negative). False negatives are correct but
					// suboptimal; we flag them so the predicate stays tight.
					if err == nil && got == tt.label {
						t.Fatalf("canFastPathLabel=false but Build returned input unchanged (false negative)")
					}
				})
			}
		})
	}
}

// TestLabelNamerBuildZeroAlloc asserts that Build is allocation-free on the
// fast path. The chosen input is already Prometheus-compliant; it must be
// returned unchanged with no heap allocations.
func TestLabelNamerBuildZeroAlloc(t *testing.T) {
	namer := LabelNamer{}
	got := testing.AllocsPerRun(100, func() {
		_, _ = namer.Build("http_method_total")
	})
	if got > 0 {
		t.Fatalf("Build allocated %f times per run on the fast path, want 0", got)
	}
}
