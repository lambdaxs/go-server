package main

import (
	"github.com/labstack/echo"
	go_server "github.com/lambdaxs/go-server"
)

func main() {
	app := go_server.New("test")

	app.HttpSrv.GET("/", func(c echo.Context) error {
		return c.JSON(200, "success")
	})

	app.Run()
}
