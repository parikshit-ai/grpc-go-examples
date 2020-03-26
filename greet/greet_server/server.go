package main

import (
	"context"
	"fmt"
	"io"
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
func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Println("Long greet function is invoke")
	result := "Hello "
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// finished reading client stream
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalln("Error while reading client stream Err :", err)
		}
		firstName := req.GetGreeting().GetFirstName()
		result += firstName + "! "
	}
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Println("Greet everyone from server")
	result := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// done with the request send response
			return nil
		}
		if err != nil {
			fmt.Println("Error while reciving from  server err : ", err)
		}
		firstname := req.GetGreeting().GetFirstName()
		result += "Hello " + firstname + " |"
		sendErr := stream.Send(&greetpb.GreetEveryoneResponse{
			Result: result,
		})
		if sendErr != nil {
			log.Fatalln("Error while reciving the data from the server Err ", sendErr)
			return sendErr
		}
	}
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
