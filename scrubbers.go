package golden

import (
	"fmt"
	"regexp"
)

/*
Scrubber modifies the subject, usually to replace not deterministic data for
some fixed replacement, but you can create and apply scrubbers for cleaning the
subject or obfuscate sensible data

The basic Scrubber is a generic Scrubber that searches for a regex pattern and replaces it
*/
type Scrubber struct {
	pattern     string
	replacement string
}

func NewScrubber(pattern, replacement string) Scrubber {
	return Scrubber{
		pattern:     pattern,
		replacement: replacement,
	}
}

func (b Scrubber) Clean(subject string) string {
	re := regexp.MustCompile(b.pattern)
	return re.ReplaceAllString(subject, b.replacement)
}

/*

 Custom Scrubbers

*/

/*
CreditCard obfuscates credit card numbers
*/
func CreditCard() Scrubber {
	return NewScrubber(
		"\\d{4}-\\d{4}-\\d{4}-",
		"****-****-****-",
	)
}

/*
Format will help to limit the scrubbing to string that matches the format. Use
the placeholder %s to indicate the part that you want to be scrubbed

	golden.Format("Credit Card: %s", golden.CreditCard())

will apply the CreditCard obfuscation only if it finds a Credit Card Number
after a "Credit Card: " string. This way you can avoid scrubbing parts of the
subject that you don't want to touch.
*/
func Format(f string, s Scrubber) Scrubber {
	return NewScrubber(
		fmt.Sprintf(f, s.pattern),
		fmt.Sprintf(f, s.replacement),
	)
}
