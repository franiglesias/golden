package golden

import (
	"bytes"
	"encoding/json"
	"sort"
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
	var output string
	if _, ok := subject.(string); ok {
		output = subject.(string)
	} else {
		rawSubject, err := json.MarshalIndent(subject, prefix, indent)
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
	result, err := sortJSONFields(str)
	if err != nil {
		return prettyJSON.String()
	}
	return result
}

func sortJSONFields(jsonStr string) (string, error) {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return "", err
	}

	sortMapFields(data)

	prettyJSON, err := json.MarshalIndent(data, prefix, indent)
	if err != nil {
		return "", err
	}

	return string(prettyJSON), nil
}

func sortMapFields(m map[string]interface{}) {
	for _, v := range m {
		if nestedMap, ok := v.(map[string]interface{}); ok {
			sortMapFields(nestedMap)
		}
	}

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	sortedMap := make(map[string]interface{})
	for _, k := range keys {
		sortedMap[k] = m[k]
	}

	for k := range m {
		delete(m, k)
	}
	for k, v := range sortedMap {
		m[k] = v
	}
}
