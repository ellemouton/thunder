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

var templates = template.Must(template.ParseGlob("assets/*.html"))

func main() {
	s, err := newState()
	if err != nil {
		log.Fatalf("newState: %s", err)
	}

	//ops.StartLoops(s)

	r := newRouter(s)

	go serverGRPCForever(s)

	log.Println("Serving on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

var grpcAddr = flag.String("grpc_address", ":8082", "elle grpc address")

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