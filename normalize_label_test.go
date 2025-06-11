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
	"fmt"
	"testing"
)

func TestNormalizeLabel(t *testing.T) {
	tests := []struct {
		label    string
		expected string
	}{
		{"", ""},
		{"label:with:colons", "label_with_colons"},
		{"LabelWithCapitalLetters", "LabelWithCapitalLetters"},
		{"label!with&special$chars)", "label_with_special_chars_"},
		{"label_with_foreign_characters_字符", "label_with_foreign_characters___"},
		{"label.with.dots", "label_with_dots"},
		{"123label", "key_123label"},
		{"_label_starting_with_underscore", "key_label_starting_with_underscore"},
		{"__label_starting_with_2underscores", "__label_starting_with_2underscores"},
	}

	for i, test := range tests {
		for _, utf8Allowed := range []bool{true, false} {
			t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
				labelNamer := LabelNamer{UTF8Allowed: utf8Allowed}
				result := labelNamer.Build(test.label)

				if utf8Allowed {
					if test.label != result {
						t.Errorf("labelNamer.Build(%q), got %q, want %q", test.label, result, test.label)
					}
				} else {
					if test.expected != result {
						t.Errorf("labelNamer.Build(%q), got %q, want %q", test.label, result, test.expected)
					}
				}
			})
		}
	}
}
