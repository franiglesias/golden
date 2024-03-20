package golden_test

import (
	"github.com/franiglesias/golden"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCharDiffReporter(t *testing.T) {

	reporter := golden.CharDiffReporter{}

	t.Run("show no differences", func(t *testing.T) {
		result := reporter.Differences("Same content", "Same content")
		assert.Equal(t, "No differences found.", result)
	})

	t.Run("show char differences", func(t *testing.T) {
		result := reporter.Differences("Wanted this.", "Gotten that.")
		assert.Contains(t, result, "Differences found:")
		assert.Contains(t, result, "(~~Wanted this~~)(++Gotten that++)")
	})
}

func TestLineDiffReporter(t *testing.T) {

	reporter := golden.LineDiffReporter{}

	t.Run("show no differences", func(t *testing.T) {
		result := reporter.Differences("Same content", "Same content")
		assert.Equal(t, "No differences found.", result)
	})

	t.Run("show char differences", func(t *testing.T) {
		result := reporter.Differences("Wanted this.", "Gotten that.")
		assert.Contains(t, result, "Differences found:")
		assert.Contains(t, result, "-Wanted this.")
		assert.Contains(t, result, "+Gotten that.")
	})
}

func TestBetterDiffReporter(t *testing.T) {

	reporter := golden.NewBetterDiffReporter()

	t.Run("show no differences", func(t *testing.T) {
		result := reporter.Differences("Same content", "Same content")

		assert.Equal(t, "No differences found.", result)
	})

	t.Run("show differences", func(t *testing.T) {
		result := reporter.Differences("Wanted this.", "Gotten that.")

		assert.Contains(t, result, "Differences found:")
		assert.Contains(t, result, "- Wanted this.")
		assert.Contains(t, result, "+ Gotten that.")
	})
}
