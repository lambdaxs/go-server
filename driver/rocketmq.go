package driver

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/lambdaxs/rocketmq-client-go"
    "github.com/lambdaxs/rocketmq-client-go/primitive"
    "github.com/lambdaxs/rocketmq-client-go/producer"
)

type RocketMQConfig struct {
    NameServers []string
    Retry int
    GroupName string
}

type RocketMQProducer struct {
    producer rocketmq.Producer
}

func NewRocketProducer(cfg RocketMQConfig) (pro *RocketMQProducer,err error) {
    pro = &RocketMQProducer{}
    if cfg.Retry == 0 {
        cfg.Retry = 2
    }
    pro.producer,err = rocketmq.NewProducer(
        producer.WithNameServer(cfg.NameServers),
        producer.WithRetry(cfg.Retry),
        producer.WithGroupName(cfg.GroupName),
        )

    if err != nil {
        return
    }
    if err = pro.producer.Start();err != nil {
        return
    }
    return
}

func (p *RocketMQProducer)SendObject(topic string, obj interface{}) error {
    buf,err := json.Marshal(obj)
    if err != nil {
        return err
    }
    return p.SendMsg(topic, buf)
}

func (p *RocketMQProducer)SendMsg(topic string, content []byte) error {
    res,err := p.producer.SendSync(context.Background(), primitive.NewMessage(topic,
        content),
    )
    if err != nil {
        return err
    }
    if res.Status != primitive.SendOK {
        return fmt.Errorf("send status:%d", res.Status)
    }
    return nil
}

func (p *RocketMQProducer)Close() {
    p.producer.Shutdown()
}