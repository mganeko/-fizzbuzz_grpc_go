//
// Package main implements a FizzBuzz server
//

package main

import (
	"context"
	"log"
	"net"
	"strconv"
	"io"
	"google.golang.org/grpc"
	//pb ".." // without go modules
	pb "fizzbuzz_go" // with go modules
)

const (
	port = ":50052"
)

type server struct {
	//pb.UnimplementedFizzbuzzServer //  undefined: fizzbuzz.UnimplementedFizzbuzzServer
	pb.FizzBuzzServer
}

func fizzbuzz(x int) string {
	var s string
	if x%15 == 0 {
		s = "FizzBuzz"
	} else if x%3 == 0 {
		s = "Fizz"
	} else if x%5 == 0 {
		s = "Buzz"
	} else {
		s = strconv.Itoa(x)
	}

	return s
}

// Single FizzBuzz
func (s *server) SingleFizzBuzz(ctx context.Context, in *pb.FizzBuzzRequest) (*pb.FizzBuzzReply, error) {
	x := in.GetX()
	str := fizzbuzz(int(x))
	log.Printf("SingleFizzBuzz Received: %v, return:%v", x, str)
	return &pb.FizzBuzzReply{Result: str}, nil
}

// Loop FizzBuzz (1 to Request.x)
func (s *server)LoopFizzBuzz(in *pb.FizzBuzzRequest, stream pb.FizzBuzz_LoopFizzBuzzServer) error {
	x := int(in.GetX())
	log.Printf("LoopFizzBuzz Received: %v", x)
	for i := 1; i <= x; i++ {
		str := fizzbuzz(i)
		log.Printf("%v --> %v", i, str)
		err := stream.Send(&pb.FizzBuzzReply{Result: str})
		if err != nil {
			return err
		}
	}
	return nil
}

// Multiple FizzBuzz
func (s *server)MultiRequestSingleResult(stream pb.FizzBuzz_MultiRequestSingleResultServer) error {
	var str string
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.FizzBuzzReply{Result: str})
		}
		if err != nil {
			return err
		}

		x := in.GetX()
		log.Printf("MultiRequestSingleResult Received: %v", x)
		str += fizzbuzz(int(x))
		str += " "
	}
}

// Multiple FizzBuzz
func (s *server) MultiFizzBuzz(stream pb.FizzBuzz_MultiFizzBuzzServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		x := in.GetX()
		str := fizzbuzz(int(x))
		log.Printf("MultiFizzBuzz Received: %v --> %v", x, str)

		err = stream.Send(&pb.FizzBuzzReply{Result: str})
		if err != nil {
			return err
		}
	}
}


// main
func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterFizzBuzzServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
