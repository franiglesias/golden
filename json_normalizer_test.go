package golden_test

import (
	"golden"
	"gotest.tools/v3/assert"
	"testing"
)

/*
TestJsonNormalizer: given that JsonNormalizer relies on json.Marshall we can
trust that any object should be correctly marshaled. This test documents how to
use it, and that the output is cleaned of leading and trailing space chars, so
we can be reasonably sure that the normalized version will be consistent and the
snapshot can be compared with the subject without showing irrelevant differences.
*/
func TestJsonNormalizer(t *testing.T) {
	tests := []struct {
		name    string
		subject any
		want    string
	}{
		{
			name:    "should normalize string",
			subject: "This is a string",
			want:    "This is a string",
		},
		{
			name:    "should normalize number to string",
			subject: 123.45,
			want:    "123.45",
		},
		{
			name: "should normalize slice",
			subject: []string{
				"Item 1",
				"Item 2",
			},
			want: `[
  "Item 1",
  "Item 2"
]`,
		},
		{name: "should remove leading and trailing spaces",
			subject: "   This is a string   ",
			want:    "This is a string",
		},

		{name: "should remove leading and trailing new lines",
			subject: "\nThis is a string\n",
			want:    "This is a string",
		},
	}

	for _, tt := range tests {
		normalizer := golden.JsonNormalizer{}
		normalized, _ := normalizer.Normalize(tt.subject)
		assert.Equal(t, tt.want, normalized)
	}
}
