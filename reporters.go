package golden

import (
	godiff "codeberg.org/h7c/go-diff"
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
		return noDifferences
	}
	diffs := diff.LineDiff(want, got)
	return fmt.Sprintf(diffHeaderFormat, diffs)
}

type BetterDiffReporter struct {
	color bool
}

func NewBetterDiffReporter() BetterDiffReporter {
	return BetterDiffReporter{
		color: true,
	}
}

func NewBetterDiffReporterWithoutColor() BetterDiffReporter {
	return BetterDiffReporter{
		color: false,
	}
}

func (b BetterDiffReporter) Differences(want, got string) string {
	w := godiff.NewFileFromString(want)
	g := godiff.NewFileFromString(got)

	if w.IsDifferentFrom(g) {
		diffs := godiff.GetDiff(w, g, b.color)
		return fmt.Sprintf(diffHeaderFormat, diffs)
	}

	return noDifferences
}
