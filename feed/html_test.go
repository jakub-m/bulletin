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
