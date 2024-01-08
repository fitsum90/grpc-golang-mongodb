// client/main.go
package main

import (
	"context"
	"log"
	"google.golang.org/grpc"
	pb "github.com/VENOLD/grpc/grpc/product"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	// Create a gRPC client
	client := pb.NewProductClient(conn)

	// Create a sample ProductRequest
	request := &pb.ProductRequest{
		ProductName:        "SampleProduct",
		ProductDescription: "This is a sample product",
		ProductPrice:       "19.99",
		// Add any other fields as needed
	}

	// Call the AddProduct method
	response, err := client.AddProduct(context.Background(), request)
	if err != nil {
		log.Fatalf("could not add product: %v", err)
	}

	// Print the response
	log.Printf("Server Response: %s", response.Result)
}
