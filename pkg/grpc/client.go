package grpc

import (
	"context"
	"log"

	pb "github.com/0xm1thrandir/orderbook/proto"
	grpcLib "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

func StartClient(address string) {
	conn, err := grpcLib.Dial(address, grpcLib.WithInsecure(), 
grpcLib.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewMidpointServiceClient(conn)
	log.Println("Connected to server")
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	stream, err := client.GetMidpoint(ctx, &pb.MidpointRequest{})
	if err != nil {
		log.Fatalf("Failed to get midpoint: %v", err)
	}
        log.Println("Connected to the server") // Add this line

	for {
    res, err := stream.Recv()
    if err != nil {
        // Add this block to handle DeadlineExceeded errors
        if status.Code(err) == codes.DeadlineExceeded {
            log.Printf("Deadline exceeded, trying again")
            continue
        }
        log.Printf("Error receiving midpoint: %v", err)
        break
    }
    log.Printf("Received midpoint: %f", res.GetMidpoint())
}

}
