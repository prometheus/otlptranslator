package otlptranslator

import "testing"

var labelBenchmarkInputs = []struct {
	name  string
	label string
}{
	{
		name:  "empty label",
		label: "",
	},
	{
		name:  "label with colons",
		label: "label:with:colons",
	},
	{
		name:  "label with capital letters",
		label: "LabelWithCapitalLetters",
	},
	{
		name:  "label with special characters",
		label: "label!with&special$chars)",
	},
	{
		name:  "label with foreign characters",
		label: "label_with_foreign_characters_字符",
	},
	{
		name:  "label with dots",
		label: "label.with.dots",
	},
	{
		name:  "label starting with digits",
		label: "123label",
	},
	{
		name:  "label starting with underscores",
		label: "_label_starting_with_underscore",
	},
	{
		name:  "label starting with 2 underscores",
		label: "__label_starting_with_2underscores",
	},
}

func BenchmarkNormalizeLabel(b *testing.B) {
	for _, input := range labelBenchmarkInputs {
		b.Run(input.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				NormalizeLabel(input.label)
			}
		})
	}
}
