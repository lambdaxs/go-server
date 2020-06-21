package main

import (
	"context"
	"fmt"
	pb "github.com/lambdaxs/go-server/example/discover/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
)

func main() {
	//schema, err := discover.GetConnSchema("127.0.0.1:8500", "HelloService")
	//if err != nil {
	//	panic(err)
	//}

	//conn, err := grpc.Dial(fmt.Sprintf("%s:///HelloService", schema),
	//	grpc.WithInsecure(),
	//	grpc.WithBalancerName(roundrobin.Name))
	//if err != nil {
	//	panic(err)
	//}
	conn, err := grpc.Dial("127.0.0.1:8081",
		grpc.WithInsecure(),
		grpc.WithBalancerName(roundrobin.Name))
	if err != nil {
		panic(err)
	}

	client := pb.NewHelloServerClient(conn)

	resp, err := client.SayHello(context.Background(), &pb.SayHelloReq{
		Name: "xiaosaaa",
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(resp.GetReply())
}
