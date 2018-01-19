package main

import (
	"net"
	"fmt"
	"log"
	"google.golang.org/grpc"
	"github.com/Sovianum/turbonetwork/nodeservice/pb"
	"github.com/Sovianum/turbonetwork/nodeservice/server"
	"context"
)

var port = 8082

func main() {
	lis, serverErr := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if serverErr != nil {
		log.Fatalf("failed to listen: %v", serverErr)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterNodeServiceServer(grpcServer, &server.Server{})
	go grpcServer.Serve(lis)

	conn, clientErr := grpc.Dial(fmt.Sprintf("127.0.0.1:%d", port), grpc.WithInsecure())
	if clientErr != nil {
		log.Fatal("Failed to connect")
	}

	client := pb.NewNodeServiceClient(conn)
	resp, err := client.CreateNodes(context.Background(), &pb.CreateRequest{})
	if err != nil {
		log.Fatalf("Failed to get response: %s", err.Error())
	}
	log.Printf("Succeeded %v", *resp)
}
