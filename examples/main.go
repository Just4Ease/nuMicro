package main

import (
	"fmt"
	"time"

	"github.com/Just4Ease/nuMicro"
	"github.com/Just4Ease/nuMicro/broker"
)

func main() {
	go func() {
		time.Sleep(5 * time.Second)
		i, err := broker.Request("ExampleSVC.sample", nil, nil)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(i)
	}()
	serviceName := "ExampleSVC"
	nuMicro.Init(serviceName, nil, Action, broker.Addrs("nats://demo.nats.io"))
}

func Action(serviceName string) {
	_, _ = broker.Respond("ExampleSVC.sample", func(event broker.RequestEvent) interface{} {

		result := make(map[string]string)
		result["username"] = "just4ease"
		result["email"] = "justicenefe@gmail.com"
		result["phone"] = "+2347056031137"
		result["website"] = "https://justicenefe.com"

		return result
	})
}
