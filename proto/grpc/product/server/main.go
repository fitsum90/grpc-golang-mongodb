// server/main.go
package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	pb "github.com/VENOLD/grpc/grpc/product"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ProductServer is the server type for handling Product service requests
type ProductServer struct {
	pb.UnimplementedProductServer
	mongoClient *mongo.Client
}

// AddProduct is the implementation of the AddProduct RPC method
func (s *ProductServer) AddProduct(ctx context.Context, req *pb.ProductRequest) (*pb.ProductResponse, error) {
	// Insert product data into MongoDB
	collection := s.mongoClient.Database("Product-grpc").Collection("products")
	_, err := collection.InsertOne(ctx, req)
	if err != nil {
		log.Printf("Error inserting into MongoDB: %v", err)
		return nil, err
	}

	result := "Product added successfully"
	return &pb.ProductResponse{Result: result}, nil
}

func main() {
	// Connect to MongoDB
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://127.0.0.1"))
	log.Println("Connected to MongoDB")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()



	// Create a gRPC server
	server := grpc.NewServer()

	// Register the ProductServer with the gRPC server
	pb.RegisterProductServer(server, &ProductServer{mongoClient: mongoClient})

	// Register reflection service on gRPC server.
	reflection.Register(server)

	// Listen on a port
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("gRPC server is running on :50051")
	// Serve the gRPC server
	if err := server.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
