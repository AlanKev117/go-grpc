package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/AlanKev117/go-grpc/greet/greetpb"
	"google.golang.org/grpc"
)

// server defines the behaviour behind the grpc server
type server struct{}

// Greet returns a response that includes the names provided by the request req.
// It needs a context as the first argument to work.
// In case of error, the second value returned will be different to nil.
func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {

	fmt.Printf("Greet function invoked with %v\n", req)

	firstName := req.GetGreeting().GetFirstName()
	secondName := req.GetGreeting().GetSecondName()
	resultString := fmt.Sprintf("Hello, %v %v", firstName, secondName)
	result := &greetpb.GreetResponse{
		Result: resultString,
	}

	return result, nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
