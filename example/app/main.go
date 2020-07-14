package main

import (
    "context"
    "fmt"
    "github.com/labstack/echo"
    go_server "github.com/lambdaxs/go-server"
    hello "github.com/lambdaxs/go-server/example/discover/pb"
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

type user struct {
    ID   int64  `json:"id"`
    Name string `json:"name"`
    Age  int64  `json:"age"`
}

type ConfigData struct {
    Server struct{
        Name string
    }
}

func main() {
    app := go_server.New("test")

    config := &ConfigData{}
    if err := app.ParseConfig(config);err != nil {
        panic(err)
    }

    fmt.Println(config.Server.Name)

    psqlDB := go_server.Model("test")
    mysqlDB := go_server.Model("db")

    srv := app.HttpServer()
    srv.GET("/", func(c echo.Context) error {
    	list := make([]user, 0)
        psqlDB.Table("user").Find(&list)
        return c.JSON(200, list)
    })

	srv.GET("/mysql", func(c echo.Context) error {
		list := make([]user, 0)
		mysqlDB.Table("user").Find(&list)
		return c.JSON(200, list)
	})

    app.RegisterGRPCServer(func(srv *grpc.Server) {
        hello.RegisterHelloServerServer(srv, &SayHelloServer{})
    })


    app.Run()
}
