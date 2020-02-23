// Package mdns provides a multicast dns registry
package mdns

import (
	"context"

	"github.com/Just4Ease/nuMicro/registry"
)

// NewRegistry returns a new mdns registry
func NewRegistry(opts ...__registry.Option) __registry.Registry {
	return __registry.NewRegistry(opts...)
}

// Domain sets the mdnsDomain
func Domain(d string) __registry.Option {
	return func(o *__registry.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, "mdns.domain", d)
	}
}
