package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/AlanKev117/go-grpc/calculator/calculatorpb"
	"google.golang.org/grpc"
)

// server defines the behaviour behind the grpc server
type server struct{}

// Greet returns a response that includes the names provided by the request req.
// It needs a context as the first argument to work.
// In case of error, the second value returned will be different to nil.
func (*server) Calculate(ctx context.Context, req *calculatorpb.OperationRequest) (*calculatorpb.OperationResponse, error) {

	fmt.Printf("Calculate function invoked with %v\n", req)

	operation := req.GetOperationArgs().GetOperation()
	value1 := req.GetOperationArgs().GetValue1()
	value2 := req.GetOperationArgs().GetValue2()

	var operationResult float32

	switch operation {
	case calculatorpb.Operation_OPCODE_SUM:
		operationResult = value1 + value2
	case calculatorpb.Operation_OPCODE_SUB:
		operationResult = value1 - value2
	case calculatorpb.Operation_OPCODE_MUL:
		operationResult = value1 * value2
	case calculatorpb.Operation_OPCODE_DIV:
		operationResult = value1 / value2
	}

	result := &calculatorpb.OperationResponse{
		Result: operationResult,
	}

	return result, nil
}

func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	number := req.GetNumber()
	prime := uint32(2)
	for number > 1 {
		if number%prime == 0 {
			stream.Send(&calculatorpb.PrimeNumberDecompositionResponse{
				Prime: prime,
			})
			number /= prime
		} else {
			prime++
		}
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
