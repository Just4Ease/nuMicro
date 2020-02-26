package broker

import (
	"github.com/Just4Ease/nuMicro/utils/log"
)

type serviceSub struct {
	channel string
	queue   string
	handler Handler
	stream  SubscribeService
	closed  chan bool
	options SubscribeOptions
}

type serviceEvent struct {
	channel string
	message *Message
}

func (s *serviceEvent) Channel() string {
	return s.channel
}

func (s *serviceEvent) Message() *Message {
	return s.message
}

func (s *serviceEvent) Ack() error {
	return nil
}

func (s *serviceSub) isClosed() bool {
	select {
	case <-s.closed:
		return true
	default:
		return false
	}
}

func (s *serviceSub) run() error {
	exit := make(chan bool)
	go func() {
		select {
		case <-exit:
		case <-s.closed:
		}

		// close the stream
		_ = s.stream.Close()
	}()

	for {
		// TODO: do not fail silently
		msg, err := s.stream.Recv()
		if err != nil {
			log.Debugf("Streaming error for subscription to topic %s: %v", s.Channel(), err)

			// close the exit channel
			close(exit)

			// don't return an error if we unsubscribed
			if s.isClosed() {
				return nil
			}

			// return stream error
			return err
		}

		// TODO: handle error
		_, _ = s.handler(&serviceEvent{
			channel: s.channel,
			message: &Message{
				Header: msg.Header,
				Body:   msg.Body,
			},
		})
	}
}

func (s *serviceSub) Options() SubscribeOptions {
	return s.options
}

func (s *serviceSub) Channel() string {
	return s.channel
}

func (s *serviceSub) Unsubscribe() error {
	select {
	case <-s.closed:
		return nil
	default:
		close(s.closed)
	}
	return nil
}
