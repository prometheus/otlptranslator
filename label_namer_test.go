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
	label         string
	sanitized     string
	wantEscapeErr bool
	wantUTF8err   bool
}{
	{label: "", sanitized: "", wantEscapeErr: true, wantUTF8err: true},
	{label: "__", sanitized: "__", wantEscapeErr: false, wantUTF8err: false},
	{label: "label:with:colons", sanitized: "label_with_colons"},
	{label: "LabelWithCapitalLetters", sanitized: "LabelWithCapitalLetters"},
	{label: "label!with&special$chars)", sanitized: "label_with_special_chars_"},
	{label: "label_with_foreign_characters_字符", sanitized: "label_with_foreign_characters_"},
	{label: "label.with.dots", sanitized: "label_with_dots"},
	{label: "123label", sanitized: "key_123label"},
	{label: "_label_starting_with_underscore", sanitized: "key_label_starting_with_underscore"},
	{label: "__label_starting_with_2underscores", sanitized: "__label_starting_with_2underscores"},
	{label: "ようこそ", sanitized: "", wantEscapeErr: true},
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
		})
	}
}

func TestNormalizeLabelPreserveMultipleUnderscores(t *testing.T) {
	testCases := []struct {
		name    string
		namer   LabelNamer
		label   string
		want    string
		wantErr bool
	}{
		{
			name:  "label with double underscores preserved",
			namer: LabelNamer{UTF8Allowed: false, PreserveMultipleUnderscores: true},
			label: "label__with__double__underscores",
			want:  "label__with__double__underscores",
		},
		{
			name:  "label with multiple underscores and special chars preserved",
			namer: LabelNamer{UTF8Allowed: false, PreserveMultipleUnderscores: true},
			label: "label.name__with&&special##chars",
			want:  "label_name__with__special__chars",
		},
		{
			name:  "reserved label with multiple underscores preserved",
			namer: LabelNamer{UTF8Allowed: false, PreserveMultipleUnderscores: true},
			label: "__reserved__label__name__",
			want:  "__reserved__label__name__",
		},
		{
			name:  "reserved label preserved even when PreserveMultipleUnderscores false",
			namer: LabelNamer{UTF8Allowed: false, PreserveMultipleUnderscores: false},
			label: "__reserved__label__name__",
			want:  "__reserved__label__name__",
		},
		{
			name:  "label with trailing multiple underscores",
			namer: LabelNamer{UTF8Allowed: false, PreserveMultipleUnderscores: true},
			label: "trailing_underscores___",
			want:  "trailing_underscores___",
		},
		// Default behavior (collapse multiple underscores)
		{
			name:  "label with double underscores collapsed by default",
			namer: LabelNamer{UTF8Allowed: false, PreserveMultipleUnderscores: false},
			label: "label__with__double__underscores",
			want:  "label_with_double_underscores",
		},
		{
			name:  "label with special chars and multiple underscores collapsed",
			namer: LabelNamer{UTF8Allowed: false, PreserveMultipleUnderscores: false},
			label: "label.name__with&&special##chars",
			want:  "label_name_with_special_chars",
		},
		// UTF8Allowed = true should ignore PreserveMultipleUnderscores
		{
			name:  "UTF8 mode ignores PreserveMultipleUnderscores setting",
			namer: LabelNamer{UTF8Allowed: true, PreserveMultipleUnderscores: false},
			label: "label__with__special.chars",
			want:  "label__with__special.chars",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.namer.Build(tt.label)
			if tt.wantErr {
				if err == nil {
					t.Errorf("LabelNamer.Build(%q) = %q, want error", tt.label, got)
				}
				return
			}
			if err != nil {
				t.Errorf("LabelNamer.Build(%q) returned error %v, want nil", tt.label, err)
				return
			}
			if got != tt.want {
				t.Errorf("LabelNamer.Build(%q) = %q, want %q", tt.label, got, tt.want)
			}
		})
	}
}
