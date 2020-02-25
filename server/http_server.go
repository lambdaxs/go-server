package server

import (
	"fmt"
	"os"
	"time"

	"github.com/labstack/echo"
	"github.com/lambdaxs/go-server/discover"
	"github.com/lambdaxs/go-server/lib/local"
	"github.com/lambdaxs/go-server/log"
	"go.uber.org/zap"
)

type HttpServer struct {
	Host        string
	Port        int
	ConsulAddr  string
	ServiceName string
}

func (h *HttpServer) Schema() string {
	return "http"
}

func (h *HttpServer) StartEchoServer(serverFunc func(srv *echo.Echo)) {
	app := echo.New()
	if serverFunc != nil {
		serverFunc(app)
	}
	if h.Host == "" { //默认使用内网ip
		h.Host = local.LocalIP()
	}
	if h.ConsulAddr == "" {
		h.ConsulAddr = os.Getenv("CONSUL_ADDR")
	}
	//http服务注册
	if h.ServiceName != "" && h.ConsulAddr != "" {
		cr := discover.ConsulRegister{
			DCName:  "",
			Address: h.ConsulAddr,
			Ttl:     time.Second * 15}
		if err := cr.Register(discover.RegisterInfo{
			Host:        h.Host,
			Port:        h.Port,
			ServiceName: fmt.Sprintf("%s:%s", h.Schema(), h.ServiceName),
			UpdateTime:  time.Second * 5}); err != nil {
			panic(err)
		}
	}
	address := fmt.Sprintf("%s:%d", h.Host, h.Port)
	log.Default().Info(
		"⇨ start http server",
		zap.String("address", address))
	if err := app.Start(address); err != nil {
		fmt.Println("start echo server error:" + err.Error())
	}
}
