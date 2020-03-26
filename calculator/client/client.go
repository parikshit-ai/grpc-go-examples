package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

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
	// doUnary(c)
	// doServerStream(c)
	// doClientStream(c)
	doBiDirectional(c)
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

func doServerStream(c calculatorpb.CalculateClient) {
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
}
func doClientStream(c calculatorpb.CalculateClient) {
	fmt.Println("client started streaming")
	requests := []*calculatorpb.GetAvgRequest{
		&calculatorpb.GetAvgRequest{
			N: 3,
		},
		&calculatorpb.GetAvgRequest{
			N: 5,
		},
		&calculatorpb.GetAvgRequest{
			N: 6,
		},
		&calculatorpb.GetAvgRequest{
			N: 5,
		},
		&calculatorpb.GetAvgRequest{
			N: 77,
		},
	}
	stream, err := c.GetAvg(context.Background())
	if err != nil {
		fmt.Println("Error while getting getAvg stream Err:", err)
	}
	for _, req := range requests {
		fmt.Println("Request send via stream is ", req)
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Println("Error while getting the res", err)
	}
	fmt.Println("Got the response : ", res)
}
func doBiDirectional(c calculatorpb.CalculateClient) {
	stream, err := c.GetMax(context.Background())
	if err != nil {
		log.Fatalln("Error while creating stram Err", err)
	}
	requests := []*calculatorpb.GetMaxRequest{
		&calculatorpb.GetMaxRequest{
			N: 3,
		},
		&calculatorpb.GetMaxRequest{
			N: 5,
		},
		&calculatorpb.GetMaxRequest{
			N: 6,
		},
		&calculatorpb.GetMaxRequest{
			N: 5,
		},
		&calculatorpb.GetMaxRequest{
			N: 77,
		},
	}
	waitc := make(chan struct{}) //to block
	// to send the data
	go func() {
		for i, req := range requests {
			fmt.Println("Request send is ", req, i)
			stream.Send(req)
			time.Sleep(time.Second)
		}
		stream.CloseSend()
	}()
	// to recive the request
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("EOF")
				break
			}
			if err != nil {
				log.Fatalln("Error While reciving the data Err", err)
			}
			fmt.Println("Max no is ", res)
		}
		close(waitc)
	}()
	<-waitc
}
