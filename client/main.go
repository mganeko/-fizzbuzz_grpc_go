//
// Package main implements a FizzBuzz server
//

package main

import (
	"context"
	"log"
	"io"
	"time"
	"google.golang.org/grpc"
	pb ".."
)

const (
	//port = ":50052"
	address     = "localhost:50052"
)


func callSingleFizzBuzz(client pb.FizzBuzzClient, ctx context.Context) {
	// -- Single FizzBuzz --
	for i := 1; i <= 20; i++ {
		result, err := client.SingleFizzBuzz(ctx, &pb.FizzBuzzRequest{X: int32(i)})
		if err != nil {
			log.Fatalf("call SingleFizzBuzz Error: %v", err)
		}

		log.Printf("%d --> %v\n", i, result.Result)
	}
}

func callLoopFizzBuzz(client pb.FizzBuzzClient, ctx context.Context) {
	request := &pb.FizzBuzzRequest{X: 30}
	stream, err := client.LoopFizzBuzz(ctx, request)
	if err != nil {
		log.Fatalf("call FizzBuzzRequest Error: %v", err)
	}
	for {
		result, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.FizzBuzzRequest(_) = _, %v", client, err)
		}
		log.Println(result.Result)
	}
}

func callMultiFizzBuzz(client pb.FizzBuzzClient, ctx context.Context) {
	stream, err := client.MultiRequestSingleResult(ctx)
	if err != nil {
		log.Fatalf("%v.MultiRequestSingleResult(_) = _, %v", client, err)
	}
	for i := 1; i <= 15; i++ {
		err := stream.Send(&pb.FizzBuzzRequest{X: int32(i)})
		if err != nil {
			log.Fatalf("%v.Send(%v) = %v", stream, i, err)
		}
	}
	reply, err2 := stream.CloseAndRecv()
	if err2 != nil {
		log.Fatalf("%v.CloseAndRecv() got error %v, want %v", stream, err2, nil)
	}
	log.Printf("multiRequestSingleResult result: %v", reply.Result)
}

// Multiple FizzBuzz
//	MultiFizzBuzz(ctx context.Context, opts ...grpc.CallOption) (FizzBuzz_MultiFizzBuzzClient, error)
func callMultFizzBuzzStream(client pb.FizzBuzzClient, ctx context.Context) {
	stream, err := client.MultiFizzBuzz(ctx)
	if err != nil {
		log.Fatalf("%v.MultFizzBuzz(_) = _, %v", client, err)
	}

	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a note : %v", err)
			}
			log.Println(in.Result)
		}
	}()

	for i := 1; i <= 50; i++ {
		err := stream.Send(&pb.FizzBuzzRequest{X: int32(i)})
		if err != nil {
			log.Fatalf("%v.Send(%v) = %v", stream, i, err)
		}
	}

	stream.CloseSend()
	<-waitc
}

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewFizzBuzzClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// -- Single FizzBuzz --
	callSingleFizzBuzz(client, ctx)

	// --- Loop FizzBuzz (1 to Request.x) ---
	log.Println("-- call LoopFizzBuzz --")
	callLoopFizzBuzz(client, ctx)

	// --- Multiple FizzBuzz ---
	log.Println("-- call MultiRequestSingleResult --")
	callMultiFizzBuzz(client, ctx)

	// --- Multiple FizzBuzz ---
	log.Println("-- call MultiFizzBuzz (stream, stream) --")
	callMultFizzBuzzStream(client, ctx)
}