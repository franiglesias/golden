package golden

import (
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"regexp"
)

type Scrubber interface {
	Clean(subject string) string
}

type baseScrubber struct {
	target      string
	replacement string
	context     string
}

func (s baseScrubber) Clean(subject string) string {
	panic("Implement Clean method")
}

/*
RegexpScrubber modifies the subject, usually to replace not deterministic data for
some fixed replacement, but you can create and apply scrubbers for cleaning the
subject or obfuscate sensible data

The basic RegexpScrubber is a generic RegexpScrubber that searches for a regex target and replaces it
*/
type RegexpScrubber struct {
	baseScrubber
}

func NewScrubber(pattern, replacement string, opts ...ScrubberOption) RegexpScrubber {
	s := baseScrubber{
		target:      pattern,
		replacement: replacement,
		context:     "%s",
	}

	for _, opt := range opts {
		opt(&s)
	}
	return RegexpScrubber{baseScrubber: s}
}

func (b RegexpScrubber) Clean(subject string) string {
	re := regexp.MustCompile(fmt.Sprintf(b.context, b.target))
	return re.ReplaceAllString(subject, fmt.Sprintf(b.context, b.replacement))
}

type PathScrubber struct {
	baseScrubber
}

func NewPathScrubber(pattern, replacement string, opts ...ScrubberOption) PathScrubber {
	s := baseScrubber{
		target:      "",
		replacement: replacement,
		context:     pattern,
	}

	for _, opt := range opts {
		opt(&s)
	}
	return PathScrubber{baseScrubber: s}
}

func (s PathScrubber) Clean(subject string) string {
	r := gjson.Get(subject, s.context)
	if !r.Exists() {
		return subject
	}
	scrubbed, err := sjson.Set(subject, s.context, s.replacement)
	if err != nil {
		return subject
	}
	return scrubbed
}

/*

## Custom Scrubbers

*/

/*
CreditCard obfuscates credit card numbers

	ccScrubber := golden.CreditCard()
*/
func CreditCard(opts ...ScrubberOption) RegexpScrubber {
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

	fixedULIDScrubber := golden.ULID(golden.Replacement("01HNB10NSJS26X2RTERPZTM0KB"))
	anotherPlaceHolder := golden.ULID(golden.Replacement("[ULID here]]"))
*/
func ULID(opts ...ScrubberOption) RegexpScrubber {
	return NewScrubber(
		"[0-9A-Za-z]{26}",
		"<ULID>",
		opts...,
	)
}

/*

## RegexpScrubber options

*/

type ScrubberOption func(s *baseScrubber)

/*
Replacement define a custom replacement for any specialized RegexpScrubber

	ulidScrubber := golden.ULID(golden.Replacement("[ULID comes here]"))
*/
func Replacement(r string) ScrubberOption {
	return func(s *baseScrubber) {
		s.replacement = r
	}
}

/*
Format will help to limit the scrubbing to string that matches the context. Use
the placeholder %s to indicate the part that you want to be scrubbed

	ccScrubber := golden.CreditCard(golden.Format("Credit Card: %s"))

will apply the CreditCard obfuscation only if it finds a Credit Card Number
after a "Credit Card: " string. This way you can avoid scrubbing parts of the
subject that you don't want to touch.
*/
func Format(f string) ScrubberOption {
	return func(s *baseScrubber) {
		s.context = f
	}
}
