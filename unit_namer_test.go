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

package otlptranslator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnitNamer_Build(t *testing.T) {
	tests := []struct {
		name     string
		namer    UnitNamer
		unit     string
		expected string
	}{
		// Basic unit mappings
		{
			name: "empty unit",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "",
			expected: "",
		},
		{
			name: "whitespace only unit",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "   ",
			expected: "",
		},

		// Time units
		{
			name: "seconds unit",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "s",
			expected: "seconds",
		},
		{
			name: "milliseconds unit",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "ms",
			expected: "milliseconds",
		},
		{
			name: "microseconds unit",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "us",
			expected: "microseconds",
		},
		{
			name: "nanoseconds unit",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "ns",
			expected: "nanoseconds",
		},
		{
			name: "minutes unit",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "min",
			expected: "minutes",
		},
		{
			name: "hours unit",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "h",
			expected: "hours",
		},
		{
			name: "days unit",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "d",
			expected: "days",
		},

		// Byte units
		{
			name: "bytes unit",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "By",
			expected: "bytes",
		},
		{
			name: "kibibytes unit",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "KiBy",
			expected: "kibibytes",
		},
		{
			name: "mebibytes unit",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "MiBy",
			expected: "mebibytes",
		},
		{
			name: "gibibytes unit",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "GiBy",
			expected: "gibibytes",
		},
		{
			name: "tibibytes unit",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "TiBy",
			expected: "tibibytes",
		},
		{
			name: "kilobytes unit",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "KBy",
			expected: "kilobytes",
		},
		{
			name: "megabytes unit",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "MBy",
			expected: "megabytes",
		},
		{
			name: "gigabytes unit",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "GBy",
			expected: "gigabytes",
		},
		{
			name: "terabytes unit",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "TBy",
			expected: "terabytes",
		},

		// SI units
		{
			name: "meters unit",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "m",
			expected: "meters",
		},
		{
			name: "volts unit",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "V",
			expected: "volts",
		},
		{
			name: "amperes unit",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "A",
			expected: "amperes",
		},
		{
			name: "joules unit",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "J",
			expected: "joules",
		},
		{
			name: "watts unit",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "W",
			expected: "watts",
		},
		{
			name: "grams unit",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "g",
			expected: "grams",
		},

		// Misc units
		{
			name: "celsius unit",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "Cel",
			expected: "celsius",
		},
		{
			name: "hertz unit",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "Hz",
			expected: "hertz",
		},
		{
			name: "percent unit",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "%",
			expected: "percent",
		},
		{
			name: "dimensionless unit (1)",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "1",
			expected: "",
		},

		// Per units - basic
		{
			name: "per second",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "1/s",
			expected: "per_second",
		},
		{
			name: "per minute",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "1/m",
			expected: "per_minute",
		},
		{
			name: "per hour",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "1/h",
			expected: "per_hour",
		},
		{
			name: "per day",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "1/d",
			expected: "per_day",
		},
		{
			name: "per week",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "1/w",
			expected: "per_week",
		},
		{
			name: "per month",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "1/mo",
			expected: "per_month",
		},
		{
			name: "per year",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "1/y",
			expected: "per_year",
		},

		// Combined main unit and per unit
		{
			name: "bytes per second",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "By/s",
			expected: "bytes_per_second",
		},
		{
			name: "kilobytes per second",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "KBy/s",
			expected: "kilobytes_per_second",
		},
		{
			name: "megabytes per hour",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "MBy/h",
			expected: "megabytes_per_hour",
		},
		{
			name: "requests per minute",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "requests/m",
			expected: "requests_per_minute",
		},
		{
			name: "watts per hour",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "W/h",
			expected: "watts_per_hour",
		},

		// Custom units
		{
			name: "custom main unit",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "custom_unit",
			expected: "custom_unit",
		},
		{
			name: "custom per unit",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "1/custom",
			expected: "per_custom",
		},
		{
			name: "custom main and per unit",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "custom/custom",
			expected: "custom_per_custom",
		},

		// Units with braces (should be ignored)
		{
			name: "main unit with braces ignored",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "{custom}",
			expected: "",
		},
		{
			name: "per unit with braces ignored",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "By/{custom}",
			expected: "bytes",
		},
		{
			name: "main unit with braces ignored, valid per unit",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "{custom}/s",
			expected: "per_second",
		},
		{
			name: "both units with braces ignored",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "{custom}/{invalid}",
			expected: "",
		},

		// Special characters in units (UTF8Allowed = false)
		{
			name: "unit with special characters cleaned (UTF8Allowed = false)",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "custom@unit#with$chars",
			expected: "custom_unit_with_chars",
		},
		{
			name: "per unit with special characters cleaned (UTF8Allowed = false)",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "By/custom@per#unit",
			expected: "bytes_per_custom_per_unit",
		},
		{
			name: "multiple consecutive special chars replaced (UTF8Allowed = false)",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "unit@@##$$name",
			expected: "unit_name",
		},
		{
			name: "multiple consecutive underscores normalized (UTF8Allowed = false)",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "unit__with__multiple__underscores",
			expected: "unit_with_multiple_underscores",
		},
		{
			name: "unit with only special characters (UTF8Allowed = false)",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "@#$%",
			expected: "",
		},

		// Special characters in units (UTF8Allowed = true)
		{
			name: "unit with special characters cleaned (UTF8Allowed = true)",
			namer: UnitNamer{
				UTF8Allowed: true,
			},
			unit:     "custom@unit#with$chars",
			expected: "custom@unit#with$chars",
		},
		{
			name: "per unit with special characters cleaned (UTF8Allowed = true)",
			namer: UnitNamer{
				UTF8Allowed: true,
			},
			unit:     "By/custom@per#unit",
			expected: "bytes_per_custom@per#unit",
		},
		{
			name: "multiple consecutive special chars replaced (UTF8Allowed = true)",
			namer: UnitNamer{
				UTF8Allowed: true,
			},
			unit:     "unit@@##$$name",
			expected: "unit@@##$$name",
		},
		{
			name: "multiple consecutive underscores normalized (UTF8Allowed = true)",
			namer: UnitNamer{
				UTF8Allowed: true,
			},
			unit:     "unit__with__multiple__underscores",
			expected: "unit__with__multiple__underscores",
		},
		{
			name: "unit with only special characters (UTF8Allowed = true)",
			namer: UnitNamer{
				UTF8Allowed: true,
			},
			unit:     "@#$%",
			expected: "@#$%",
		},

		// Edge cases with slashes
		{
			name: "only slash",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "/",
			expected: "",
		},
		{
			name: "main unit with trailing slash",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "By/",
			expected: "bytes",
		},
		{
			name: "slash with empty per unit",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "seconds/",
			expected: "seconds",
		},
		{
			name: "leading slash with per unit",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "/s",
			expected: "per_second",
		},

		// Edge cases with underscores and normalization
		{
			name: "unit starting with underscore",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "_custom_unit",
			expected: "custom_unit",
		},
		{
			name: "unit ending with underscore",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "custom_unit_",
			expected: "custom_unit",
		},
		{
			name: "per unit with underscores normalized",
			namer: UnitNamer{
				UTF8Allowed: false,
			},
			unit:     "By/_custom__unit_",
			expected: "bytes_per_custom_unit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.namer.Build(tt.unit)
			require.Equal(t, tt.expected, got)
		})
	}
}
