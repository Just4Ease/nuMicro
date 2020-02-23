package nuMicro

import (
	"github.com/Just4Ease/nuMicro/broker"
	"github.com/Just4Ease/nuMicro/utils/log"
)

type mountServices func(serviceName string)

func Init(serviceName string, f mountServices) {
	if err := broker.Connect(); err != nil {
		log.Fatal(err, " Failed to start broker. Ensure all is well!")
	}
	if err := broker.Connect(); err != nil {
		log.Fatal(err)
	}

	forever := make(chan bool)
	f(serviceName)
	<-forever
}
