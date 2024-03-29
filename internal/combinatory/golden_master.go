package combinatory

import (
	"fmt"
	"strings"
)

/*
Master generates a subject running tests for every parameter combination. We need to provide:

* A function that wraps the subject under test
* An slice of slices of possible parameters

Should be a method of Golden. Not return
*/
func Master(f Wrapper, values ...[]any) []GM {

	all := Generate(values)
	var r []GM
	for idx, combination := range all {
		t := GM{
			Id:     idx + 1,
			Params: joinCommaSeparated(combination),
			Output: f(combination...),
		}
		r = append(r, t)
	}
	return r
}

type GM struct {
	Id     int
	Params string
	Output any
}

func joinCommaSeparated(a []any) string {
	result := make([]string, len(a))
	for i, item := range a {
		result[i] = fmt.Sprintf("%v", item)
	}
	return strings.Join(result, ", ")
}

type Wrapper func(args ...any) any
