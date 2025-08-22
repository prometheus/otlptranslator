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
	label               string
	sanitized           string
	sanitizedPermissive string
	wantEscapeErr       bool
	wantUTF8err         bool
}{
	{
		label:               "",
		sanitized:           "",
		sanitizedPermissive: "",
		wantEscapeErr:       true,
		wantUTF8err:         true,
	},
	{
		label:               "__",
		sanitized:           "__",
		sanitizedPermissive: "__",
		wantEscapeErr:       false,
		wantUTF8err:         false,
	},
	{
		label:               "label:with:colons",
		sanitized:           "label_with_colons",
		sanitizedPermissive: "label_with_colons",
	},
	{
		label:               "LabelWithCapitalLetters",
		sanitized:           "LabelWithCapitalLetters",
		sanitizedPermissive: "LabelWithCapitalLetters",
	},
	{
		label:               "label!with&special$chars)",
		sanitized:           "label_with_special_chars_",
		sanitizedPermissive: "label_with_special_chars_",
	},
	{
		label:               "label_with_foreign_characters_字符",
		sanitized:           "label_with_foreign_characters___",
		sanitizedPermissive: "label_with_foreign_characters___",
	},
	{
		label:               "label.with.dots",
		sanitized:           "label_with_dots",
		sanitizedPermissive: "label_with_dots",
	},
	{
		label:               "123label",
		sanitized:           "key_123label",
		sanitizedPermissive: "key_123label",
	},
	{
		label:               "_label_starting_with_underscore",
		sanitized:           "key_label_starting_with_underscore",
		sanitizedPermissive: "_label_starting_with_underscore",
	},
	{
		label:               "__label_starting_with_2underscores",
		sanitized:           "__label_starting_with_2underscores",
		sanitizedPermissive: "__label_starting_with_2underscores",
	},
	{
		label:               "ようこそ",
		sanitized:           "",
		sanitizedPermissive: "",
		wantEscapeErr:       true,
	},
}

func TestNormalizeLabel(t *testing.T) {
	for _, tt := range labelTestCases {
		t.Run(tt.label, func(t *testing.T) {
			labelNamer := LabelNamer{}
			got, err := labelNamer.Build(tt.label)
			if tt.wantEscapeErr {
				if err == nil {
					t.Errorf("LabelNamer.Build(%q) returned nil err, wanted one, return value %v", tt.label, got)
				}
			} else if err != nil {
				t.Errorf("LabelNamer.Build(%q) returned error %v, want nil", tt.label, err)
			}
			if got != tt.sanitized {
				t.Errorf("LabelNamer.Build(%q) = %q, want %q", tt.label, got, tt.sanitized)
			}
			if tt.wantEscapeErr {
				return
			}
			labelNamer.PermissiveLabelSanitization = true
			got, err = labelNamer.Build(tt.label)
			if err != nil {
				t.Errorf("LabelNamer.Build(%q) (permissive) returned error %v, want nil", tt.label, err)
			}
			if got != tt.sanitizedPermissive {
				t.Errorf("LabelNamer.Build(%q) (permissive) = %q, want %q", tt.label, got, tt.sanitized)
			}
		})
	}
}

func TestNormalizeLabelUTF8Allowed(t *testing.T) {
	for _, tt := range labelTestCases {
		t.Run(tt.label, func(t *testing.T) {
			labelNamer := LabelNamer{UTF8Allowed: true}
			got, err := labelNamer.Build(tt.label)
			if tt.wantUTF8err {
				if err == nil {
					t.Errorf("LabelNamer.Build(%q) returned nil err, wanted one, return value %v", tt.label, got)
				}
			} else if err != nil {
				t.Errorf("LabelNamer.Build(%q) returned error %v, want nil", tt.label, err)
			}
			if got != tt.label {
				t.Errorf("LabelNamer.Build(%q) = %q, want %q", tt.label, got, tt.label)
			}
			if tt.wantEscapeErr {
				return
			}
			labelNamer.PermissiveLabelSanitization = true
			got, err = labelNamer.Build(tt.label)
			if err != nil {
				t.Errorf("LabelNamer.Build(%q) (permissive) returned error %v, want nil", tt.label, err)
			}
			if got != tt.label {
				t.Errorf("LabelNamer.Build(%q) (permissive) = %q, want %q", tt.label, got, tt.sanitized)
			}
		})
	}
}
