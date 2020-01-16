package server

import (
    "fmt"
    "github.com/lambdaxs/go-server/discover"
    "google.golang.org/grpc"
    "net"
    "time"
)

type GRPCServer struct {
    Host string
    Port int
    ConsulAddr string
    ServiceName string
}

func (g *GRPCServer)Schema() string {
    return "grpc"
}

func (g *GRPCServer)StartGRPCServer(registerFunc func(srv *grpc.Server), option ...grpc.ServerOption){
    s := grpc.NewServer(option...)
    if registerFunc != nil {
        registerFunc(s)
    }
    addr := fmt.Sprintf("%s:%d", g.Host, g.Port)
    lis,err := net.Listen("tcp", addr)
    if err != nil {
        fmt.Println("start tcp listen error:"+err.Error())
        return
    }
    //grpc服务注册
    if g.ServiceName != "" && g.ConsulAddr != "" {
        cr := discover.ConsulRegister{
            DCName:  "",
            Address: g.ConsulAddr,
            Ttl:     time.Second * 15,}
        if err := cr.Register(discover.RegisterInfo{
            Host:           g.Host,
            Port:           g.Port,
            ServiceName:    fmt.Sprintf("%s:%s", g.Schema(), g.ServiceName),
            UpdateTime: time.Second});err != nil {
            panic(err)
        }
    }


    fmt.Printf("⇨ grpc server started on \x1b[0;%dm%s\n\x1b[0m",32, fmt.Sprintf("%s:%d", g.Host, g.Port))
    if err := s.Serve(lis);err != nil {
        fmt.Println("start grpc server error:"+err.Error())
    }
}
