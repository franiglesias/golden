package golden_test

import (
	"github.com/franiglesias/golden"
	"gotest.tools/v3/assert"
	"testing"
)

/*
TestJsonNormalizer: given that JsonNormalizer relies on json.Marshal we can
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
		{
			name:    "should normalize json string as is",
			subject: `{"object":{"id":"12345", "name":"My Object", "count":1234, "validated": true, "other": {"remark": "accept"}}}`,
			want: `{
  "object": {
    "id": "12345",
    "name": "My Object",
    "count": 1234,
    "validated": true,
    "other": {
      "remark": "accept"
    }
  }
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			normalizer := golden.JsonNormalizer{}
			normalized, _ := normalizer.Normalize(tt.subject)
			assert.Equal(t, tt.want, normalized)
		})
	}
}
