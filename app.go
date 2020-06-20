package go_server

import (
    "github.com/labstack/echo"
    "github.com/lambdaxs/go-server/confu"
    "github.com/lambdaxs/go-server/driver/mysql_client"
    "github.com/lambdaxs/go-server/driver/redis_client"
    "github.com/lambdaxs/go-server/server"
    "io/ioutil"
)


type appConfig struct {
    HttpServer struct{
        Host string
        Port int
    }
    GrpcServer struct{
        Host string
        Port int
    }
    Mysql map[string]mysql_client.MysqlDB `toml:"mysql"`
    Redis map[string]redis_client.RedisDB `toml:"redis"`
    Content string
}

var AppConfig *appConfig

func New(serviceName string){

    configPath,remoteConfigPath := confu.ParseFlag()
    if configPath != "" {
        if err := confu.InitWithFilePath(configPath, AppConfig);err != nil {
            panic("local config file load err:"+err.Error())
        }
        buf,err := ioutil.ReadFile(configPath)
        if err != nil {
            panic("local config file load err:"+err.Error())
        }
        AppConfig.Content = string(buf)
    }

    if remoteConfigPath != "" {
         buf,err := confu.InitWithRemotePath(remoteConfigPath, AppConfig, "")
         if err != nil {
             panic("load remote config err:"+err.Error())
         }
         AppConfig.Content = string(buf)
    }

    httpSrv := server.HttpServer{
        Host:        AppConfig.HttpServer.Host,
        Port:        AppConfig.HttpServer.Port,
        ServiceName: serviceName,
    }

    httpSrv.StartEchoServer(func(srv *echo.Echo) {
        srv.GET("/" , handler.GiftServer.Submit)
    })
}