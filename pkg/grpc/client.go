package grpc

import (
	"context"

	"google.golang.org/grpc"
)

func GetMidpoint(client MidpointClient) (float64, error) {
	req := &MidpointRequest{}
	resp, err := client.GetMidpoint(context.Background(), req)
	if err != nil {
		return 0, err
	}

	return resp.Midpoint, nil
}

func ConnectClient(address string) (MidpointClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	return NewMidpointClient(conn), nil
}

