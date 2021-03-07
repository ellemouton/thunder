package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/ellemouton/thunder/syncchain/server"
	"github.com/ellemouton/thunder/syncchain/syncchainpb"
)

var grpcAddr = flag.String("grpc_address", ":8081", "syncchain grpc address")

func main() {
	fmt.Println("syncchain is running")

	lis, err := net.Listen("tcp", *grpcAddr)
	if err != nil {
		log.Fatal(err)
	}

	syncchainServer := server.New()
	grpcServer := grpc.NewServer()

	syncchainpb.RegisterSyncchainServer(grpcServer, syncchainServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to server syncchain", err)
	}

}
