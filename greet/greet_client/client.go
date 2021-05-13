package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/AlanKev117/go-grpc/greet/greetpb"
	"google.golang.org/grpc"
)

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting unary gRPC...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName:  "Alan",
			SecondName: "Kevin",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Greet from server: %v", err)
	}
	log.Printf("Response from server Greet: %v", res.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting server streaming gRPC...")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName:  "Alan",
			SecondName: "Kevin",
		},
	}

	resStream, err := c.GreetManyTimes(context.Background(), req)

	if err != nil {
		log.Fatalf("Error while calling GreetManyTimes: %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while reading from stream: %v", err)
		}
		res := msg.GetResult()
		log.Printf("Response: %v\n", res)
	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting a client streaming gRPC operation...")

	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName:  "Alan",
				SecondName: "Kevin",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName:  "Dani",
				SecondName: "Elías",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName:  "Esteban",
				SecondName: "Damián",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())

	if err != nil {
		log.Fatalf("error while calling LongGreet from server: %v", err)
	}

	for i, req := range requests {
		log.Printf("Sending request via LongGreet stream %v", i)
		stream.Send(req)
		time.Sleep(200 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving LongGreet response: %v", err)
	}
	fmt.Println("Response from LongGreet: ")
	fmt.Println(res.GetResult())
}

func doBiDirectionalStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting a client streaming gRPC operation...")

	stream, err := c.GreetEveryone(context.Background())

	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
	}

	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName:  "Alan",
				SecondName: "Kevin",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName:  "Dani",
				SecondName: "Elías",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName:  "Esteban",
				SecondName: "Damián",
			},
		},
	}

	waitc := make(chan struct{})

	// Sending messages to the server
	go func() {
		for _, req := range requests {
			log.Printf("Sending message: %v\n", req)
			stream.Send(req)
			time.Sleep(200 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	// Receiving responses from the server
	go func() {
		defer close(waitc)
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving response from server: %v", err)
				break
			}
			fmt.Printf("Greeting received: %v\n", res.GetResult())
		}
	}()
	<-waitc
}

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Couldn't connect: %v", err)
	}
	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)

	doUnary(c)
	doServerStreaming(c)
	doClientStreaming(c)
	doBiDirectionalStreaming(c)
}
