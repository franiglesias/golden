package golden

import (
	"encoding/json"
	"strings"
)

/*
JsonNormalizer is the default subject normalization tool for golden. You
can introduce custom normalizers. Take into account the trimming of the output
to ensure that the normalization process doesn't introduce undesirable leading
or trailing characters that could lead to irrelevant differences between the
subject and the snapshot.
*/
type JsonNormalizer struct {
}

func (n JsonNormalizer) Normalize(subject any) (string, error) {
	rawSubject, err := json.MarshalIndent(subject, "", "  ")
	if err != nil {
		return "", err
	}
	trimmed := strings.Trim(string(rawSubject), `" \n`)
	return trimmed, nil
}

//TODO: consider support for scrubbers here
