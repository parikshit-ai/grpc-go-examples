package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"time"

	"github.com/parikshit-ai/go-proto/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct{}

func (*server) Add(ctx context.Context, in *calculatorpb.Request) (*calculatorpb.Response, error) {
	fmt.Println("inside Add request is", in)
	out := &calculatorpb.Response{
		Ans: in.A + in.B,
	}
	return out, nil
}

func (*server) PrimeNoDecomposition(req *calculatorpb.PrimeNoDecompositionRequest, stream calculatorpb.Calculate_PrimeNoDecompositionServer) error {
	fmt.Println("started prime no decomposetion")
	var k int32 = 2
	n := req.GetN()
	for n > 1 {
		if n%k == 0 {
			res := &calculatorpb.PrimeNoDecompositionResponse{
				N: k,
			}
			n = n / k
			stream.Send(res)
			time.Sleep(time.Second)
		} else {
			k = k + 1
		}
	}
	return nil
}

func (*server) GetAvg(stream calculatorpb.Calculate_GetAvgServer) error {
	fmt.Println("Starting getavg from serverside")
	var ans float32
	var cnt int32
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// send the response
			ans /= float32(cnt)
			res := calculatorpb.GetAvgResponse{
				N: ans,
			}
			return stream.SendAndClose(&res)
		}
		if err != nil {
			log.Fatalln("Error while reciving the stream Err: ", err)
		}
		ans += float32(req.GetN())
		cnt++
	}
}

func (*server) GetMax(stream calculatorpb.Calculate_GetMaxServer) error {
	fmt.Println("Inside getMax serverside")
	var result int32 = -1
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// end from client side
			return nil
		}
		if err != nil {
			log.Fatalln("Error while reciving the data ", err)
		}
		if req.GetN() > result {
			errSend := stream.Send(&calculatorpb.GetMaxResponse{
				N: req.GetN(),
			})
			if errSend != nil {
				log.Fatalln("Error while sending the data to client")
				return errSend
			}
			result = req.GetN()
		}
	}
}

func (*server) GetSqureRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	fmt.Println("Inside server method of getSquareRoot")
	n := req.GetN()
	if n < 0 {
		// thorw an error
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Recived a negetive number, %v", n),
		)
	}
	return &calculatorpb.SquareRootResponse{
		N: math.Sqrt(float64(n)),
	}, nil
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
