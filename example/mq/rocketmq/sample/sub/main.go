package main

import (
	"context"
	"fmt"
	"github.com/lambdaxs/rocketmq-client-go"
	"github.com/lambdaxs/rocketmq-client-go/consumer"
	"github.com/lambdaxs/rocketmq-client-go/primitive"
)

func main() {
	c, err := rocketmq.NewPushConsumer(
		consumer.WithGroupName("testGroup"),
		consumer.WithNameServer([]string{"127.0.0.1:9876"}),
		consumer.WithConsumerModel(consumer.BroadCasting),
		consumer.WithInstance("instance1"),
	)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if err := c.Subscribe("test1", consumer.MessageSelector{}, func(context context.Context, ext ...*primitive.MessageExt) (result consumer.ConsumeResult, e error) {
		for _, msg := range ext {
			fmt.Printf("receive msg:%s id:%s\n", string(msg.Body), msg.MsgId)
		}
		return consumer.ConsumeSuccess, nil
	}); err != nil {
		fmt.Println("sub error" + err.Error())
	}

	c.Start()
	defer c.Shutdown()

	select {}
}
