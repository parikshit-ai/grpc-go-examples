package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/parikshit-ai/go-proto/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Println("Greet function is invoked with req: ", req)
	firstName := req.GetRequest().GetFirstName()
	result := &greetpb.GreetResponse{
		Result: firstName,
	}
	return result, nil
}
func (*server) GreetManyTimes(req *greetpb.GreetManyTimeRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		result := "Hello " + firstName + " " + strconv.Itoa(i) + " Times."
		res := &greetpb.GreetManyTimeResponse{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(time.Second)
	}
	return nil
}
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
