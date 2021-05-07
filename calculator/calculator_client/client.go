package main

import (
	"context"
	"log"

	"github.com/AlanKev117/go-grpc/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Couldn't connect: %v", err)
	}
	defer conn.Close()

	c := calculatorpb.NewCalculatorServiceClient(conn)
	req := &calculatorpb.OperationRequest{
		OperationArgs: &calculatorpb.OperationArgs{
			Operation: calculatorpb.Operation_OPCODE_DIV,
			Value1:    12.0,
			Value2:    12.0,
		},
	}

	res, err := c.Calculate(context.Background(), req)

	if err != nil {
		log.Fatalf("Error while calling Calculate from server: %v", err)
	}
	log.Printf("Response from server Calculate: %v", res.Result)

}
