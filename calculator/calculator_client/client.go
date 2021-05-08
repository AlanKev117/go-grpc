package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/AlanKev117/go-grpc/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func doOperation(c calculatorpb.CalculatorServiceClient, opcode calculatorpb.Operation, value1 float32, value2 float32) {

	req := &calculatorpb.OperationRequest{
		OperationArgs: &calculatorpb.OperationArgs{
			Operation: opcode,
			Value1:    value1,
			Value2:    12.0,
		},
	}

	res, err := c.Calculate(context.Background(), req)

	if err != nil {
		log.Fatalf("Error while calling Calculate from server: %v", err)
	}
	log.Printf("Response from server Calculate: %v", res.Result)
}

func doGetPrimeFactors(c calculatorpb.CalculatorServiceClient, number uint32) {
	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: number,
	}

	resStream, err := c.PrimeNumberDecomposition(context.Background(), req)

	if err != nil {
		log.Fatalf("error while calling PrimerNumberDecomposition: %v", err)
	}

	primes := []uint32{}

	for {
		res, err := resStream.Recv()
		if err == io.EOF {
			log.Printf("No more prime factors left for %v\n", number)
			break
		}
		if err != nil {
			log.Fatalf("Error while receiving next prime factor: %v\n", err)
		}
		primes = append(primes, res.GetPrime())
	}

	fmt.Printf("Prime factors of %v are: \n", number)
	fmt.Printf("%v\n", primes)
}

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Couldn't connect: %v", err)
	}
	defer conn.Close()

	c := calculatorpb.NewCalculatorServiceClient(conn)

	doOperation(c, calculatorpb.Operation_OPCODE_SUM, 23, 54)
	doGetPrimeFactors(c, 1)
}
