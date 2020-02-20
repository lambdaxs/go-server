package nats_client

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

//初始化nats连接
func SetConnOptions(opts []nats.Option) []nats.Option {
	totalWait := 10 * time.Minute
	reconnectDelay := time.Second

	opts = append(opts, nats.ReconnectWait(reconnectDelay))                        //设置重连等待时间
	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))         //设置最大重试次数
	opts = append(opts, nats.DisconnectErrHandler(func(conn *nats.Conn, e error) { //打印链接断开日志
		log.Printf("nats client disconnect:%s \n", e.Error())
	}))
	opts = append(opts, nats.ReconnectHandler(func(conn *nats.Conn) { //打印重连日志
		log.Println("nats client reconnect")
	}))
	opts = append(opts, nats.ClosedHandler(func(conn *nats.Conn) { //打印链接断开日志
		log.Println("nats client exit:%v", conn.LastError())
	}))
	return opts
}
