package main

import (
	"fmt"

	"github.com/Just4Ease/nuMicro/broker"
)

func main() {
	_ = broker.Connect()

	_, _ = broker.Respond("ExampleSVC.sample", func(event broker.Event) interface{} {

		result := make(map[string]string)
		result["username"] = "just4ease"
		result["email"] = "justicenefe@gmail.com"
		result["phone"] = "+2347056031137"
		result["website"] = "https://justicenefe.com"

		return result
	}, nil)



	i, err := broker.Request("ExampleSVC.sample", nil, nil)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(i)
}
