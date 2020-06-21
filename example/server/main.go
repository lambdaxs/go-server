package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/labstack/echo"
	"github.com/lambdaxs/go-server/confu"
	hello "github.com/lambdaxs/go-server/example/discover/pb"
	"github.com/lambdaxs/go-server/server"
)

//SayHelloServer server
type SayHelloServer struct {
}

//SayHello hanler
func (s *SayHelloServer) SayHello(ctx context.Context, req *hello.SayHelloReq) (resp *hello.SayHelloResp, err error) {
	resp = &hello.SayHelloResp{
		Reply: "",
	}
	resp.Reply = fmt.Sprintf("%s:%s", req.GetName(), "hello")
	return
}

func main() {
	confu.ParseFlag()

	var port int64
	flag.Int64Var(&port, "port", 0, "")
	flag.Parse()

	//logger := log.NewLogger(log.Config{
	//    Development: true,
	//})
	//logger.Info("info",zap.String("key", "value"))

	httpSrv := server.HttpServer{
		Host:        "127.0.0.1",
		Port:        int(port),
		ConsulAddr:  "127.0.0.1:8500",
		ServiceName: "test",
	}

	httpSrv.StartEchoServer(func(srv *echo.Echo) {
		srv.GET("/", func(c echo.Context) error {
			return c.JSON(200, "success")
		})
	})

	//grpcServer := server.GRPCServer{
	//	Host: "127.0.0.1",
	//	Port: 9093,
	//}
	//grpcServer.StartGRPCServer(func(srv *grpc.Server) {
	//	hello.RegisterHelloServerServer(srv, &SayHelloServer{})
	//})
}
