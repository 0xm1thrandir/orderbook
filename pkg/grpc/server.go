package grpc

import (
	"context"
	"net"

	"google.golang.org/grpc"
)

type MidpointServer struct {
	MidpointServer
}

func (s *MidpointServer) GetMidpoint(ctx context.Context, req *MidpointRequest) (*MidpointResponse, error) {
	// Calculate the aggregated midpoint and return it as a response.
	midpoint := 0.0 // Replace this with the actual calculation.
	return &MidpointResponse{Midpoint: midpoint}, nil
}

func StartServer(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	RegisterMidpointServer(server, &MidpointServer{})

	return server.Serve(listener)
}

