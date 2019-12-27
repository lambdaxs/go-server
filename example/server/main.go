package main

import (
    "context"
    "flag"
    "fmt"
    "github.com/labstack/echo"
    hello "github.com/lambdaxs/go-server/example/discover/pb"
    "github.com/lambdaxs/go-server/server"
    "google.golang.org/grpc"
)

type SayHelloServer struct {
}

func (s *SayHelloServer) SayHello(ctx context.Context, req *hello.SayHelloReq) (resp *hello.SayHelloResp, err error) {
    resp = &hello.SayHelloResp{
        Reply: "",
    }
    resp.Reply = fmt.Sprintf("%s:%s", req.GetName(), "hello")
    return
}

func main() {
    var isTest bool
    var remoteConfig string
    flag.BoolVar(&isTest, "t", false, "pre-test server config")
    flag.StringVar(&remoteConfig, "remote-config", "", "remote server config")
    flag.Parse()
    if isTest {
        fmt.Println(remoteConfig)
        return
    }

    httpSrv := server.HttpServer{
        Host:        "127.0.0.1",
        Port:        9000,
        ConsulAddr: "127.0.0.1:8500",
        ServiceName: "HelloService",
    }
    go httpSrv.StartEchoServer(func(srv *echo.Echo) {
        srv.GET("/", func(c echo.Context) error {
            return c.JSON(200, "hello world");
        })
    })

    grpcSrv := server.GRPCServer{
        Host:       "127.0.0.1",
        Port:       9002,
        ConsulAddr: "127.0.0.1:8500",
        ServiceName: "HelloService",
    }
    grpcSrv.StartGRPCServer(func(srv *grpc.Server) {
        hello.RegisterHelloServerServer(srv, &SayHelloServer{})
    })
}
