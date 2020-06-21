package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

func main() {
	var nc *nats.Conn
	var err error
	nc, err = nats.Connect("nats://49.235.146.124:4222", nats.MaxReconnects(6), nats.ReconnectHandler(func(conn *nats.Conn) {
		fmt.Println("重新连接")
		nc = conn
	}), nats.DisconnectErrHandler(func(conn *nats.Conn, e error) {
		fmt.Println("连接断线")
		fmt.Println(e.Error())
	}))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	sc, err := stan.Connect("c1", "server-host-1",
		stan.NatsConn(nc))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer sc.Close()

	if err := sc.Publish("foo", []byte("hello world")); err != nil {
		fmt.Println(err.Error())
		return
	}

	if err := sc.Publish("bar", []byte("hello world!!!")); err != nil {
		fmt.Println(err.Error())
		return
	}
}
