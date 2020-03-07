package broker

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Just4Ease/nuMicro/codec/msgpack"
	"github.com/gofrs/uuid"
	"github.com/nats-io/nats.go"
)

type natsBroker struct {
	sync.Once
	sync.RWMutex

	// indicate if we're connected
	connected bool

	// address to bind routes to
	addresses []string
	// servers for the client
	servers []string

	// client connection and nats opts
	conn  *nats.Conn
	opts  Options
	nopts nats.Options

	// should we drain the connection
	drain   bool
	closeCh chan error

	//// embedded server
	//server *server.Server
	// configure to use local server
	local bool
	// server exit channel
	exit chan bool
}

type subscriber struct {
	s    *nats.Subscription
	opts SubscribeOptions
}

type publication struct {
	c string
	m *Message
}

type req struct {
	c string
	m *RequestInput
}

func (r req) Channel() string {
	return r.c
}
func (r req) Message() *RequestInput {
	return r.m
}

func (r req) Ack() error {
	// Our server auto acks.
	return nil
}

func (p *publication) Channel() string {
	return p.c
}

func (p *publication) Message() *Message {
	return p.m
}

func (p *publication) Ack() error {
	// nats does not support acking
	return nil
}

func (s *subscriber) Options() SubscribeOptions {
	return s.opts
}

func (s *subscriber) Channel() string {
	return s.s.Subject
}

func (s *subscriber) Unsubscribe() error {
	return s.s.Unsubscribe()
}

//func (n *natsBroker) Address() string {
//	n.RLock()
//	defer n.RUnlock()
//
//	if n.server != nil {
//		return n.server.ClusterAddr().String()
//	}
//
//	if n.conn != nil && n.conn.IsConnected() {
//		return n.conn.ConnectedUrl()
//	}
//
//	if len(n.addrs) > 0 {
//		return n.addrs[0]
//	}
//
//	return "127.0.0.1:-1"
//}

func (n *natsBroker) setAddresses(addresses []string) []string {
	//nolint:prealloc
	var connectionAddresses []string
	for _, address := range addresses {
		if len(address) == 0 {
			continue
		}
		if !strings.HasPrefix(address, "nats://") {
			address = "nats://" + address
		}
		connectionAddresses = append(connectionAddresses, address)
	}
	// if there's no address and we weren't told to
	// embed a local server then use the default url
	if len(connectionAddresses) == 0 && !n.local {
		connectionAddresses = []string{nats.DefaultURL}
	}
	return connectionAddresses
}

func (n *natsBroker) Connect() error {
	n.Lock()
	defer n.Unlock()

	if !n.connected {
		// create exit chan
		n.exit = make(chan bool)

		// start the local server

		// set to connected
	}

	status := nats.CLOSED
	if n.conn != nil {
		status = n.conn.Status()
	}

	switch status {
	case nats.CONNECTED, nats.RECONNECTING, nats.CONNECTING:
		return nil
	default: // DISCONNECTED or CLOSED or DRAINING
		opts := n.nopts
		opts.DrainTimeout = 1 * time.Second
		opts.AsyncErrorCB = n.onAsyncError
		opts.DisconnectedErrCB = n.onDisconnectedError
		opts.ClosedCB = n.onClose
		//opts.Servers = n.servers
		opts.Secure = n.opts.Secure
		opts.TLSConfig = n.opts.TLSConfig

		// secure might not be set
		if n.opts.TLSConfig != nil {
			opts.Secure = true
		}

		c, err := opts.Connect()
		if err != nil {
			return err
		}
		n.conn = c

		n.connected = true

		return nil
	}
}

func (n *natsBroker) Disconnect() error {
	n.RLock()
	defer n.RUnlock()

	if !n.connected {
		return nil
	}

	// drain the connection if specified
	if n.drain {
		_ = n.conn.Drain()
	}

	// close the client connection
	n.conn.Close()

	// shutdown the local server
	// and deregister
	//if n.server != nil {
	//	select {
	//	case <-n.exit:
	//	default:
	//		close(n.exit)
	//	}
	//}

	// set not connected
	n.connected = false

	return nil
}

func (n *natsBroker) Init(opts ...Option) error {
	n.setOption(opts...)
	return nil
}

func (n *natsBroker) Options() Options {
	return n.opts
}

func (n *natsBroker) Publish(channel string, msg *Message, opts ...PublishOption) error {
	b, err := n.opts.Codec.Marshal(msg)
	if err != nil {
		return err
	}
	n.RLock()
	defer n.RUnlock()
	return n.conn.Publish(channel, b)
}

