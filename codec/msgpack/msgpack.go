package msgpack

import (
	"encoding/json"
	"io"

	"github.com/Just4Ease/nuMicro/codec"
)

type Codec struct {
	Conn    io.ReadWriteCloser
	Encoder *json.Encoder
	Decoder *json.Decoder
}

func (c *Codec) ReadHeader(m *codec.Message, t codec.MessageType) error {
	return nil
}

func (c *Codec) ReadBody(b interface{}) error {
	if b == nil {
		return nil
	}
	return c.Decoder.Decode(b)
}

func (c *Codec) Write(m *codec.Message, b interface{}) error {
	if b == nil {
		return nil
	}
	return c.Encoder.Encode(b)
}

func (c *Codec) Close() error {
	return c.Conn.Close()
}

func (c *Codec) String() string {
	return "msgpack"
}

func NewCodec(c io.ReadWriteCloser) codec.Codec {
	//return &Codec{
	//	Conn:    c,
	//	Decoder: msgpack.NewDecoder(c),
	//	Encoder: msgpack.NewDecoder(c),
	//}
	return nil
}
