package main

import (
	"context"
	"fmt"
	"io"
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
	req := &calculatorpb.PrimeNoDecompositionRequest{
		N: 210,
	}
	resSream, _ := c.PrimeNoDecomposition(context.Background(), req)
	for {
		msg, err := resSream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln("someting went wrong 500 with err:", err)
		}
		fmt.Println(msg.GetN())
	}
	// doUnary(c)
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
