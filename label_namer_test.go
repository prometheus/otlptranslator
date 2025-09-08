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
	"testing"
)

var labelTestCases = []struct {
	label                        string
	sanitized                    string
	sanitizedMultipleUnderscores string
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
	{label: "ようこそ", sanitized: ""},
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
}

func TestBuildLabel(t *testing.T) {
	for _, tt := range labelTestCases {
		t.Run(tt.label, func(t *testing.T) {
			t.Run("Not preserving multiple underscores", func(t *testing.T) {
				labelNamer := LabelNamer{}
				got, _ := labelNamer.Build(tt.label)
				if got != tt.sanitized {
					t.Errorf("LabelNamer.Build(%q) = %q, want %q", tt.label, got, tt.sanitized)
				}
			})
			t.Run("Preserving multiple underscores", func(t *testing.T) {
				labelNamer := LabelNamer{PreserveMultipleUnderscores: true}
				got, _ := labelNamer.Build(tt.label)
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
			got, _ := labelNamer.Build(tt.label)
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

func TestBuildLabel_Errors(t *testing.T) {
	labelTestCases := []struct {
		label         string
		wantEscapeErr bool
		wantUTF8err   bool
	}{
		{label: "", wantEscapeErr: true, wantUTF8err: true},
		{label: "__", wantEscapeErr: true, wantUTF8err: true},
		{label: "ようこそ", wantEscapeErr: true}, // Would be __ if UTF-8 isn't allowed
	}
	for _, tt := range labelTestCases {
		t.Run(tt.label, func(t *testing.T) {
			if tt.wantEscapeErr {
				labelNamer := LabelNamer{}
				got, err := labelNamer.Build(tt.label)
				if err == nil {
					t.Errorf("LabelNamer.Build(%q) returned nil err, wanted one, return value %v", tt.label, got)
				}
			}

			if tt.wantUTF8err {
				labelNamer := LabelNamer{UTF8Allowed: true}
				got, err := labelNamer.Build(tt.label)
				if err == nil {
					t.Errorf("LabelNamer.Build(%q) returned nil err, wanted one, return value %v", tt.label, got)
				}
			}
		})
	}
}
