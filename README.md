nuMicro
=========

####A micro services framework for event driven distributed systems.


[![GoDoc](https://godoc.org/github.com/just4ease/nuMicro?status.svg)](https://godoc.org/github.com/just4ease/nuMicro)


https://pkg.go.dev/github.com/Just4Ease/nuMicro

Installation
------------

```
go get github.com/Just4Ease/nuMicro
```

Usage
-----

The Init `nuMicro.Init` is a an entry point into starting the Microservice.


```go
package main

import (

"fmt"


"github.com/Just4Ease/nuMicro"
"github.com/Just4Ease/nuMicro/broker"
) 
var serviceName string

func main()  {
  serviceName = "UsersSVC" // The name of your microservice here.
  nuMicro.Init(serviceName, EventsHandle, ActionsHandle)
}

func BuildChannelName(channelName string) string {
    s := fmt.Sprintf("%s.%s", serviceName, channelName)
    return s
}


func EventsHandle(serviceName string)  {
 // Handle your events here for pub/sub pattern
}

func ActionsHandle(serviceName string)  {
 // Handle your actions here for pub/sub pattern
_, _ = broker.Respond(BuildChannelName("Contact"), func(event broker.Event) interface{} {
        // Here, this could be a function that processes something to respond when it is being called on the contact channel.
		result := make(map[string]string)
		result["username"] = "just4ease"
		result["email"] = "justicenefe@gmail.com"
		result["phone"] = "+2347056031137"
		result["website"] = "https://justicenefe.com"

		return result
	}, nil)
}
```
#NOTE: 
----
- nuMicro uses NATS.io as it's broker but you can extend and configure your message broker

TODO:
----
- Mount service discovery

- Mount the probe, the probe is used to get an auto-documentation for each action call in the microservice.
It would return the input argument and output response sample for each action call so that the "client" can send in the right parameters.

- Use `msgpack` for internal service/client to service/client encoding and decoding of messages. It's faster than JSON, BSON and Gob, gateways can use whatever they desire.   

- Enable auto scaling ( Still an experimental feature right now. )

- Embed internal NATS server for fallback if a cloud nats goes down, and republish it's public NATS address so that it can be discoverable within 10 - 100ms without noticeable downtime. 

- Enable the circuit breaker
- Enable an extra layer for idempotence


License
-------

The MIT License (MIT)

See [LICENSE](https://github.com/just4ease/nuMicro/blob/master/LICENSE) for details.


[repo-url]: https://github.com/Just4Ease/nuMicro
