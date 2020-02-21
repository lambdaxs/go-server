package discover

import (
    "fmt"
    "github.com/hashicorp/consul/api"
    "os"
    "os/signal"
    "strconv"
    "strings"
    "syscall"
    "time"
)

type RegisterInfo struct {
    Host string
    Port int
    ServiceName string
    Tags []string
    UpdateTime time.Duration
}

type RegisterServer interface {
    Register(info RegisterInfo) error
    UnRegister(info RegisterInfo) error
}

type ConsulRegister struct {
    DCName string
    Address string
    Ttl time.Duration
}

func (c *ConsulRegister)Register(info RegisterInfo) error {
    config := api.DefaultConfig()
    config.Datacenter = c.DCName
    config.Address = c.Address
    client,err := api.NewClient(config)
    if err != nil {
        return err
    }
    //注册服务
    serviceID := genServiceID(info.ServiceName, info.Host, info.Port)
    regRequest := api.AgentServiceRegistration{
        ID:                serviceID,
        Name:              info.ServiceName,
        Tags:              info.Tags,
        Port:              info.Port,
        Address:           info.Host,
    }
    if err := client.Agent().ServiceRegister(&regRequest);err != nil {
        return err
    }
    //设置健康检查
    check := api.AgentServiceCheck{TTL: c.Ttl.String(), Status: api.HealthPassing}
    if err := client.Agent().CheckRegister(
        &api.AgentCheckRegistration{
            ID: serviceID,
            Name: info.ServiceName,
            ServiceID: serviceID,
            AgentServiceCheck: check});err != nil {
        return fmt.Errorf("init check error:%s", err.Error())
    }
    //进程优雅退出
    go func() {
        ch := make(chan os.Signal, 1)
        signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
        x := <-ch
        fmt.Println("consul quit register:",x)
        // un-register service
        if err := c.UnRegister(info);err != nil {
            fmt.Println("consul quit register error:",err.Error())
        }
        s, _ := strconv.Atoi(fmt.Sprintf("%d", x))
        os.Exit(s)
    }()

    //定时刷新服务
    go func() {
        ticker := time.NewTicker(info.UpdateTime)
        for {
            <-ticker.C
            err = client.Agent().UpdateTTL(serviceID, "", check.Status)
            if err != nil {//服务注销后会抛出500异常
                msg := err.Error()
                if !strings.HasPrefix(msg, " Unexpected response code: 500") {
                    fmt.Println("consul update ttl err:",err.Error())
                }
            }
        }
    }()
    return nil
}

func (c *ConsulRegister)UnRegister(info RegisterInfo) error {
    serviceID := genServiceID(info.ServiceName, info.Host, info.Port)

    config := api.DefaultConfig()
    config.Datacenter = c.DCName
    config.Address = c.Address
    client,err := api.NewClient(config)
    if err != nil {
        return err
    }
    //注销服务
    if err := client.Agent().ServiceDeregister(serviceID);err != nil {
        return err
    }
    //注销健康检查
    if err := client.Agent().CheckDeregister(serviceID);err != nil {
        return err
    }
    return nil
}

func genServiceID(name string, host string, port int) string {
    return fmt.Sprintf("%s-%s-%d", name, host, port)
}