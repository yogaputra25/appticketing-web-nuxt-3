package model

import "encoding/json"

func jsonMarshal(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}

func jsonUnmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
