package bytes

import (
	"errors"
)

type Marshaller struct{}

type Message struct {
	Header map[string]string
	Body   []byte
}

func (n Marshaller) Marshal(v interface{}) ([]byte, error) {
	switch ve := v.(type) {
	case *[]byte:
		return *ve, nil
	case []byte:
		return ve, nil
	case *Message:
		return ve.Body, nil
	}
	return nil, errors.New("invalid message")
}

func (n Marshaller) Unmarshal(d []byte, v interface{}) error {
	switch ve := v.(type) {
	case *[]byte:
		*ve = d
	case *Message:
		ve.Body = d
	}
	return errors.New("invalid message")
}

func (n Marshaller) String() string {
	return "bytes"
}
