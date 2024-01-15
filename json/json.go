package json

import jsoniter "github.com/json-iterator/go"

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func Unmarshal(data []byte, val interface{}) error {
	return json.Unmarshal(data, val)
}

func Marshal(val interface{}) ([]byte, error) {
	return json.Marshal(val)
}
