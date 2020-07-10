package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo"
	go_server "github.com/lambdaxs/go-server"
	hello "github.com/lambdaxs/go-server/example/discover/pb"
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
	app := go_server.New("test")

	app.HttpServer().GET("/", func(c echo.Context) error {
		return c.JSON(200, "success")
	})

	hello.RegisterHelloServerServer(app.GRPCServer(), &SayHelloServer{})


	app.Run()
}
