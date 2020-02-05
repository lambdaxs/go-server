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

    /**声明队列，
    队列名
    是否持久化，
    是否排外，只允许该channel访问该对了，用于一个队列只能有一个消费者来消费
    是否自动删除 消费完删除
    */
    q,err := channel.QueueDeclare("hello", true, false, false, false, nil)
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    /**发布消息
    交换机
    队列名
    路由
     */
    if err := channel.Publish("", q.Name, false, false, amqp.Publishing{
        ContentType:     "text/plain",
        Body:            []byte("Hello world!!"),
    });err != nil {
       fmt.Println(err.Error())
        return
    }
    fmt.Println("send success")
}
