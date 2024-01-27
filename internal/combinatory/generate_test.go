package combinatory_test

import (
	"github.com/franiglesias/golden/internal/combinatory"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
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
