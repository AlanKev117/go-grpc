package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"github.com/AlanKev117/go-grpc/greet/greetpb"
	"google.golang.org/grpc"
)

// server defines the behaviour behind the grpc server
type server struct{}

// Greet returns a response that includes the names provided by the request req.
// It needs a context as the first argument to work.
// In case of error, the second value returned will be different to nil.
func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {

	fmt.Printf("Greet called with %v\n", req)

	firstName := req.GetGreeting().GetFirstName()
	secondName := req.GetGreeting().GetSecondName()
	resultString := fmt.Sprintf("Hello, %v %v", firstName, secondName)
	result := &greetpb.GreetResponse{
		Result: resultString,
	}

	return result, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {

	fmt.Printf("GreetManyTimes called with: %v\n", req)

	firstName := req.GetGreeting().GetFirstName()
	secondName := req.GetGreeting().GetSecondName()

	for i := 0; i < 10; i++ {
		res_string := fmt.Sprintf("Hello %v, %v %v", i, firstName, secondName)
		res := &greetpb.GreetManyTimesResponse{
			Result: res_string,
		}
		stream.Send(res)
		time.Sleep(200 * time.Millisecond)
	}

	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Println("LongGreet called with client streaming request...")
	reqArgs := []string{}
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			res := &greetpb.LongGreetResponse{
				Result: fmt.Sprintf("Hello to all of you, %s", strings.Join(reqArgs, ", ")),
			}
			return stream.SendAndClose(res)
		}
		if err != nil {
			log.Fatalf("error while reading from client stream: %v", err)
		}
		firstName := req.GetGreeting().GetFirstName()
		secondName := req.GetGreeting().GetSecondName()
		reqArgs = append(reqArgs, fmt.Sprintf("%v %v", firstName, secondName))
	}
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
