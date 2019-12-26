package main

import (
    "context"
    "fmt"
    "github.com/lambdaxs/go-server/discover"
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
    "net"
    "time"
    pb "github.com/lambdaxs/go-server/demo/discover/pb"
)

type server struct {

}

func (s *server)SayHello(ctx context.Context, req *pb.SayHelloReq) (resp *pb.SayHelloResp,err error) {
    resp = &pb.SayHelloResp{
        Reply:                "",
    }
    resp.Reply = fmt.Sprintf("hello:%s", req.GetName())
    return
}

func main(){
    host := "127.0.0.1"
    port := 8080
    consulPort := 8500

    listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(host), port, "",})
    if err != nil {
        fmt.Println(err.Error())
    }
    s := grpc.NewServer()

    // register service
    cr := discover.ConsulRegister{
        DCName:  "",
        Address: fmt.Sprintf("%s:%d", host, consulPort),
        Ttl:     time.Second * 15,
    }
    if err := cr.Register(discover.RegisterInfo{
        Host:           host,
        Port:           port,
        ServiceName:    "HelloService",
        UpdateTime: time.Second});err != nil {
            panic(err)
    }

    pb.RegisterHelloServerServer(s, &server{})
    reflection.Register(s)
    if err := s.Serve(listen); err != nil {
        fmt.Println("failed to serve:" + err.Error())
    }

}