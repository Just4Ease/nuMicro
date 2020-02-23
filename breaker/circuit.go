package breaker

import (
	"time"

	breaker "github.com/sony/gobreaker"
)

var cb *breaker.CircuitBreaker

type Counts breaker.Counts

type State breaker.State

type Settings struct {
	Name          string
	MaxRequests   uint32
	Interval      time.Duration
	Timeout       time.Duration
	ReadyToTrip   func(counts Counts) bool
	OnStateChange func(name string, from State, to State)
}

func NewCircuitBreaker(st *breaker.Settings) *breaker.CircuitBreaker {
	return &breaker.CircuitBreaker{}
}
