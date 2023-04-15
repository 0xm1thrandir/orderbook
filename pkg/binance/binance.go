package binance

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)


// UpdateCallback is a function that will be called when there's an update to the order book
type UpdateCallback func()

var updateCallback UpdateCallback

// RegisterUpdateCallback registers a callback function that will be called when there's an update to the order book
func RegisterUpdateCallback(callback UpdateCallback) {
	updateCallback = callback
}



const BaseURL = "https://api.binance.com"

type OrderBook struct {
	Bids [][2]string `json:"bids"`
	Asks [][2]string `json:"asks"`
}


func StartWebSocket(symbol string) error {
	// Start the depth stream connection
	err := ConnectDepthStream(symbol, func(message []byte) {
		// Unmarshal the depth stream message and update the local order book
		// You need to implement the logic to update the local order book here

		// Notify the main program when there's an update to the order book
		if updateCallback != nil {
			updateCallback()
		}
	})
	if err != nil {
		return err
	}

	// Start the book ticker stream connection
	err = ConnectBookTickerStream(symbol, func(message []byte) {
		// Unmarshal the book ticker stream message and update the local order book
		// You need to implement the logic to update the local order book here

		// Notify the main program when there's an update to the order book
		if updateCallback != nil {
			updateCallback()
		}
	})
	if err != nil {
		return err
	}

	return nil
}



func GetOrderBook(symbol string, limit int) (*OrderBook, error) {
	url := fmt.Sprintf("%s/api/v3/depth?symbol=%s&limit=%d", BaseURL, symbol, limit)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var orderBook OrderBook
	if err := json.Unmarshal(body, &orderBook); err != nil {
		return nil, err
	}

	return &orderBook, nil
}

var wsDialer = websocket.DefaultDialer

func ConnectDepthStream(symbol string, handleMessage func([]byte)) error {
	url := fmt.Sprintf("wss://stream.binance.com:9443/ws/%s@depth", symbol)
	conn, _, err := wsDialer.Dial(url, nil)
	if err != nil {
		return err
	}

	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		handleMessage(message)
	}

	return nil
}

func ConnectBookTickerStream(symbol string, handleMessage func([]byte)) error {
	url := fmt.Sprintf("wss://stream.binance.com:9443/ws/%s@bookTicker", symbol)
	conn, _, err := wsDialer.Dial(url, nil)
	if err != nil {
		return err
	}

	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		handleMessage(message)
	}

	return nil
}

