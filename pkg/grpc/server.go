package grpc
import (
	"log"
	"net"
	"strconv"
	"sync"

	pb "github.com/0xm1thrandir/orderbook/pkg/grpc/pb"
	grpcLib "google.golang.org/grpc"
)

type CustomMidpointServer struct {
	pb.UnimplementedMidpointServiceServer
	clients   map[string]pb.MidpointService_MidpointServer
	clientMtx sync.Mutex
}

func (s *CustomMidpointServer) GetMidpoint(_ *pb.MidpointRequest, stream pb.MidpointService_MidpointServer) error {
	clientID := strconv.Itoa(len(s.clients))
	s.clientMtx.Lock()
	s.clients[clientID] = stream
	s.clientMtx.Unlock()

	// Keep the connection open.
	<-stream.Context().Done()

	// Remove the client from the list when the connection is closed.
	s.clientMtx.Lock()
	delete(s.clients, clientID)
	s.clientMtx.Unlock()

	return nil
}

func (s *CustomMidpointServer) SendMidpoint(midpoint float64) {
	s.clientMtx.Lock()
	defer s.clientMtx.Unlock()

	midpointMessage := &pb.MidpointResponse{Midpoint: midpoint}
	for _, client := range s.clients {
		client.Send(midpointMessage)
	}
}

func NewServer() *CustomMidpointServer {
	return &CustomMidpointServer{
		clients: make(map[string]pb.MidpointService_MidpointServer),
	}
}

func StartServer(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	server := grpcLib.NewServer()
	pb.RegisterMidpointServiceServer(server, NewServer())

	log.Printf("gRPC server started on %s", address)
	return server.Serve(listener)
}
