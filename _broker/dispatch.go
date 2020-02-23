package _broker

import (
	"fmt"

	"github.com/vmihailenco/msgpack"
)

func (b *Broker) Dispatch(to string, action string, message Message) {
	defer b.Flush()
	body, _ := msgpack.Marshal(message)
	subject := fmt.Sprintf("%s.%s", to, action)
	if err := b.connection.Publish(subject, body); err != nil {
		fmt.Println(err.Error(), " Error Dispatching message")
	}
}
