package grpc_client

import (
	"fmt"
	"github.com/lambdaxs/go-server/discover"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
)

type GRPCClient struct {
	ConsulAddr  string //127.0.0.1:8500
	ServiceName string
}

//获取grpc连接
func (g *GRPCClient) GetConn() (*grpc.ClientConn, error) {
	schema, err := discover.GetConnSchema(g.ConsulAddr, g.ServiceName)
	if err != nil {
		panic(err)
	}
	conn, err := grpc.Dial(fmt.Sprintf("%s:///%s", schema, g.ServiceName),
		grpc.WithInsecure(),
		grpc.WithBalancerName(roundrobin.Name))
	if err != nil {
		panic(err)
	}
	return conn, nil
}
