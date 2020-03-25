package main

import (
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
}
