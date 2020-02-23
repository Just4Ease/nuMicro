package noop

import (
	"github.com/Just4Ease/nuMicro/debug/log"
)

type noop struct{}

func (n *noop) Read(...log.ReadOption) ([]log.Record, error) {
	return nil, nil
}

func (n *noop) Write(log.Record) error {
	return nil
}

func (n *noop) Stream() (log.Stream, error) {
	return nil, nil
}

func NewLog(opts ...log.Option) log.Log {
	return new(noop)
}
