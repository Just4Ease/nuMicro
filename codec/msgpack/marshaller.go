package msgpack

import (
	"github.com/vmihailenco/msgpack"
)

type Marshaller struct{}

func (m Marshaller) Marshal(v interface{}) ([]byte, error) {
	return msgpack.Marshal(v)
}

func (m Marshaller) Unmarshal(d []byte, v interface{}) error {
	b, err := msgpack.Marshal(&d)
	if err != nil {
		return err
	}

	if err := msgpack.Unmarshal(b, &v); err != nil {
		return err
	}

	return nil
}

func (m Marshaller) String() string {
	return "msgpack"
}
