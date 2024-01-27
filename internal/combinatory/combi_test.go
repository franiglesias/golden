package combinatory_test

import (
	"github.com/franiglesias/golden/internal/combinatory"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func MyFunc(part string, times int) string {
	return strings.Repeat(part, times)
}

func Border(title string, part string, span int) string {
	width := span*2 + len(title) + 2
	top := strings.Repeat(part, width)
	body := part + strings.Repeat(" ", span) + title + strings.Repeat(" ", span) + part
	return top + "\n" + body + "\n" + top + "\n"
}

func TestExample(t *testing.T) {
	toTest := func(args ...any) any {
		return MyFunc(args[0].(string), args[1].(int))
	}

	parts := []any{"-", "=", "*", "#"}
	times := []any{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	master := combinatory.Master(toTest, parts, times)
	assert.Len(t, master, 44)
}

func TestWithThreeParametersExample(t *testing.T) {
	toTest := func(args ...any) any {
		return Border(args[0].(string), args[1].(string), args[2].(int))
	}

	titles := []any{"Example 1", "Example long enough", "Another thing"}
	parts := []any{"-", "=", "*", "#"}
	times := []any{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	master := combinatory.Master(toTest, titles, parts, times)
	assert.Len(t, master, 132)
}
