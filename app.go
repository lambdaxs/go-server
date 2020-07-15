package go_server

import (
    "context"
    "fmt"
    "github.com/BurntSushi/toml"
    "github.com/gomodule/redigo/redis"
    grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
    "github.com/jinzhu/gorm"

    "github.com/lambdaxs/go-server/govern/confu"
    "github.com/lambdaxs/go-server/govern/log"
    "github.com/lambdaxs/go-server/govern/monitor"

    "github.com/lambdaxs/go-server/driver/mysql_client"
    "github.com/lambdaxs/go-server/driver/psql_client"
    "github.com/lambdaxs/go-server/driver/redis_client"

    "github.com/grpc-ecosystem/go-grpc-middleware"
    "github.com/lambdaxs/go-server/middleware"

    "github.com/labstack/echo"
    "github.com/lambdaxs/go-server/server"

    "go.uber.org/zap"
    "google.golang.org/grpc"
    "io/ioutil"
    "os"
    "os/signal"
    "reflect"
    "syscall"
    "time"
)

type appServer struct {
    ServiceName   string
    AppConfig     *appConfig
    ConfigContent string

    httpSrv *echo.Echo
    gRPCSrv *grpc.Server

    dbMap    map[string]*gorm.DB
    redisMap map[string]*redis.Pool

    serverListen chan struct{}
    stopSign     chan string
}

type appConfig struct {
    HttpServer struct {
        Host       string
        Port       int
        ConsulAddr string
    }
    GrpcServer struct {
        Host       string
        Port       int
        ConsulAddr string
    }
    Log     logConfig
    Monitor struct {
        SystemClose bool
        HttpClose   bool
        GrpcClose   bool
    }
    Tracer struct {
        AgentAddr string
        HttpClose bool
        GrpcClose bool
    }
    Mysql map[string]mysql_client.MysqlConfig `toml:"mysql"`
    Psql  map[string]psql_client.PsqlConfig   `toml:"psql"`
    Redis map[string]redis_client.RedisConfig `toml:"redis"`
}

type logConfig struct {
    log.Config
    HttpClose bool
    GrpcClose bool
}

var defaultServer *appServer

func Default() *appServer {
    return defaultServer
}

func New(serviceName string) *appServer {

    app := &appServer{
        ServiceName:   serviceName,
        AppConfig:     &appConfig{},
        ConfigContent: "",

        httpSrv: nil,
        gRPCSrv: nil,

        dbMap:    map[string]*gorm.DB{},
        redisMap: map[string]*redis.Pool{},

        serverListen: make(chan struct{}, 1),
        stopSign:     make(chan string, 1),
    }

    app.initConfig()

    app.initLogger()

    app.initMonitor()

    app.initTracer()

    app.initSource()

    defaultServer = app

    return app
}

func (app *appServer) HttpServer() *echo.Echo {
    if defaultServer.httpSrv == nil {
        app.initHttpServer()
        return defaultServer.httpSrv
    }
    return defaultServer.httpSrv
}

func (app *appServer) RegisterGRPCServer(reg func(srv *grpc.Server), opts ...grpc.ServerOption) *grpc.Server {
    if defaultServer.gRPCSrv == nil {
        app.initGRPCServer(reg, opts...)
        return defaultServer.gRPCSrv
    }
    return defaultServer.gRPCSrv
}

func (app *appServer) Run() {

    //watch system signal
    app.watchExit()

    msg := <-app.stopSign

    //  graceful close server, default timeout 5s
    if app.httpSrv != nil {
        ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
        if err := app.httpSrv.Shutdown(ctx); err != nil {
            log.Default().Error("stop http server error", zap.String("error", err.Error()))
        } else {
            log.Default().Info("stop http server success")
        }
    }

    // graceful close GRPC server
    if app.gRPCSrv != nil {
        app.gRPCSrv.GracefulStop()
        log.Default().Info("stop GRPC server success")
    }

    // close conn fro database source
    for _, conn := range app.dbMap {
        _ = conn.Close()
    }
    for _, conn := range app.redisMap {
        _ = conn.Close()
    }

    time.Sleep(time.Millisecond * 500)
    log.Default().Info("stop server:" + msg)
}

// load config
func (app *appServer) initConfig() {
    //init config  will be use flag.Parse() function
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

    //load remote config from consul
    if remoteConfigPath != "" {
        buf, err := confu.InitWithRemotePath(remoteConfigPath, app.AppConfig, "")
        if err != nil {
            panic("load remote config err:" + err.Error())
        }
        app.ConfigContent = string(buf)
    }
}

// init logger
func (app *appServer) initLogger() {
    if !reflect.DeepEqual(app.AppConfig.Log, logConfig{}) {
        log.SetLogger(log.NewLogger(app.AppConfig.Log.Config))
    }
}

// init system monitor
func (app *appServer) initMonitor() {
    if !app.AppConfig.Monitor.SystemClose {
        middleware.InitSystemMonitor()
    }
}

// init trace
func (app *appServer) initTracer() {
    // todo
}

