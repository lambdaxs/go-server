package main

import (
    "fmt"
    "github.com/isayme/go-amqp-reconnect/rabbitmq"
)

func main() {
    rabbitmq.Debug = true

    conn,err := rabbitmq.Dial("amqp://admin:123456@49.235.146.124:5672/")
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

    //声明队列
    q,err := channel.QueueDeclare("hello", true, false, false, false, nil)
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    //批量接收消息
    msgs,err := channel.Consume(q.Name, "c1", false, false, false, false, nil)
    if err != nil {
        fmt.Println(err.Error())
        return
    }


    go func() {
        for msg := range msgs {
            fmt.Println("接受消息"+ string(msg.Body))
            msg.Ack(false)
        }
    }()

    select {}
}
