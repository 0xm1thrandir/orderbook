package grpc

import (
	"log"
	"net"
	"sync"
	"context"

	pb "github.com/0xm1thrandir/orderbook/proto"
	grpcLib "google.golang.org/grpc"
)

type CustomMidpointServer struct {
	pb.UnimplementedMidpointServiceServer
	clients      map[string]pb.MidpointService_GetMidpointServer
	clientMutex  sync.Mutex
	clientID     int
	updateSignal chan float64
}

func NewServer() *CustomMidpointServer {
	return &CustomMidpointServer{
		clients:      make(map[string]pb.MidpointService_GetMidpointServer),
		clientMutex:  sync.Mutex{},
		clientID:     0,
		updateSignal: make(chan float64),
	}
}

func (s *CustomMidpointServer) GetMidpoint(stream pb.MidpointService_GetMidpointServer) error {
    return &pb.MidpointResponse{
        Midpoint: 0.0, // Replace 0.0 with the actual midpoint value you want to return
    }, nil
}


func (s *CustomMidpointServer) SendMidpoint(midpoint float64) {
	s.clientMutex.Lock()
	defer s.clientMutex.Unlock()

	for _, client := range s.clients {
		res := &pb.MidpointResponse{
			Midpoint: midpoint,
		}
		client.Send(res)
	}
}

func (s *CustomMidpointServer) Start(address string) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpcLib.NewServer()
	pb.RegisterMidpointServiceServer(grpcServer, s)

	log.Printf("Starting gRPC server on %s", address)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
