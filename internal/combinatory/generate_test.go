package combinatory_test

import (
	"reflect"
	"testing"

	"github.com/franiglesias/golden/internal/combinatory"
	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	tests := []struct {
		name string
		args [][]any
		want [][]any
	}{
		{
			name: "should do nothing if empty",
			args: [][]any{},
			want: [][]any{},
		},
		{
			name: "should be same combination as input if only one parameter",
			args: [][]any{{1, 2, 3}},
			want: [][]any{{1}, {2}, {3}},
		},
		{
			name: "should be combine more parameters",
			args: [][]any{
				{1, 2, 3},
				{"a", "b"},
				{"#", "%"},
			},
			want: [][]any{
				{1, "a", "#"},
				{2, "a", "#"},
				{3, "a", "#"},
				{1, "b", "#"},
				{2, "b", "#"},
				{3, "b", "#"},
				{1, "a", "%"},
				{2, "a", "%"},
				{3, "a", "%"},
				{1, "b", "%"},
				{2, "b", "%"},
				{3, "b", "%"},
			},
		},

		{
			name: "should combine 5 parameters with 2 values each",
			args: [][]any{{"a", "b"}, {"c", "d"}, {"e", "f"}, {"g", "h"}, {"j", "k"}},
			want: [][]any{
				{"a", "c", "e", "g", "j"},
				{"b", "c", "e", "g", "j"},
				{"a", "d", "e", "g", "j"},
				{"b", "d", "e", "g", "j"},
				{"a", "c", "f", "g", "j"},
				{"b", "c", "f", "g", "j"},
				{"a", "d", "f", "g", "j"},
				{"b", "d", "f", "g", "j"},
				{"a", "c", "e", "h", "j"},
				{"b", "c", "e", "h", "j"},
				{"a", "d", "e", "h", "j"},
				{"b", "d", "e", "h", "j"},
				{"a", "c", "f", "h", "j"},
				{"b", "c", "f", "h", "j"},
				{"a", "d", "f", "h", "j"},
				{"b", "d", "f", "h", "j"},
				{"a", "c", "e", "g", "k"},
				{"b", "c", "e", "g", "k"},
				{"a", "d", "e", "g", "k"},
				{"b", "d", "e", "g", "k"},
				{"a", "c", "f", "g", "k"},
				{"b", "c", "f", "g", "k"},
				{"a", "d", "f", "g", "k"},
				{"b", "d", "f", "g", "k"},
				{"a", "c", "e", "h", "k"},
				{"b", "c", "e", "h", "k"},
				{"a", "d", "e", "h", "k"},
				{"b", "d", "e", "h", "k"},
				{"a", "c", "f", "h", "k"},
				{"b", "c", "f", "h", "k"},
				{"a", "d", "f", "h", "k"},
				{"b", "d", "f", "h", "k"},
			},
		},

		{
			name: "should be combine two parameters",
			args: [][]any{
				{1, 2, 3},
				{"a", "b"},
			},
			want: [][]any{
				{1, "a"},
				{2, "a"},
				{3, "a"},
				{1, "b"},
				{2, "b"},
				{3, "b"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := combinatory.Generate(tt.args)
			assert.EqualValues(t, tt.want, got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Generate() = %v, want %v", got, tt.want)
			}
		})
	}
}
