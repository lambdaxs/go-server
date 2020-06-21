package go_server

import (
	"context"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/lambdaxs/go-server/confu"
	"github.com/lambdaxs/go-server/driver/mysql_client"
	"github.com/lambdaxs/go-server/driver/redis_client"
	"github.com/lambdaxs/go-server/log"
	"github.com/lambdaxs/go-server/server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type appServer struct {
	ServiceName string
	HttpSrv     *echo.Echo
	GRPCSrv     *grpc.Server
	grpcOptions []grpc.ServerOption

	DBMap         map[string]*gorm.DB
	RedisMap      map[string]*redis.Pool
	AppConfig     *appConfig
	ConfigContent string

	serverListen chan struct{}
	StopSign     chan string
}

type appConfig struct {
	HttpServer struct {
		Host string
		Port int
	}
	GrpcServer struct {
		Host string
		Port int
	}
	Log   log.Config
	Mysql map[string]mysql_client.MysqlDB `toml:"mysql"`
	Redis map[string]redis_client.RedisDB `toml:"redis"`
}

func New(serviceName string) *appServer {

	app := &appServer{
		ServiceName:   serviceName,
		HttpSrv:       nil,
		GRPCSrv:       nil,
		grpcOptions: []grpc.ServerOption{},

		DBMap:         map[string]*gorm.DB{},
		RedisMap:      map[string]*redis.Pool{},
		AppConfig:     &appConfig{},
		ConfigContent: "",
		serverListen:  make(chan struct{}, 1),
		StopSign:      make(chan string, 1),
	}

	app.initConfig()
	app.initLogger()
	app.initSource()
	app.initHttpServer()
	app.initGRPCServer()

	return app
}

func (app *appServer) Run() {

	//监听信号
	app.watchExit()

	msg := <-app.StopSign

	// 优雅关闭http服务器,默认超时5s
	if app.HttpSrv != nil {
		ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
		if err := app.HttpSrv.Shutdown(ctx); err != nil {
			log.Default().Error("stop http server error", zap.String("error", err.Error()))
		} else {
			log.Default().Info("stop http server success")
		}
	}

	//优雅关闭GRPC服务
	if app.GRPCSrv != nil {
		app.GRPCSrv.GracefulStop()
		log.Default().Info("stop GRPC server success")
	}

	// 优雅关闭数据库资源
	for _, conn := range app.DBMap {
		conn.Close()
	}
	for _, conn := range app.RedisMap {
		conn.Close()
	}

	time.Sleep(time.Millisecond * 500)
	log.Default().Info("stop server:" + msg)
}

// 加载配置
func (app *appServer) initConfig() {
	//初始化配置 会引用flag.Parse()方法
	configPath, remoteConfigPath := confu.ParseFlag()
	if configPath != "" {
		if err := confu.InitWithFilePath(configPath, app.AppConfig); err != nil {
			panic("local config file load err:" + err.Error())
		}
		buf, err := ioutil.ReadFile(configPath)
		if err != nil {
			panic("local config file load err:" + err.Error())
		}

		app.ConfigContent = string(buf)
	}

	//加载远端配置
	if remoteConfigPath != "" {
		buf, err := confu.InitWithRemotePath(remoteConfigPath, app.AppConfig, "")
		if err != nil {
			panic("load remote config err:" + err.Error())
		}
		app.ConfigContent = string(buf)
	}
}

// 初始化日志输出
func (app *appServer) initLogger() {
	if app.AppConfig.Log.FilePath != "" {
		log.SetLogger(log.NewLogger(app.AppConfig.Log))
	}
}

// 加载数据库资源
func (app *appServer) initSource() {
	//初始化数据库
	if len(app.AppConfig.Mysql) != 0 {
		for name, dbConfig := range app.AppConfig.Mysql {
			conn, err := dbConfig.ConnectGORMDB()
			if err != nil {
				panic(fmt.Sprintf("db init err:%s %s dsn:%s", err.Error(), name, dbConfig.DSN))
			}
			app.DBMap[name] = conn
			log.Default().Info(fmt.Sprintf("init db success:%s", name))
		}
	}

	//初始化redis
	if len(app.AppConfig.Redis) != 0 {
		for name, dbConfig := range app.AppConfig.Redis {
			pool := dbConfig.ConnectRedisPool()
			app.RedisMap[name] = pool
			log.Default().Info(fmt.Sprintf("init redis success:%s", name))
		}
	}
}

func (app *appServer) initHttpServer() {
	//启动HTTP服务器
	if app.AppConfig.HttpServer.Port != 0 {
		httpSrv := server.HttpServer{
			Host:        app.AppConfig.HttpServer.Host,
			Port:        app.AppConfig.HttpServer.Port,
			ServiceName: app.ServiceName,
		}

		go httpSrv.StartEchoServer(func(srv *echo.Echo) {
			app.HttpSrv = srv
			app.serverListen <- struct{}{}
		})
		log.Default().Info(fmt.Sprintf("start http server:%s:%d", httpSrv.Host, httpSrv.Port))
		<-app.serverListen
	}
}

// GRPC server options
func (app *appServer) SetGRPCOptions(option ...grpc.ServerOption) {
	for _, opt := range option {
		app.grpcOptions = append(app.grpcOptions, opt)
	}
}

func (app *appServer) initGRPCServer() {
	//启动GRPC服务器
	if app.AppConfig.GrpcServer.Port != 0 {
		grpcSrv := server.GRPCServer{
			Host:        app.AppConfig.GrpcServer.Host,
			Port:        app.AppConfig.GrpcServer.Port,
			ServiceName: app.ServiceName,
		}

		go grpcSrv.StartGRPCServer(func(srv *grpc.Server) {
			app.GRPCSrv = srv
			app.serverListen <- struct{}{}
		}, app.grpcOptions...)
		log.Default().Info(fmt.Sprintf("start grpc server:%s:%d", grpcSrv.Host, grpcSrv.Port))
		<-app.serverListen
	}
}

func (a *appServer) watchExit() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		sig := <-sigs
		a.StopSign <- sig.String()
	}()
}
