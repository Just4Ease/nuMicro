package json

import (
	"encoding/json"
)

type Marshaller struct{}

func (j Marshaller) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (j Marshaller) Unmarshal(d []byte, v interface{}) error {
	return json.Unmarshal(d, v)
}

func (j Marshaller) String() string {
	return "json"
}
