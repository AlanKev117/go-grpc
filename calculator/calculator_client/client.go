package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/AlanKev117/go-grpc/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func doOperation(c calculatorpb.CalculatorServiceClient, opcode calculatorpb.Operation, value1 float32, value2 float32) {
	log.Printf("Executing Operation %v\n", opcode)
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
	fmt.Printf("Response from server Calculate: %v", res.Result)
}

func doGetPrimeFactors(c calculatorpb.CalculatorServiceClient, number uint32) {
	log.Printf("Calculating prime factors for %v", number)
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

func doCalculateAverage(c calculatorpb.CalculatorServiceClient, numbers []int32) {

	log.Printf("Calculating average for %v numbers", len(numbers))

	stream, err := c.ComputeAverage(context.Background())
	for _, number := range numbers {
		log.Printf("sending %v to server", number)
		stream.Send(&calculatorpb.ComputeAverageRequest{
			Number: number,
		})
		time.Sleep(100 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving average: %v", err)
	}

	average := res.GetAverage()

	fmt.Printf("Average for %v is: %v\n", numbers, average)
}

func doGetMaximumValues(c calculatorpb.CalculatorServiceClient, numbers []int32) {
	log.Printf("Calculating max values for %v numbers\n", len(numbers))

	stream, err := c.FindMaximum(context.Background())

	if err != nil {
		log.Fatalf("Error while getting stream from client: %v", err)
	}

	waitc := make(chan struct{})

	// Sending each value
	go func() {
		for _, number := range numbers {
			fmt.Printf("Sending %v to server.\n", number)
			stream.Send(&calculatorpb.FindMaximumRequest{
				Number: number,
			})
			time.Sleep(200 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	// Receiving and handling new max value
	go func() {
		defer close(waitc)
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving server response: %v", err)
			}
			nextMaximum := res.GetMaximum()
			fmt.Printf("Received new max value: %v\n", nextMaximum)
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

	c := calculatorpb.NewCalculatorServiceClient(conn)

	doOperation(c, calculatorpb.Operation_OPCODE_SUM, 23, 54)
	doGetPrimeFactors(c, 1)
	doCalculateAverage(c, []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	doGetMaximumValues(c, []int32{1, 5, 3, 6, 2, 20})
}
