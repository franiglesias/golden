package golden

import (
	"bytes"
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

const indent = "  "
const prefix = ""

func (n JsonNormalizer) Normalize(subject any) (string, error) {
	var rawSubject []byte
	var output string
	if _, ok := subject.(string); ok {
		output = subject.(string)
	} else {
		var err error
		rawSubject, err = json.MarshalIndent(subject, prefix, indent)
		output = string(rawSubject)
		if err != nil {
			return "", err
		}
	}

	output = strings.Trim(strings.Trim(output, "\n"), `" `)
	return prettyPrint(output), nil
}

/*
prettyPrint prettify valid json if detected it so the snapshot is more readable
to humans. If not, return the string as is.
*/
func prettyPrint(str string) string {
	var prettyJSON bytes.Buffer
	// If not valid json return as is
	if err := json.Indent(&prettyJSON, []byte(str), prefix, indent); err != nil {
		return str
	}
	// Return pretty json
	return prettyJSON.String()
}
