package grpc

import (
	"context"
	"log"
	"time"

	pb "github.com/0xm1thrandir/orderbook/proto"
	grpcLib "google.golang.org/grpc"
)

func StartClient(address string) {
	conn, err := grpcLib.Dial(address, grpcLib.WithInsecure(), grpcLib.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewMidpointServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	stream, err := client.GetMidpoint(ctx, &pb.MidpointRequest{})
	if err != nil {
		log.Fatalf("Failed to get midpoint: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err != nil {
			log.Printf("Error receiving midpoint: %v", err)
			break
		}
		log.Printf("Received midpoint: %f", res.GetMidpoint())
	}
}

