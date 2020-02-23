package _broker

import (
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/vmihailenco/msgpack"
)

func New(service string) *Broker {
	connection, err := nats.Connect(nats.DefaultURL)
	failOnError(err, "Failed to connect to NATS")
	return &Broker{
		ServiceName: service,
		connection:  connection,
		registry:    make(Registry, 0), // Initialize an empty registry first.
	}
}

func (b *Broker) RegisterAction(name string, fn ActionHandler) {
	b.Lock()
	defer b.Unlock()
	b.registry[name] = fn
	forever := make(chan bool)
	go b.Listen(&name)
	<-forever
}

func (b *Broker) ListServices() {
	b.Lock()
	defer b.Unlock()
	fmt.Println(b.registry)
}

func (b *Broker) Discover(caller string) []string {
	b.Lock()
	defer b.Unlock()

	actions := make([]string, 0)

	for name := range b.registry {
		actions = append(actions, name)
	}

	return actions
}

func (b *Broker) Listen(name *string) {
	subject := fmt.Sprintf("%s.%s", b.ServiceName, *name)
	if sub, err := b.connection.Subscribe(subject, func(msg *nats.Msg) {
		if b.registry[*name] == nil {
			response := _utils.Result{
				Success: false,
				Message: "Sorry, you have an invalid method call",
				Error:   _utils.ActionNotAllowed,
			}
			message, _ := msgpack.Marshal(response)
			if err := msg.Respond(message); err != nil {
				failOnError(err, "Error responding to "+*name)
			}
		}
		action := b.registry[*name]
		payload := make(map[string]interface{})
		_ = _utils.PackRaw(msg.Data, &payload)
		result := action(payload)
		message, e := msgpack.Marshal(result)
		if e != nil {
			fmt.Println(e, " Error Encoding.")
		}
		if err := msg.Respond(message); err != nil {
			failOnError(err, "Error responding to "+*name)
		}
	}); err != nil {
		fmt.Println(err.Error(), " Error subscribing to chan.")
	} else {
		fmt.Println(sub, " Sub Status.")
	}
}

//func HandleMessage() {
//
//}

//func ExecuteAction(b Broker, message Message) {
//	to := message.From
//	if b.registry[message.Action] == nil {
//		message.To = to
//		message.From = b.ServiceName
//		message.Payload = utils.Result{
//			Success: false,
//			Message: "Sorry, you have tried an invalid action",
//			Error:   utils.ActionNotAllowed,
//		}
//		go b.Dispatch(to, "", message)
//	}
//
//	action := b.registry[message.Action]
//	input := make(map[string]interface{})
//	_ = utils.Pack(message.Payload, &input)
//	result := action(input)
//	output := Message{
//		From:     b.ServiceName,
//		Action:   message.ReplyWith,
//		Payload:  result,
//		TimeSent: time.Now(),
//	}
//	go b.Dispatch(to, "", output)
//}

func (b *Broker) Request(to string, message []byte, replyTo *string) {
	defer b.Flush()
	if err := b.connection.PublishRequest(to, *replyTo, message); err != nil {
		fmt.Println(err.Error(), " Error Requesting message")
	}
}

func (b *Broker) Flush() {
	_ = b.connection.Flush()
}

func (b *Broker) OnResponse(response string) {
	defer b.Flush()
	_, _ = b.connection.Subscribe(response, func(msg *nats.Msg) {
		//hander()
		r := make(map[string]interface{})
		_ = _utils.PackRaw(msg.Data, &r)

		fmt.Println(r, " REs.")
		//fmt.Println(r, " Response In.")
	})
}
