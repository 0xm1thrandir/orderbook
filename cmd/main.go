package main

import (
	"log"
        "github.com/0xm1thrandir/orderbook/pkg/binance"
        "github.com/0xm1thrandir/orderbook/pkg/grpc"
)

func main() {
	// Initialize the order book and gRPC server
	orderBook := &binance.OrderBook{}
	grpcServer := grpc.NewServer()

	// Start the Binance WebSocket connections
	err := binance.StartWebSocket("btcusdt") // Use the appropriate trading pair
	if err != nil {
		log.Fatalf("Failed to start WebSocket: %v", err)
	}


	// Send the aggregated midpoint to gRPC clients whenever the order book is updated
	updateMidpoint := func() {
		midpoint := orderBook.Midpoint()
		grpcServer.SendMidpoint(midpoint)
	}


	// Register the updateMidpoint function as a callback for WebSocket updates
	binance.RegisterUpdateCallback(updateMidpoint)

	// Start the gRPC server
	grpcServer.Start()
}

