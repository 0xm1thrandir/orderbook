package grpc

import (
	"log"
	"net"
	"strconv"
	"sync"

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

func NewCustomMidpointServer() *CustomMidpointServer {
	return &CustomMidpointServer{
		clients:      make(map[string]pb.MidpointService_GetMidpointServer),
		clientMutex:  sync.Mutex{},
		clientID:     0,
		updateSignal: make(chan float64),
	}
}

func (s *CustomMidpointServer) GetMidpoint(req *pb.MidpointRequest, stream pb.MidpointService_GetMidpointServer) error {
    // Existing log statement
    log.Println("Client connected")

    s.clientMutex.Lock()
    id := strconv.Itoa(s.clientID)
    s.clients[id] = stream
    s.clientID++
    s.clientMutex.Unlock()
    log.Printf("New client connected: %s", id)

    // Wait for updates and send them to the client
    for midpoint := range s.updateSignal {
        res := &pb.MidpointResponse{
            Midpoint: midpoint,
        }

        log.Printf("Sending midpoint to client %s: %f", id, res.GetMidpoint()) // Add this line
        if err := stream.Send(res); err != nil {
            log.Printf("Error sending midpoint to client %s: %v", id, err)
            s.clientMutex.Lock()
            delete(s.clients, id)
            s.clientMutex.Unlock()
            return err
        }
        log.Printf("Midpoint sent to client %s: %f", id, res.GetMidpoint()) // Add this line
    }

    return nil
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
