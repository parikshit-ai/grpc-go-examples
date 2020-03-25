package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/parikshit-ai/go-proto/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Add(ctx context.Context, in *calculatorpb.Request) (*calculatorpb.Response, error) {
	fmt.Println("inside Add request is", in)
	out := &calculatorpb.Response{
		Ans: in.A + in.B,
	}
	return out, nil
}
func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalln("Error while listening err: ", err)
	}
	s := grpc.NewServer()
	calculatorpb.RegisterCalculateServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalln("Error occur while serving err: ", err)
	}
	fmt.Print("Server is running ")
}
