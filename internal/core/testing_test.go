package core

import (
	"regexp"
	"testing"

	"github.com/alecthomas/assert"
)

func TestGoldenIgnoreLines(t *testing.T) {
	original := `
Line1
Line2
Line3
Line4`
	expected := `
Line1
Line4`
	actual, err := goldenReplacePatterns(original, GoldenReplacement{
		Pattern:     regexp.MustCompile("Line2\nLine3\n"),
		Replacement: "",
	})
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)

	expected2 := `
Line4
Line3
Line2
Line1`
	actual2, err := goldenReplacePatterns(original,
		GoldenReplacement{
			Pattern:     regexp.MustCompile("(?s)(Line1).*(Line2).*(Line3).*(Line4)"),
			Replacement: "$4\n$3\n$2\n$1",
		})
	assert.Nil(t, err)
	assert.Equal(t, expected2, actual2)
}
