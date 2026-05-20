package main

import (
	"booking-service/cmd/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		panic(err)
	}
}
