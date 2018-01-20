package main

import (
	"context"
	"fmt"
	"github.com/Sovianum/turbonetwork/nodeservice/pb"
	"github.com/Sovianum/turbonetwork/nodeservice/server"
	"github.com/Sovianum/turbonetwork/nodeservice/server/factories"
	"google.golang.org/grpc"
	"log"
	"net"
)

var port = 8082

func main() {
	lis, serverErr := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if serverErr != nil {
		log.Fatalf("failed to listen: %v", serverErr)
	}

	grpcServer := grpc.NewServer()
	gteServer := server.NewGTEServer()

	pb.RegisterNodeServiceServer(grpcServer, gteServer)
	go grpcServer.Serve(lis)

	conn, clientErr := grpc.Dial(fmt.Sprintf("127.0.0.1:%d", port), grpc.WithInsecure())
	if clientErr != nil {
		log.Fatal("Failed to connect")
	}

	client := pb.NewNodeServiceClient(conn)

	createReq, _ := server.GetCreateRequest([]string{"node"}, []string{factories.PressureLossNodeType}, []map[string]float64{
		{"sigma": 1},
	})
	resp, err := client.CreateNodes(context.Background(), createReq)
	if err != nil {
		log.Fatalf("Failed to get response: %s", err.Error())
	}
	log.Printf("Succeeded %v", *resp)

	resp1, err1 := client.Process(context.Background(), &pb.Identifiers{
		Ids:[]*pb.NodeIdentifier{resp.Items[0].Identifiers[0]},
	})
	if err1 != nil {
		log.Fatalf("Failed to get response: %s", err1.Error())
	}
	log.Printf("Succeeded %v", *resp1)
}
