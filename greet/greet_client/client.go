package main

import (
	"context"
	"fmt"
	"log"

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
	doUnary(c)
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
		log.Fatalln("UNexpectEd Error while fetcing the data ", err)
	}
	fmt.Println(res)
}
