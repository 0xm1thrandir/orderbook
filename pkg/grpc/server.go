package grpc

import (
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
	pb "github.com/0xm1thrandir/orderbook/proto"
)

type CustomMidpointServer struct {
	pb.UnimplementedMidpointServiceServer
	mu          sync.Mutex
	clients     map[int]pb.MidpointService_GetMidpointServer
	nextClientID int
}

func (s *CustomMidpointServer) GetMidpoint(req *pb.MidpointRequest, stream pb.MidpointService_GetMidpointServer) error {
	s.mu.Lock()
	id := s.nextClientID
	s.nextClientID++
	s.clients[id] = stream
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.clients, id)
		s.mu.Unlock()
	}()

	<-stream.Context().Done()
	return nil
}

func (s *CustomMidpointServer) SendMidpoint(midpoint float64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, client := range s.clients {
		err := client.Send(&pb.MidpointResponse{Midpoint: midpoint})
		if err != nil {
			log.Printf("Failed to send midpoint to client: %v", err)
		}
	}
}

func NewServer() *CustomMidpointServer {
	return &CustomMidpointServer{
		clients: make(map[int]pb.MidpointService_GetMidpointServer),
	}
}

func StartServer(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	pb.RegisterMidpointServiceServer(server, NewServer())

	return server.Serve(listener)
}
