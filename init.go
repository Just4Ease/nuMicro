package nuMicro

import (
	"fmt"
	"log"

	"github.com/Just4Ease/nuMicro/broker"
)

type mountServices func(serviceName string)

func Init(serviceName string, eventsHandler mountServices) {
	if err := broker.Connect(); err != nil {
		log.Fatal(err, " Failed to start broker. Mayday!, Mayday! Call the NATS officer, ensure all is well!")
	}
	forever := make(chan bool)
	if eventsHandler != nil {
		go eventsHandler(serviceName)
	}
	if eventsHandler == nil {
		fmt.Println("Event handler has not been set, no events will be listened for.")
	}

	log.Println(fmt.Sprintf("Running nuMicro as : %s", serviceName))
	<-forever
}