func (n *natsBroker) Request(channel string, msg *Message, opts ...PublishOption) (interface{}, error) {
	id, _ := uuid.NewV4()
	replyAlias := fmt.Sprintf("%s", id)
	var result interface{}
	wg := sync.WaitGroup{}
	b, err := n.opts.Codec.Marshal(msg)
	if err != nil {
		return nil, err
	}
	n.RLock()
	defer n.RUnlock()
	wg.Add(1)
	go func(r *interface{}) {
		_, _ = n.conn.Subscribe(replyAlias, func(msg *nats.Msg) {
			defer wg.Done()
			_ = n.opts.Codec.Unmarshal(msg.Data, &r)
		})
	}(&result)
	_ = n.conn.PublishRequest(channel, replyAlias, b)
	wg.Wait()
	return result, nil
}

func (n *natsBroker) Subscribe(channel string, handler Handler, opts ...SubscribeOption) (Subscriber, error) {
	if n.conn == nil {
		return nil, errors.New("not connected")
	}

	opt := SubscribeOptions{
		AutoAck: true,
		Context: context.Background(),
	}

	for _, o := range opts {
		if o != nil {
			o(&opt)
		}
	}

	fn := func(msg *nats.Msg) {
		var m Message
		if err := n.opts.Codec.Unmarshal(msg.Data, &m); err != nil {
			return
		}
		_ = handler(&publication{m: &m, c: msg.Subject})
	}

	var sub *nats.Subscription
	var err error

	n.RLock()
	if len(opt.Queue) > 0 {
		sub, err = n.conn.QueueSubscribe(channel, opt.Queue, fn)
	} else {
		sub, err = n.conn.Subscribe(channel, fn)
	}
	n.RUnlock()
	if err != nil {
		return nil, err
	}
	//&subscriber{s: sub, opts: opt}
	return &subscriber{s: sub, opts: opt}, nil
}

func (n *natsBroker) Respond(channel string, handler ActionHandle) (Subscriber, error) {
	if n.conn == nil {
		return nil, errors.New("not connected")
	}

	opt := SubscribeOptions{
		AutoAck: true,
		Context: context.Background(),
	}

	fn := func(msg *nats.Msg) {
		var m RequestInput
		if err := n.opts.Codec.Unmarshal(msg.Data, &m); err != nil {
			return
		}
		i := handler(&req{m: &m, c: msg.Subject})
		out, _ := n.opts.Codec.Marshal(i)
		_ = msg.Respond(out)
	}

	var sub *nats.Subscription
	var err error

	n.RLock()
	if len(opt.Queue) > 0 {
		sub, err = n.conn.QueueSubscribe(channel, opt.Queue, fn)
	} else {
		sub, err = n.conn.Subscribe(channel, fn)
	}
	n.RUnlock()
	if err != nil {
		return nil, err
	}
	//&subscriber{s: sub, opts: opt}
	return &subscriber{s: sub, opts: opt}, nil
}

func (n *natsBroker) String() string {
	return "nats"
}

func (n *natsBroker) setOption(opts ...Option) {
	for _, o := range opts {
		o(&n.opts)
	}

	n.Once.Do(func() {
		n.nopts = nats.GetDefaultOptions()
	})

	// local embedded server
	n.local = true
	// set to drain
	n.drain = true

	if !n.opts.Secure {
		n.opts.Secure = n.nopts.Secure
	}

	if n.opts.TLSConfig == nil {
		n.opts.TLSConfig = n.nopts.TLSConfig
	}

	n.addresses = n.setAddresses(n.opts.Addresses)
}

func (n *natsBroker) onClose(conn *nats.Conn) {
	n.closeCh <- nil
}

func (n *natsBroker) onDisconnectedError(conn *nats.Conn, err error) {
	n.closeCh <- err
}

func (n *natsBroker) onAsyncError(conn *nats.Conn, sub *nats.Subscription, err error) {
	// There are kinds of different async error nats might callback, but we are interested
	// in ErrDrainTimeout only here.
	if err == nats.ErrDrainTimeout {
		n.closeCh <- err
	}
}

func NewBroker(opts ...Option) *natsBroker {
	options := Options{
		// Default codec
		Codec:   msgpack.Marshaller{},
		Context: context.Background(),
	}

	n := &natsBroker{
		opts:    options,
		closeCh: make(chan error),
	}
	n.setOption(opts...)

	return n
}
