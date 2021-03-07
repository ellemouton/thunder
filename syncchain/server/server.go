package server

import (
	"context"
	"fmt"

	pb "github.com/ellemouton/thunder/syncchain/syncchainpb"
)

var _ pb.SyncchainServer = (*Server)(nil)

type Server struct{}

func New() *Server {
	return &Server{}
}

func (srv *Server) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	return &pb.PingResponse{
		Msg: fmt.Sprintf("got your ping: %s", req.Msg),
	}, nil
}
