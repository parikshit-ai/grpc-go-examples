package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/parikshit-ai/go-proto/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalln("Unable to dial to localhost:50051 with err", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	fmt.Printf("server created %v", c)
	// doUnary(c)
	// doStreaming(c)
	doClientStreaming(c)
}
func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Started doing client streaming")
	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalln("Error while calling longGreet", err)
	}
	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Pari",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Pari1",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Pari2",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Pari3",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Pari4",
			},
		},
	}
	for _, req := range requests {
		fmt.Println("Sending stream current data is ", req)
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Println("Error while getting response from long greet Err :", err)
	}
	fmt.Printf("LONG GREET RESPONSE IS %+v", res)
}
func doStreaming(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetManyTimeRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Streaming Pari",
			LastName:  "Streaming Singh",
		},
	}
	fmt.Println("Server Streaming rpc")
	resStream, _ := c.GreetManyTimes(context.Background(), req)
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error while straming :", err)
		}
		fmt.Println("Response from GreetManyTImes ", msg.GetResult())
	}
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("calling doUnary in RPC..")
	req := &greetpb.GreetRequest{
		Request: &greetpb.Greeting{
			FirstName: "Pari",
			LastName:  "Singh",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalln("UnexpectEd Error while fetcing the data ", err)
	}
	fmt.Println(res)
}
