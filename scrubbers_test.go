package golden_test

import (
	"github.com/franiglesias/golden"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBasicRegexpScrubbing(t *testing.T) {
	t.Run("should not replace anything if no match", func(t *testing.T) {
		subject := "A string not suspicions of contain anything to remove"
		scrubber := golden.NewRegexpScrubber("\\d{2}-\\d{2}-\\d{2}", "24-01-15")
		result := scrubber.Clean(subject)
		assert.Equal(t, subject, result)
	})

	t.Run("should replace dates", func(t *testing.T) {
		subject := "The next days 24-01-30, 24-02-03 and 24-02-10 we will be closed."
		scrubber := golden.NewRegexpScrubber("\\d{2}-\\d{2}-\\d{2}", "24-01-15")
		result := scrubber.Clean(subject)
		expected := "The next days 24-01-15, 24-01-15 and 24-01-15 we will be closed."
		assert.Equal(t, expected, result)
	})
}
