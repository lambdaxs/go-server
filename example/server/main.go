package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/labstack/echo"
	"github.com/lambdaxs/go-server/confu"
	hello "github.com/lambdaxs/go-server/example/discover/pb"
	"github.com/lambdaxs/go-server/lib/validate"
	"github.com/lambdaxs/go-server/server"
	"google.golang.org/grpc"
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
	flag.Parse()

	//logger := log.NewLogger(log.Config{
	//    Development: true,
	//})
	//logger.Info("info",zap.String("key", "value"))

	httpSrv := server.HttpServer{
		Host: "127.0.0.1",
		Port: 9002,
		ConsulAddr: "127.0.0.1:8500",
		ServiceName:"test",
	}

	go httpSrv.StartEchoServer(func(srv *echo.Echo) {
		srv.POST("/", func(c echo.Context) error {
			reqModel := new(struct {
				UID   int64  `json:"uid" form:"uid" validate:"required"`
				Age   int64  `json:"age" form:"age" validate:"required,gte=0,lte=130"`
				Email string `json:"email" validate:"required,email"`
				Code  string `json:"code" validate:"required,len=4"`
				Plat  string `json:"plat" validate:"required,oneof=ios android"`
			})
			if err := c.Bind(reqModel); err != nil {
				return c.JSON(200, err.Error())
			}
			if err := validate.Struct(reqModel); err != nil {
				return c.JSON(200, err.Error())
			}
			return c.JSON(200, "success")
		})
	})

	grpcServer := server.GRPCServer{
		Host: "127.0.0.1",
		Port: 9093,
	}
	grpcServer.StartGRPCServer(func(srv *grpc.Server) {
		hello.RegisterHelloServerServer(srv, &SayHelloServer{})
	})
}
