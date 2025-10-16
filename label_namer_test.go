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
	label                                    string
	sanitized                                string
	sanitizedMultipleUnderscores             string
	wantSanitizationError                    error
	wantSanitizationMultipleUnderscoresError error
	wantUTF8Error                            error
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
	{label: "123label", sanitized: "key_123label"},
	{label: "_label_starting_with_underscore", sanitized: "_label_starting_with_underscore"},
	{
		label:                        "__label_starting_with_2underscores",
		sanitized:                    "_label_starting_with_2underscores",
		sanitizedMultipleUnderscores: "__label_starting_with_2underscores",
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

// TestBuildLabel_Underscores confirms that `key_` is only prepended to label
// names starting with an underscore if UnderscoreLabelSanitization is true.
// Labels starting with a number always get `key_` prepended so they are valid
// Prometheus labels.
func TestBuildLabel_Underscores(t *testing.T) {
	labelTestCases := []struct {
		label                string
		sanitized            string
		sanitizedUnderscores string
	}{
		{
			label:                "regular label",
			sanitized:            "regular_label",
			sanitizedUnderscores: "regular_label",
		},
		{
			label:                "123label",
			sanitized:            "key_123label",
			sanitizedUnderscores: "key_123label",
		},
		{
			label:                "_label_starting_with_underscore",
			sanitized:            "_label_starting_with_underscore",
			sanitizedUnderscores: "key_label_starting_with_underscore",
		},
	}
	for _, tt := range labelTestCases {
		t.Run(tt.label, func(t *testing.T) {
			labelNamer := LabelNamer{}
			got, err := labelNamer.Build(tt.label)
			if err != nil {
				t.Errorf("LabelNamer.Build(%q) returned err: %v", tt.label, got)
			}
			if got != tt.sanitized {
				t.Errorf("LabelNamer.Build(%q) = %q, want %q", tt.label, got, tt.sanitized)
			}

			labelNamer.UnderscoreLabelSanitization = true
			got, err = labelNamer.Build(tt.label)
			if err != nil {
				t.Errorf("LabelNamer.Build(%q) (underscore mode) returned err: %v", tt.label, got)
			}

			if got != tt.sanitizedUnderscores {
				t.Errorf("LabelNamer.Build(%q) (underscore mode) = %q, want %q", tt.label, got, tt.sanitizedUnderscores)
			}
		})
	}
}
