// Package broker is an interface used for asynchronous messaging
package broker

// Broker is an interface used for asynchronous messaging.
type Broker interface {
	Init(...Option) error
	Options() Options
	Address() string
	Connect() error
	Disconnect() error
	Publish(channel string, m *Message, opts ...PublishOption) error
	Request(channel string, m *Message, opts ...PublishOption) (interface{}, error)
	Respond(channel string, h ActionHandle, opts ...SubscribeOption) (Subscriber, error)
	Subscribe(channel string, h Handler, opts ...SubscribeOption) (Subscriber, error)
	String() string
}

// Handler is used to process messages via a subscription of a channel.
// The handler is passed a publication interface which contains the
// message and optional Ack method to acknowledge receipt of the message.
type Handler func(Event) error
type ActionHandle func(RequestEvent) interface{}

type Message struct {
	Header map[string]string
	Body   []byte
}

/**
 * For requests
 */
type RequestInput struct {
	Headers map[string]string
	Body    interface{}
	Params  map[string]interface{}
	Query   map[string]interface{}
}

type RequestEvent interface {
	Channel() string
	Message() *RequestInput
	Ack() error
}

type SubscribeService interface {
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Recv() (*Message, error)
}

// Event is given to a subscription handler for processing
type Event interface {
	Channel() string
	Message() *Message
	Ack() error
}

// Subscriber is a convenience return type for the Subscribe method
type Subscriber interface {
	Options() SubscribeOptions
	Channel() string
	Unsubscribe() error
}

var (
	DefaultBroker = NewBroker()
)

func Init(opts ...Option) error {
	return DefaultBroker.Init(opts...)
}

func Connect() error {
	return DefaultBroker.Connect()
}

func Disconnect() error {
	return DefaultBroker.Disconnect()
}

func Publish(channel string, msg *Message, opts ...PublishOption) error {
	return DefaultBroker.Publish(channel, msg, opts...)
}

func Subscribe(channel string, handler Handler, opts ...SubscribeOption) (Subscriber, error) {
	return DefaultBroker.Subscribe(channel, handler, opts...)
}

func Respond(channel string, handler ActionHandle) (Subscriber, error) {
	return DefaultBroker.Respond(channel, handler)
}
func Request(channel string, msg *Message, opts ...PublishOption) (interface{}, error) {
	return DefaultBroker.Request(channel, msg, opts...)
}

func String() string {
	return DefaultBroker.String()
}
