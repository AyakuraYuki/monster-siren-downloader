package json

import (
	"encoding/json"

	"github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
)

var JSON = jsoniter.ConfigCompatibleWithStandardLibrary

func RegisterFuzzyDecoders() {
	extra.RegisterFuzzyDecoders()
}

func Stringify(v any) string {
	if v == nil {
		return "null"
	}
	raw, _ := JSON.MarshalToString(v)
	return raw
}

func Prettify(v any) string {
	if v == nil {
		return "null"
	}
	bs, _ := json.MarshalIndent(v, "", "    ")
	return string(bs)
}
