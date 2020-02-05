package main

import (
    "fmt"
    "github.com/streadway/amqp"
)

func main() {
    conn,err := amqp.Dial("amqp://admin:123456@49.235.146.124:5672/")
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    defer conn.Close()

    channel,err := conn.Channel()
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    defer channel.Close()

}
