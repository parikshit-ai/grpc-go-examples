package main

import (
	"log"
	"net"

	"github.com/parikshit-ai/go-proto/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct{}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalln("Faild to listen ", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalln("Failed to serve ", err)
	}

}
