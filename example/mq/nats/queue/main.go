package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
)

func main() {
	url := "nats://49.235.146.124:4222"
	opts := []nats.Option{nats.Name("publisher")}
	conn, err := nats.Connect(url, opts...)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()

	conn.QueueSubscribe("money", "people", func(msg *nats.Msg) {
		fmt.Println(string(msg.Data))
	})

	select {}
}
