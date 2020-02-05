package main

import (
    "fmt"
    "github.com/nats-io/stan.go"
    "time"
)

type Connection struct {
    stan.Conn
    reconnectChan map[string]chan struct{}
}

func Connect(cid string, clientID string, options ...stan.Option) (*Connection, error) {
    c := &Connection{Conn: nil, reconnectChan: make(map[string]chan struct{})}

    var closeConnectFunc func(conn stan.Conn, e error)
    closeConnectFunc = func(conn stan.Conn, e error) {
        go func() {
            for {
                time.Sleep(time.Second * 3)
                options = append(options, stan.SetConnectionLostHandler(closeConnectFunc))
                sc, err := stan.Connect(cid, clientID, options...)
                if err == nil { // 重连成功
                    c.Conn = sc
                    fmt.Println("重连成功")
                    for _, v := range c.reconnectChan {
                        v <- struct{}{}
                    }
                    break
                }
            }
        }()
    }

    options = append(options, stan.SetConnectionLostHandler(closeConnectFunc))
    sc, err := stan.Connect(cid, clientID, options...)
    if err != nil {
        return nil, err
    }
    c.Conn = sc
    return c, nil
}

func (c *Connection)Publish(subject string, data []byte) error{
    return c.Conn.Publish(subject, data)
}

func (c *Connection)PublishAsync(subject string, data []byte, ah stan.AckHandler) (string, error) {
    return c.Conn.PublishAsync(subject, data, ah)
}

func (c *Connection)Subscribe(subject string, cb stan.MsgHandler, opts ...stan.SubscriptionOption) (stan.Subscription, error){
    var rs stan.Subscription
    var err error

    callsub := func() error {
        rs, err = c.Conn.Subscribe(subject, cb, opts...)
        if err != nil {
            return err
        }
        return nil
    }
    if err = callsub(); err != nil {
        return rs, err
    }
    c.reconnectChan["sub_"+subject] = make(chan struct{}, 1)

    go func() {
        for {
            select {
            case <-c.reconnectChan["sub_"+subject]: //触发重连
                fmt.Println("重新订阅:" + subject)
                if err = callsub(); err != nil {
                    time.Sleep(time.Second * 3)
                    continue
                }
            }
        }
    }()

    return rs, err
}


func (c *Connection) QueueSubscribe(subject, qgroup string, cb stan.MsgHandler, opts ...stan.SubscriptionOption) (stan.Subscription, error) {
    var rs stan.Subscription
    var err error

    callsub := func() error {
        rs, err = c.Conn.QueueSubscribe(subject, qgroup, cb, opts...)
        if err != nil {
            return err
        }
        return nil
    }
    if err = callsub(); err != nil {
        return rs, err
    }
    c.reconnectChan["queue_"+subject] = make(chan struct{}, 1)

    go func() {
        for {
            select {
            case <-c.reconnectChan["queue_"+subject]: //触发重连
                fmt.Println("重新订阅:" + subject)
                if err = callsub(); err != nil {
                    time.Sleep(time.Second * 3)
                    continue
                }
            }
        }
    }()

    return rs, err
}

func main() {
    sc, err := Connect("c1", "server-host-2", stan.NatsURL("nats://49.235.146.124:4222"))
    if err != nil {
        fmt.Println("err:" + err.Error())
        return
    }

    fmt.Println("初始化连接")

    startOpt := stan.DeliverAllAvailable()
    _, err = sc.QueueSubscribe("foo", "consumer_group_1", func(msg *stan.Msg) {
        fmt.Printf("接收到消息：%s time:%d seq:%d\n", string(msg.Data), msg.Timestamp, msg.Sequence)
        msg.Ack()
    }, startOpt, stan.DurableName("consumer_1"), stan.SetManualAckMode())
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    _, err = sc.QueueSubscribe("bar", "consumer_group_1", func(msg *stan.Msg) {
        fmt.Printf("接收到消息：%s time:%d seq:%d\n", string(msg.Data), msg.Timestamp, msg.Sequence)
        msg.Ack()
    }, startOpt, stan.DurableName("consumer_1"), stan.SetManualAckMode())

    select {}
}

func (c *Connection)Close() error {
    return c.Conn.Close()
}