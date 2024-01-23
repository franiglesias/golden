package golden

import "regexp"

type Scrubber interface {
	Clean(subject string) string
}

type BasicRegexpScrubber struct {
	pattern     string
	replacement string
}

func NewRegexpScrubber(pattern, replacement string) BasicRegexpScrubber {
	return BasicRegexpScrubber{
		pattern:     pattern,
		replacement: replacement,
	}
}

func (b BasicRegexpScrubber) Clean(subject string) string {
	re := regexp.MustCompile(b.pattern)
	return re.ReplaceAllString(subject, b.replacement)
}
