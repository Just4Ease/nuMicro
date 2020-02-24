package nuMicro

import (
	"fmt"
	"log"

	"github.com/Just4Ease/nuMicro/broker"
)

type mountServices func(serviceName string)

func Init(serviceName string, eventsHandler mountServices, actionsHandler mountServices) {
	if err := broker.Connect(); err != nil {
		log.Fatal(err, " Failed to start broker. Mayday!, Mayday! Call the NATS officer, ensure all is well!")
	}
	forever := make(chan bool)
	if eventsHandler != nil {
		go eventsHandler(serviceName)
	}
	if actionsHandler != nil {
		go actionsHandler(serviceName)
	}
	if eventsHandler == nil {
		fmt.Println("Event handler has not been set, no events will be listened for.")
	}
	if actionsHandler == nil {
		fmt.Println("Action handler has not been set, no actions will be called.")
	}

	log.Println(fmt.Sprintf("Running nuMicro as : %s", serviceName))
	<-forever
}
