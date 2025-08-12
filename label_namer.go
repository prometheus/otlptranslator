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
// Provenance-includes-location: https://github.com/prometheus/prometheus/blob/93e991ef7ed19cc997a9360c8016cac3767b8057/storage/remote/otlptranslator/prometheus/normalize_label.go
// Provenance-includes-license: Apache-2.0
// Provenance-includes-copyright: Copyright The Prometheus Authors
// Provenance-includes-location: https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/95e8f8fdc2a9dc87230406c9a3cf02be4fd68bea/pkg/translator/prometheus/normalize_label.go
// Provenance-includes-license: Apache-2.0
// Provenance-includes-copyright: Copyright The OpenTelemetry Authors.

package otlptranslator

import (
	"fmt"
	"strings"
	"unicode"
)

// LabelNamer is a helper struct to build label names.
// It translates OpenTelemetry Protocol (OTLP) attribute names to Prometheus-compliant label names.
//
// Example usage:
//
//	namer := LabelNamer{UTF8Allowed: false}
//	result := namer.Build("http.method") // "http_method"
type LabelNamer struct {
	UTF8Allowed bool
	// PreserveMultipleUnderscores when true, preserves multiple consecutive underscores
	// in label names when UTF8Allowed is false. This option is discouraged
	// as it violates the OpenTelemetry to Prometheus specification
	// (https://github.com/open-telemetry/opentelemetry-specification/blob/v1.38.0/specification/compatibility/prometheus_and_openmetrics.md#otlp-metric-points-to-prometheus)
	// but may be needed for compatibility with legacy systems that allow multiple underscores.
	// This option is ignored when UTF8Allowed is true.
	PreserveMultipleUnderscores bool
}

// Build normalizes the specified label to follow Prometheus label names standard.
//
// Translation rules:
//   - Replaces invalid characters with underscores
//   - Prefixes labels with invalid start characters (numbers or `_`) with "key"
//   - Preserves double underscore labels (reserved names)
//   - If UTF8Allowed is true, returns label as-is
//
// Examples:
//
//	namer := LabelNamer{UTF8Allowed: false}
//	namer.Build("http.method")     // "http_method"
//	namer.Build("123invalid")      // "key_123invalid"
//	namer.Build("__reserved__")    // "__reserved__" (preserved)
func (ln *LabelNamer) Build(label string) (normalizedName string, err error) {
	defer func() {
		if len(normalizedName) == 0 {
			err = fmt.Errorf("normalization for label name %q resulted in empty name", label)
			return
		}

		if ln.UTF8Allowed || normalizedName == label {
			return
		}

		// Check that the resulting normalized name contains at least one non-underscore character
		for _, c := range normalizedName {
			if c != '_' {
				return
			}
		}
		err = fmt.Errorf("normalization for label name %q resulted in invalid name %q", label, normalizedName)
		normalizedName = ""
	}()

	// Trivial case.
	if len(label) == 0 || ln.UTF8Allowed {
		normalizedName = label
		return
	}

	// Reserved labels (starting with exactly __) should be preserved entirely
	// This includes labels like __name__, __address__, etc.
	if len(label) >= 2 && strings.HasPrefix(label, "__") && (len(label) == 2 || label[2] != '_') {
		normalizedName = label
		return
	}

	normalizedName = sanitizeLabelName(label, ln.PreserveMultipleUnderscores)

	// Check if the sanitized name contains only underscores (invalid)
	hasNonUnderscore := false
	for _, c := range normalizedName {
		if c != '_' {
			hasNonUnderscore = true
			break
		}
	}
	if !hasNonUnderscore {
		err = fmt.Errorf("normalization for label name %q resulted in invalid name %q", label, normalizedName)
		normalizedName = ""
		return
	}

	// If label starts with a number, prepend with "key_".
	if unicode.IsDigit(rune(normalizedName[0])) {
		normalizedName = "key_" + normalizedName
	} else if strings.HasPrefix(normalizedName, "_") {
		// Add "key" prefix for labels starting with underscores (unless they are reserved)
		// Reserved labels starting with exactly __ are already handled above
		normalizedName = "key" + normalizedName
	}

	return
}

// sanitizeLabelName replaces any characters not valid according to the
// classical Prometheus label naming scheme with an underscore.
// When preserveMultipleUnderscores is true or by default, multiple consecutive underscores are preserved.
// When false, multiple consecutive underscores are collapsed to single underscores (new behavior).
func sanitizeLabelName(name string, preserveMultipleUnderscores bool) string {
	if preserveMultipleUnderscores {
		// Simple case: just replace invalid characters, preserve multiple underscores
		var b strings.Builder
		b.Grow(len(name))
		for _, r := range name {
			if isValidCompliantLabelChar(r) {
				b.WriteRune(r)
			} else {
				b.WriteRune('_')
			}
		}
		return b.String()
	}

	// Collapse multiple underscores while replacing invalid characters
	var b strings.Builder
	b.Grow(len(name))
	prevWasUnderscore := false

	for _, r := range name {
		if isValidCompliantLabelChar(r) {
			b.WriteRune(r)
			prevWasUnderscore = false
		} else if !prevWasUnderscore {
			// Invalid character - replace with underscore (collapse consecutive underscores)
			b.WriteRune('_')
			prevWasUnderscore = true
		}
		// If prevWasUnderscore is true, skip this underscore (collapse)
	}
	return b.String()
}

// isValidCompliantLabelChar checks if a rune is a valid label name character (a-z, A-Z, 0-9).
func isValidCompliantLabelChar(r rune) bool {
	return (r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		(r >= '0' && r <= '9')
}
