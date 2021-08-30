package feed

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractTextFromHTML(t *testing.T) {
	tcs := []struct {
		HTML string
		text string
	}{
		{
			HTML: ``,
			text: ``,
		},
		{
			HTML: ` <p`,
			text: ``,
		},
		{
			HTML: ` <p>broken`,
			text: `broken`,
		},
		{
			HTML: `foo`,
			text: `foo`,
		},
		{
			HTML: `<p>foo<div>bar</div></p>`,
			text: `foo bar`,
		},
		{
			HTML: "<p> foo <a> bar </a> baz \n<a> quux   </a> \r orb </p>",
			text: `foo bar baz quux orb`,
		},
	}

	for _, tc := range tcs {
		actual := ExtractTextFromHTML(tc.HTML)
		assert.Equal(t, tc.text, actual, "input: %s", tc.HTML)
	}
}

func TestTrimSentences(t *testing.T) {
	tcs := []struct {
		input    string
		expected string
	}{
		{
			input:    ``,
			expected: ``,
		},
		{
			input:    `...............................`,
			expected: `....................`,
		},
		{
			input:    `Foo bar.`,
			expected: `Foo bar.`,
		},
		{
			input:    `Foo bar. Baz quux.`,
			expected: `Foo bar. Baz quux.`,
		},
		{
			input:    `Foo bar. Baz quux. Orb fox.`,
			expected: `Foo bar. Baz quux.`,
		},
		{
			input:    `Foo bar. Baz quux. Orb fox`,
			expected: `Foo bar. Baz quux.`,
		},
		{
			input:    `Foo bar. Baz quux`,
			expected: `Foo bar. Baz quux`,
		},
		{
			// This can be improved to preserve word boundary.
			input:    `Foobar bazquux orb fox`,
			expected: `Foobar bazquux orb f`,
		},
		{
			input:    `1234567890123456789012345`,
			expected: `12345678901234567890`,
		},
		// {
		// // BUG. Rune boundary is broken
		// 	input:    `1234567890123456789ąąxxx`,
		// 	expected: `1234567890123456789ą`,
		// },
	}

	for _, tc := range tcs {
		actual := TrimSentences(tc.input, 20)
		assert.Equal(t, tc.expected, actual, "input: `%s`", tc.input)
	}
}
