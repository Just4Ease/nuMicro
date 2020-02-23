package nuMicro

import (
	"fmt"
	"log"

	"github.com/Just4Ease/nuMicro/broker"
)

type mountServices func(serviceName string)

func Init(serviceName string, f mountServices) {
	if err := broker.Connect(); err != nil {
		log.Fatal(err, " Failed to start broker. Ensure all is well!")
	}
	forever := make(chan bool)
	f(serviceName)
	// TODO: Run service pings, health checks here and discoveries etc..
	_log_ := fmt.Sprintf("Running nuMicro as : %s", serviceName)
	fmt.Println(_log_)
	<-forever
}
