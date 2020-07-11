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

type user struct {
    ID   int64  `json:"id"`
    Name string `json:"name"`
    Age  int64  `json:"age"`
}

func main() {
    app := go_server.New("test")

    psqlDB := go_server.Model("test")

    app.HttpServer().GET("/", func(c echo.Context) error {
    	list := make([]user, 0)
        psqlDB.Table("public.user").Find(&list)
        return c.JSON(200, list)
    })

    hello.RegisterHelloServerServer(app.GRPCServer(), &SayHelloServer{})

    app.Run()
}
