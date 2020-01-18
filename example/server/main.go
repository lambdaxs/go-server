package main

import (
    "context"
    "flag"
    "fmt"
    "github.com/labstack/echo"
    "github.com/lambdaxs/go-server/confu"
    hello "github.com/lambdaxs/go-server/example/discover/pb"
    "github.com/lambdaxs/go-server/server"
    "github.com/lambdaxs/go-server/lib/validate"
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
    confu.ParseFlag()
    flag.Parse()

    httpSrv := server.HttpServer{
        Host: "127.0.0.1",
        Port: 9000,
    }
    httpSrv.StartEchoServer(func(srv *echo.Echo) {
        srv.POST("/", func(c echo.Context) error {
            reqModel := new(struct {
                Uid   int64  `json:"uid" form:"uid" validate:"required"`
                Age   int64  `json:"age" form:"age" validate:"required,gte=0,lte=130"`
                Email string `json:"email" validate:"required,email"`
                Code string `json:"code" validate:"required,len=4"`
                Plat string `json:"plat" validate:"required,oneof=ios android"`
            })
            if err := c.Bind(reqModel); err != nil {
                return c.JSON(200, err.Error())
            }
            if err := validate.Struct(reqModel);err != nil {
                return c.JSON(200, err.Error())
            }
            return c.JSON(200, "success")
        })
    })
}
