package msgpack

import (
	"github.com/vmihailenco/msgpack"
)

type Marshaller struct{}

func (m Marshaller) Marshal(v interface{}) ([]byte, error) {
	return msgpack.Marshal(v)
}

func (m Marshaller) Unmarshal(d []byte, v interface{}) error {
	return msgpack.Unmarshal(d, v)
}

func (m Marshaller) String() string {
	return "msgpack"
}
