package helper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

/*
TSpy is a replacement of *testing.T for some tests of the golden library. With
it, we can spy if the Verify method fails when differences between subject and
snapshot are found
*/
type TSpy struct {
	*testing.T
	failed bool
	report string
}

func (t *TSpy) Errorf(_ string, report ...any) {
	t.failed = true
	t.report = report[0].(string)
}

func (t *TSpy) Reset() {
	t.failed = false
	t.report = ""
}

/*
AssertFailedTest allows us to spy on TSpy
*/
func AssertFailedTest(t *testing.T, gt *TSpy) {
	assert.True(t, gt.failed, "Test passed and it shouldn't")
}

func AssertPassTest(t *testing.T, gt *TSpy) {
	assert.False(t, gt.failed, "Test failed and it shouldn't")
}

func AssertReportContains(t *testing.T, g *TSpy, s string) {
	assert.Containsf(t, g.report, s, "Diff report doesn't contains expected '%s'", s)
}