// load database source
func (app *appServer) initSource() {
    // mysql
    if len(app.AppConfig.Mysql) != 0 {
        for name, dbConfig := range app.AppConfig.Mysql {
            conn, err := dbConfig.Connect()
            if err != nil {
                panic(fmt.Sprintf("db init err:%s %s dsn:%s", err.Error(), name, dbConfig.DSN))
            }
            app.dbMap[name] = conn
            log.Default().Info(fmt.Sprintf("init db success:%s", name))
        }
    }

    // postgresql
    if len(app.AppConfig.Psql) != 0 {
        for name, dbConfig := range app.AppConfig.Psql {
            conn, err := dbConfig.Connect()
            if err != nil {
                panic(fmt.Sprintf("db init err:%s %s dsn:%s", err.Error(), name, dbConfig.DSN))
            }
            app.dbMap[name] = conn
            log.Default().Info(fmt.Sprintf("init db success:%s", name))
        }
    }

    // redis
    if len(app.AppConfig.Redis) != 0 {
        for name, dbConfig := range app.AppConfig.Redis {
            pool := dbConfig.Connect()
            app.redisMap[name] = pool
            log.Default().Info(fmt.Sprintf("init redis success:%s", name))
        }
    }
}

func (app *appServer) initHttpServer() {
    //start HTTP server
    if app.AppConfig.HttpServer.Port != 0 {
        httpSrv := server.HttpServer{
            Host:        app.AppConfig.HttpServer.Host,
            Port:        app.AppConfig.HttpServer.Port,
            ConsulAddr:  app.AppConfig.HttpServer.ConsulAddr,
            ServiceName: app.ServiceName,
        }

        go httpSrv.StartEchoServer(func(srv *echo.Echo) {
            app.httpSrv = srv

            // log
            if !app.AppConfig.Log.HttpClose {
                srv.Use(middleware.HttpServerLogger())
                //支持动态调整日志等级
                srv.GET("/govern/log/get", log.Default().LogLevel)
                srv.PUT("/govern/log/update", log.Default().LogLevel)
            }

            // monitor
            if !app.AppConfig.Monitor.HttpClose {
                srv.Use(middleware.HttpServerMonitor())
                srv.GET("/govern/metrics", monitor.StartMonitorServer)
            }

            // todo open trace
            if !app.AppConfig.Tracer.HttpClose {

            }

            // default open recover
            srv.Use(middleware.HttpServerRecovery())

            app.serverListen <- struct{}{}
        })

        log.Default().Info(fmt.Sprintf("start http server:%s:%d", httpSrv.Host, httpSrv.Port))
        <-app.serverListen
    }
}

func (app *appServer) initGRPCServer(register func(srv *grpc.Server), opts ...grpc.ServerOption) {
    // start GRPC server
    if app.AppConfig.GrpcServer.Port != 0 {
        grpcSrv := server.GRPCServer{
            Host:        app.AppConfig.GrpcServer.Host,
            Port:        app.AppConfig.GrpcServer.Port,
            ServiceName: app.ServiceName,
        }

        middlewareList := make([]grpc.UnaryServerInterceptor, 0)
        // log
        if !app.AppConfig.Log.GrpcClose {
            middlewareList = append(middlewareList, middleware.GRPCServerLogger())
        }

        // monitor
        if !app.AppConfig.Monitor.GrpcClose {
            middlewareList = append(middlewareList, middleware.GRPCServerMonitor())
        }

        // 开启链路 todo 明确标示头
        if !app.AppConfig.Tracer.GrpcClose {
            middlewareList = append(middlewareList, middleware.GRPCServerTracer())
        }

        //添加recovery
        middlewareList = append(middlewareList, middleware.GRPCServerRecovery())

        opts = append(opts, grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(middlewareList...)))

        go grpcSrv.StartGRPCServer(func(srv *grpc.Server) {
            app.gRPCSrv = srv

            if register != nil {
                register(srv)
            }

            // register grpc monitor
            if !app.AppConfig.Monitor.GrpcClose {
                grpc_prometheus.Register(srv)
            }

            app.serverListen <- struct{}{}
        }, opts...)

        log.Default().Info(fmt.Sprintf("start grpc server:%s:%d", grpcSrv.Host, grpcSrv.Port))
        <-app.serverListen
    }
}

func (a *appServer) watchExit() {
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
    go func() {
        sig := <-sigs
        a.stopSign <- sig.String()
    }()
}

// todo support yaml json format
func (a *appServer) ParseConfig(i interface{}) error {
    _, err := toml.Decode(a.ConfigContent, i)
    if err != nil {
        return err
    }
    return nil
}

func ParseConfig(i interface{}) error {
    return defaultServer.ParseConfig(i)
}

func Model(name string) *gorm.DB {
    return defaultServer.dbMap[name]
}

func RedisPool(name string) *redis.Pool {
    return defaultServer.redisMap[name]
}

func ConfigContent() string {
    return defaultServer.ConfigContent
}

func HttpServer() *echo.Echo {
    return defaultServer.HttpServer()
}

func GRPCServer() *grpc.Server {
    return defaultServer.gRPCSrv
}
