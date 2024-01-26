package golden_test

import (
	"github.com/franiglesias/golden"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBasicRegexpScrubbing(t *testing.T) {
	t.Run("should not replace anything if no match", func(t *testing.T) {
		subject := "A string not suspicions of contain anything to remove"
		scrubber := golden.NewScrubber("\\d{2}-\\d{2}-\\d{2}", "24-01-15")
		result := scrubber.Clean(subject)
		assert.Equal(t, subject, result)
	})

	t.Run("should replace dates", func(t *testing.T) {
		subject := "The next days 24-01-30, 24-02-03 and 24-02-10 we will be closed."
		scrubber := golden.NewScrubber("\\d{2}-\\d{2}-\\d{2}", "24-01-15")
		result := scrubber.Clean(subject)
		expected := "The next days 24-01-15, 24-01-15 and 24-01-15 we will be closed."
		assert.Equal(t, expected, result)
	})
}

func TestCreditCard(t *testing.T) {
	scrubber := golden.CreditCard()
	subject := "Credit card: 1234-5678-9012-1234"
	assert.Equal(t, "Credit card: ****-****-****-1234", scrubber.Clean(subject))
}

func TestFormatScrubber(t *testing.T) {
	t.Run("should obfuscated only credit card number", func(t *testing.T) {
		scrubber := golden.Format("Credit card: %s", golden.CreditCard())
		subject := "Credit card: 1234-5678-9012-1234, Another code: 4561-1234-4532-6543"
		assert.Equal(t, "Credit card: ****-****-****-1234, Another code: 4561-1234-4532-6543", scrubber.Clean(subject))
	})

	t.Run("should obfuscate all numbers", func(t *testing.T) {
		scrubber := golden.CreditCard()
		subject := "Credit card: 1234-5678-9012-1234, Another code: 4561-1234-4532-6543"
		assert.Equal(t, "Credit card: ****-****-****-1234, Another code: ****-****-****-6543", scrubber.Clean(subject))
	})
}
