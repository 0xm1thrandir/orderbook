package grpc

import (
	"log"
	"net"
	"sync"

	"github.com/0xm1thrandir/orderbook/pkg/grpc/pb"
	"google.golang.org/grpc"
)

type CustomMidpointServer struct {
	pb.UnimplementedMidpointServiceServer
	mu       sync.Mutex
	clients  map[string]pb.MidpointService_MidpointServer
	nextID   int
}

func NewServer() *CustomMidpointServer {
	return &CustomMidpointServer{
		clients: make(map[string]pb.MidpointService_MidpointServer),
		nextID:  1,
	}
}

func (s *CustomMidpointServer) GetMidpoint(req *pb.MidpointRequest, stream pb.MidpointService_MidpointServer) error {
	s.mu.Lock()
	id := s.nextID
	s.nextID++
	s.clients[strconv.Itoa(id)] = stream
	s.mu.Unlock()

	// This will block the stream, waiting for it to be closed by the client.
	<-stream.Context().Done()

	s.mu.Lock()
	delete(s.clients, strconv.Itoa(id))
	s.mu.Unlock()

	return nil
}

func (s *CustomMidpointServer) SendMidpoint(midpoint float64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	response := &pb.MidpointResponse{Midpoint: midpoint}
	for _, client := range s.clients {
		if err := client.Send(response); err != nil {
			log.Printf("Failed to send midpoint to client: %v", err)
		}
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

