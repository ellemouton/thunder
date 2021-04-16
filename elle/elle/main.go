package main

import (
	"flag"
	"github.com/ellemouton/thunder/elle/ellepb"
	"github.com/ellemouton/thunder/elle/server"
	"google.golang.org/grpc"
	"html/template"
	"log"
	"net"
	"net/http"
)

var (
	templates = template.Must(template.ParseGlob("assets/*.html"))
	httpAddr = flag.String("httpAddr", ":8080", "address to serve on")
	grpcAddr = flag.String("grpc_address", ":8082", "elle grpc address")
)

func main() {
	flag.Parse()

	s, err := newState()
	if err != nil {
		log.Fatalf("newState: %s", err)
	}

	//ops.StartLoops(s)

	r := newRouter(s)

	go serverGRPCForever(s)

	log.Printf("Serving on port %s", *httpAddr)
	log.Fatal(http.ListenAndServe(*httpAddr, r))
}

func serverGRPCForever(s *State) {
	lis, err := net.Listen("tcp", *grpcAddr)
	if err != nil {
		log.Fatal(err)
	}

	elleServer := server.New(s)
	grpcServer := grpc.NewServer()

	ellepb.RegisterElleServer(grpcServer, elleServer)

	log.Printf("serving gRPC server at %s", *grpcAddr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to server elle", err)
	}
}