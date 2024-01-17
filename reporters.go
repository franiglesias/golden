package golden

import (
	"fmt"
	"github.com/andreyvit/diff"
)

const diffHeaderFormat = "\nDifferences found:\n==================\n%s\n"
const noDifferences = "No differences found."

type CharDiffReporter struct {
}

func (c CharDiffReporter) Differences(want, got string) string {
	if want == got {
		return noDifferences
	}
	diffs := diff.CharacterDiff(want, got)
	return fmt.Sprintf(diffHeaderFormat, diffs)
}

type LineDiffReporter struct{}

func (LineDiffReporter) Differences(want, got string) string {
	if want == got {
		return "No differences found."
	}
	diffs := diff.LineDiff(want, got)
	return fmt.Sprintf(diffHeaderFormat, diffs)
}
