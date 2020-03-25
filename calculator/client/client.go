package main

import (
	"context"
	"fmt"
	"log"

	"github.com/parikshit-ai/go-proto/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalln("Error while dialing err: ", err)
	}
	defer cc.Close()
	c := calculatorpb.NewCalculateClient(cc)
	doUnary(c)
	fmt.Println("Client is running")
}

func doUnary(c calculatorpb.CalculateClient) {
	fmt.Println("inside dpUnary client")
	req := &calculatorpb.Request{
		A: 10,
		B: 30,
	}
	out, _ := c.Add(context.Background(), req)
	fmt.Println(out)
}
