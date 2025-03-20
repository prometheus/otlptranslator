package prometheus

import "testing"

var labelBenchmarkInputs = []string{
	"",
	"label:with:colons",
	"LabelWithCapitalLetters",
	"label!with&special$chars)",
	"label_with_foreign_characters_字符",
	"label.with.dots",
	"123label",
	"_label_starting_with_underscore",
	"__label_starting_with_2underscores",
}

func BenchmarkNormalizeLabel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, input := range labelBenchmarkInputs {
			NormalizeLabel(input)
		}
	}
}
