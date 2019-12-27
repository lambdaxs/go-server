package server

import (
    "fmt"
    "github.com/labstack/echo"
    "github.com/lambdaxs/go-server/discover"
    "time"
)

type HttpServer struct {
    Host string
    Port int
    ConsulAddr string
    ServiceName string
}

func (h *HttpServer)Schema() string {
    return "http"
}

func (h *HttpServer)StartEchoServer(serverFunc func(srv *echo.Echo)){
    app := echo.New()
    if serverFunc != nil {
        serverFunc(app)
    }
    //http服务注册
    if h.ServiceName != "" && h.ConsulAddr != "" {
        cr := discover.ConsulRegister{
            DCName:  "",
            Address: h.ConsulAddr,
            Ttl:     time.Second * 15,}
        if err := cr.Register(discover.RegisterInfo{
            Host:           h.Host,
            Port:           h.Port,
            ServiceName:    fmt.Sprintf("%s:%s", h.Schema(), h.ServiceName),
            UpdateTime: time.Second});err != nil {
            panic(err)
        }
    }

    if err := app.Start(fmt.Sprintf("%s:%d", h.Host, h.Port));err != nil {
        fmt.Println("start echo server error:"+err.Error())
    }
}

