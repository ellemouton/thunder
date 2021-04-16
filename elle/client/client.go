package client

import (
        "context"
        "flag"
        "github.com/ellemouton/thunder/elle"
        pb "github.com/ellemouton/thunder/elle/ellepb"
        "google.golang.org/grpc"
)

var _ elle.Client = (*client)(nil)

var grpcAddr = flag.String("grpc_address", ":8082", "elle grpc address")

type client struct {
        rpcConn   *grpc.ClientConn
        rpcClient pb.ElleClient
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

        c.rpcClient = pb.NewElleClient(c.rpcConn)

        return &c, nil
}

func (c client) RequiresPayment(ctx context.Context, path string) (*elle.AssetDetails, error) {
        resp, err := c.rpcClient.RequiresPayment(ctx, &pb.RequiresPaymentRequest{
                Path: path,
        })

        if err != nil {
                return nil, err
        }

        return &elle.AssetDetails{
                ID:              resp.Id,
                RequiresPayment: resp.RequiresPayment,
                Price:           resp.Price,
                Memo:            resp.Memo,
        }, nil
}

func (c client) Close() error {
        return c.rpcConn.Close()
}


