package golden

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

/*
AssertFailedTest allows us to spy on TSpy
*/
func AssertFailedTest(t *testing.T, gt *TSpy) {
	assert.True(t, gt.failed)
}

func AssertReportContains(t *testing.T, g *TSpy, s string) {
	assert.Contains(t, g.report, s)
}
