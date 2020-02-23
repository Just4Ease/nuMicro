package json

import (
	"encoding/json"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

type Marshaller struct{}

func (j Marshaller) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (j Marshaller) Unmarshal(d []byte, v interface{}) error {
	if pb, ok := v.(proto.Message); ok {
		return jsonpb.UnmarshalString(string(d), pb)
	}
	return json.Unmarshal(d, v)
}

func (j Marshaller) String() string {
	return "json"
}
