package main

import (
	"context"
	"fmt"
	"github.com/lambdaxs/go-server/discover"
	pb "github.com/lambdaxs/go-server/example/discover/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"time"
)

type server struct {
}

func (s *server) SayHello(ctx context.Context, req *pb.SayHelloReq) (resp *pb.SayHelloResp, err error) {
	resp = &pb.SayHelloResp{
		Reply: "",
	}
	resp.Reply = fmt.Sprintf("hello:%s", req.GetName())
	return
}

func main() {
	host := "127.0.0.1"
	port := 8080
	consulPort := 8500

	s := grpc.NewServer()

	// register service
	cr := discover.ConsulRegister{
		DCName:  "",
		Address: fmt.Sprintf("%s:%d", host, consulPort),
		Ttl:     time.Second * 15}
	if err := cr.Register(discover.RegisterInfo{
		Host:        host,
		Port:        port,
		ServiceName: "HelloService",
		UpdateTime:  time.Second}); err != nil {
		panic(err)
	}

	pb.RegisterHelloServerServer(s, &server{})
	reflection.Register(s)

	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		fmt.Println(err.Error())
	}
	if err := s.Serve(listen); err != nil {
		fmt.Println("failed to serve:" + err.Error())
	}

}
