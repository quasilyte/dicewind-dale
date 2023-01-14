package jsonconfig

import (
	"encoding/json"
	"regexp"
)

var (
	trailingCommaRE = regexp.MustCompile(`,\s*([}\]])`)
	commentsRE      = regexp.MustCompile(`\s*//[^\n]*`)
)

func Unmarshal(data []byte, v any) error {
	data = trailingCommaRE.ReplaceAll(data, []byte("$1"))
	data = commentsRE.ReplaceAll(data, []byte(""))
	return json.Unmarshal(data, v)
}
