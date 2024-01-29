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
	pattern     string // The regexp pattern that describe what you are looking for
	replacement string // Default replacement
	format      string // A format string that allow you to delimite the scope for replacement
}

func NewScrubber(pattern, replacement string, opts ...ScrubberOption) Scrubber {
	s := Scrubber{
		pattern:     pattern,
		replacement: replacement,
		format:      "%s",
	}

	for _, opt := range opts {
		opt(&s)
	}
	return s
}

func (b Scrubber) Clean(subject string) string {
	re := regexp.MustCompile(fmt.Sprintf(b.format, b.pattern))
	return re.ReplaceAllString(subject, fmt.Sprintf(b.format, b.replacement))
}

/*

## Custom Scrubbers

*/

/*
CreditCard obfuscates credit card numbers
*/
func CreditCard(opts ...ScrubberOption) Scrubber {
	return NewScrubber(
		"\\d{4}-\\d{4}-\\d{4}-",
		"****-****-****-",
		opts...,
	)
}

/*
ULID replaces Unique Lexicographic Identifiers with <ULID> placeholder

	ulidScrubber := golden.ULID()

or you can specify a custom replacement

	fixedULID := golden.ULID(golden.Replacement("01HNB10NSJS26X2RTERPZTM0KB"))
	anotherPlaceHolder := golden.ULID(golden.Replacement("[ULID here]]"))
*/
func ULID(opts ...ScrubberOption) Scrubber {
	return NewScrubber("[0-9A-Za-z]{26}", "<ULID>", opts...)
}

/*

## Scrubber options

*/

type ScrubberOption func(s *Scrubber)

/*
Replacement define a custom replacement for any specialized Scrubber

	golden.ULID(golden.Replacement("[ULID comes here]"))
*/
func Replacement(r string) ScrubberOption {
	return func(s *Scrubber) {
		s.replacement = r
	}
}

/*
Format will help to limit the scrubbing to string that matches the format. Use
the placeholder %s to indicate the part that you want to be scrubbed

	golden.CreditCard(golden.Format("Credit Card: %s"))

will apply the CreditCard obfuscation only if it finds a Credit Card Number
after a "Credit Card: " string. This way you can avoid scrubbing parts of the
subject that you don't want to touch.
*/
func Format(f string) ScrubberOption {
	return func(s *Scrubber) {
		s.format = f
	}
}
