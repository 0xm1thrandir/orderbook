package grpc

import (
    "context"
    "net"
    "sync"

    "google.golang.org/grpc"
)

type Client struct {
    stream MidpointService_GetMidpointServer
}

type CustomMidpointServer struct {
    MidpointServiceServer
    clients []*Client
    mtx     sync.Mutex
}

func (s *CustomMidpointServer) GetMidpoint(req *MidpointRequest, stream MidpointService_GetMidpointServer) error {
    client := &Client{stream: stream}
    s.mtx.Lock()
    s.clients = append(s.clients, client)
    s.mtx.Unlock()

    // Keep the stream open for updates
    <-stream.Context().Done()

    return nil
}

func (s *CustomMidpointServer) SendMidpoint(midpoint float64) {
    s.mtx.Lock()
    for _, client := range s.clients {
        client.stream.Send(&MidpointResponse{Midpoint: midpoint})
    }
    s.mtx.Unlock()
}

func StartServer(address string) error {
    listener, err := net.Listen("tcp", address)
    if err != nil {
        return err
    }

    server := grpc.NewServer()
    customServer := &CustomMidpointServer{}
    RegisterMidpointServiceServer(server, customServer)

    return server.Serve(listener)
}
