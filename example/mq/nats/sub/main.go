package main

import (
    "fmt"
    "github.com/nats-io/nats.go"
)

func main() {
    url := "nats://49.235.146.124:4222"
    opts := []nats.Option{nats.Name("publisher")}
    opts = append(opts, nats.DisconnectErrHandler(func(conn *nats.Conn, e error) {
        // 重连
        fmt.Println("dis connect:"+e.Error())
    }))
    conn,err := nats.Connect(url, opts...)
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    defer conn.Close()



    subject := "money"
    conn.Subscribe(subject, func(msg *nats.Msg) {
        fmt.Println(string(msg.Data))
    })

    select {}
}
