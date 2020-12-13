## 特性

- 服务注册发现
- grpc负载均衡
- 配置中心
- 监控/日志/链路
- 常用数据库驱动

### 配置管理

- consul
- admin
- config version

### 日志管理

- EFK
- kafka
- agent monitor

### 服务监控

- prometheus
- grafna
- consul + confd

### 链路监控

- jaeger + agent + collector

### 服务注册

- consul
- agent health
- service health
- cluster alert

### 机器和部署管理

- CMDB
- ansible + Makefile + image version


### 用户管理

- user
- group
- auth grant

### 测试

- sonar
- unit



## Example


- config.toml

```toml
[httpServer]
    host = "127.0.0.1"
    port = 8000

[grpcServer]
    host = "127.0.0.1"
    port = 8001

[mysql]
    [mysql.db]
        dsn = "root:123456@tcp(127.0.0.1:3306)/dbname?charset=utf8&parseTime=True&loc=Local&readTimeout=3s"
        log = true

[redis]
    [redis.cache]
        dsn = "127.0.0.1:6379"
```

- main.go

```go
package main

import (
    "context"
    "flag"
    "fmt"
    "github.com/labstack/echo"
    hello "github.com/lambdaxs/go-server/example/discover/pb"
    "github.com/lambdaxs/go-server"
    "google.golang.org/grpc"
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
    app := go_server.New("lijia-server")
    
    //start http server
    httpSrv := app.HttpServer()
    httpSrv.GET("/hello", func(c echo.Context) error {
        return c.JSON(200,"Hello Go-Server!")
    })
    
    //start grpc server
    app.RegisterGRPCServer(func(srv *grpc.Server) {
        hello.RegisterHelloServerServer(srv, &SayHelloServer{})             
    })   

    //server start
    app.Run()
    
}

```

- 启动程序

```shell script
go run main.go --config config.toml
```

## 配置文件

- http服务配置

```toml
[httpServer]
    host = "127.0.0.1"
    port = 8000
```

- grpc服务配置

```toml
[grpcServer]
    host = "127.0.0.1"
    port = 8001
```

- mysql配置

```toml
[mysql]
    [mysql.db]
        dsn = "root:123456@tcp(127.0.0.1:3306)/dbname?charset=utf8&parseTime=True&loc=Local&readTimeout=3s"
        log = true
```

```golang
//get mysql db object
db := go_server.Model("db")

```

- redis配置

```toml
[redis]
    [redis.cache]
        dsn = "127.0.0.1:6379"
```

```golang
//get redis connect pool
pool := go_server.RedisPool("cache")
conn := pool.Get()
defer conn.close()
```