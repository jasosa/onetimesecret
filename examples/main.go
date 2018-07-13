package main

import (
	"fmt"

	"github.com/jasosa/onetimesecret"
)

func main() {
	userEmail := "Your user email here"
	apiToken := "Your api token here"

	client := onetimesecret.NewClient(userEmail, apiToken, "https://onetimesecret.com/api/v1")
	secret, value, err := client.Generate(3600)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(secret)
		fmt.Println(value)
	}
}
