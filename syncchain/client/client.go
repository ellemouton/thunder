package client

import (
	"context"
	"flag"

	"github.com/ellemouton/thunder/syncchain"
	pb "github.com/ellemouton/thunder/syncchain/syncchainpb"
	"google.golang.org/grpc"
)

var _ syncchain.Client = (*client)(nil)

var grpcAddr = flag.String("grpc_address", ":8081", "syncchain grpc address")

type client struct {
	rpcConn   *grpc.ClientConn
	rpcClient pb.SyncchainClient
}

func New() (*client, error) {
	var (
		c   client
		err error
	)

	c.rpcConn, err = grpc.Dial(*grpcAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	c.rpcClient = pb.NewSyncchainClient(c.rpcConn)

	return &c, nil
}

func (c *client) Ping(ctx context.Context, msg string) (string, error) {
	resp, err := c.rpcClient.Ping(ctx, &pb.PingRequest{Msg: msg})
	if err != nil {
		return "", err
	}

	return resp.Msg, nil
}

func (c *client) Close() error {
	return c.rpcConn.Close()
}
