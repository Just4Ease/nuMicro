package main

import (
	"fmt"

	"github.com/Just4Ease/nuMicro"
	"github.com/Just4Ease/nuMicro/broker"
)

func main() {
	i, err := broker.Request("ExampleSVC.sample", nil, nil)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(i)

	serviceName := "ExampleSVC"
	nuMicro.Init(serviceName, nil, Action, broker.Addrs("127.0.0.1:4222"))
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
