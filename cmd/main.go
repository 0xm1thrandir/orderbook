package main

import (
    "log"
    "github.com/0xm1thrandir/orderbook/pkg/binance"
    "github.com/0xm1thrandir/orderbook/pkg/grpc"
)

func main() {
    // Initialize the order book and gRPC server
    orderBook := &binance.OrderBook{}
    grpcServer := grpc.NewCustomMidpointServer()

    // Start the Binance WebSocket connections
    go func() {
        if err := binance.ConnectDepthStream(orderBook, "btcusdt", func(message []byte) {
            // Log the received depth stream message
            log.Printf("Depth stream message: %s", string(message))

            // Update the local order book
            orderBook.UpdateFromDepthStreamMessage(message)
        }); err != nil {
            log.Fatal(err)
        }
    }()

    // Send the aggregated midpoint to gRPC clients whenever the order book is updated
    updateMidpoint := func() {
        midpoint := orderBook.Midpoint()
        log.Printf("Updated midpoint: %f", midpoint) // Print the updated midpoint
        grpcServer.SendMidpoint(midpoint)
    }

    // Register the updateMidpoint function as a callback for WebSocket updates
    orderBook.RegisterUpdateCallback(updateMidpoint)

    // Start the gRPC server
    go grpcServer.Start("localhost:50051") // Pass the address to the Start method

    // Start the gRPC client
    grpc.StartClient("localhost:50051")
}
