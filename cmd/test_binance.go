package main

import (
	"encoding/json"
	"fmt"
	"log"
	"github.com/0xm1thdandir/orderbook/pkg/binance"
)

func main() {
	symbol := "btcusdt"
	limit := 5

	// Get initial order book data
	orderBook, err := binance.GetOrderBook(symbol, limit)
	if err != nil {
		log.Fatal("Error fetching order book:", err)
	}

	fmt.Println("Initial order book:", orderBook)

	// Handle depth stream messages
	handleDepthMessage := func(data []byte) {
		// Parse and process the depth message
		fmt.Println("Depth message:", string(data))
	}

	// Handle book ticker stream messages
	handleBookTickerMessage := func(data []byte) {
		// Parse and process the book ticker message
		fmt.Println("Book ticker message:", string(data))
	}

	// Connect to WebSocket streams
	go func() {
		if err := binance.ConnectDepthStream(symbol, handleDepthMessage); err != nil {
			log.Fatal("Error connecting to depth stream:", err)
		}
	}()

	go func() {
		if err := binance.ConnectBookTickerStream(symbol, handleBookTickerMessage); err != nil {
			log.Fatal("Error connecting to book ticker stream:", err)
		}
	}()

	// Keep the main function running
	select {}
}

