package binance

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
)

const BaseURL = "https://api.binance.com"

type OrderBook struct {
	Bids          [][2]string
	Asks          [][2]string
	updateCallback func()
	mu             sync.Mutex
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		Bids: make([][2]string, 0),
		Asks: make([][2]string, 0),
	}
}

func (ob *OrderBook) RegisterUpdateCallback(callback func()) {
	ob.mu.Lock()
	defer ob.mu.Unlock()
	ob.updateCallback = callback
}

func (ob *OrderBook) triggerUpdate() {
	ob.mu.Lock()
	defer ob.mu.Unlock()
	if ob.updateCallback != nil {
		ob.updateCallback()
	}
}

// Other functions and methods from your original binance.go file...


func StartWebSocket(symbol string) error {
    // Initialize the order book
    orderBook := NewOrderBook()

    // Start the depth stream connection
    err := ConnectDepthStream(orderBook, symbol, func(message []byte) {
        // Log the received depth stream message
        log.Printf("Depth stream message: %s", string(message))

        // Unmarshal the depth stream message and update the local order book
        err := orderBook.UpdateFromDepthStreamMessage(message)
        if err != nil {
            log.Println("Error updating order book from depth stream:", err)
        }

        // Notify the main program when there's an update to the order book
        orderBook.triggerUpdate()
    })
    if err != nil {
        return err
    }

    // Start the book ticker stream connection
    err = ConnectBookTickerStream(symbol, func(message []byte) {
        // Log the received book ticker stream message
        log.Printf("Book ticker stream message: %s", string(message))

        // Unmarshal the book ticker stream message and update the local order book
        err := orderBook.UpdateFromBookTickerStreamMessage(message)
        if err != nil {
            log.Println("Error updating order book from book ticker stream:", err)
        }

        // Notify the main program when there's an update to the order book
        orderBook.triggerUpdate()
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

	var rawOrderBook struct {
		Bids [][2]string `json:"bids"`
		Asks [][2]string `json:"asks"`
	}
	if err := json.Unmarshal(body, &rawOrderBook); err != nil {
		return nil, err
	}

	orderBook := &OrderBook{
		Bids: make([][2]string, len(rawOrderBook.Bids)),
                Asks: make([][2]string, len(rawOrderBook.Asks)),
	}

	for i, bid := range rawOrderBook.Bids {
		price, err := strconv.ParseFloat(bid[0], 64)
		if err != nil {
			return nil, err
		}
		amount, err := strconv.ParseFloat(bid[1], 64)
		if err != nil {
			return nil, err
		}
		orderBook.Bids[i] = [2]string{strconv.FormatFloat(price, 'f', -1, 64), strconv.FormatFloat(amount, 'f', -1, 64)}

	}

	for i, ask := range rawOrderBook.Asks {
		price, err := strconv.ParseFloat(ask[0], 64)
		if err != nil {
			return nil, err
		}
		amount, err := strconv.ParseFloat(ask[1], 64)
		if err != nil {
			return nil, err
		}
		orderBook.Asks[i] = [2]string{strconv.FormatFloat(price, 'f', -1, 64), strconv.FormatFloat(amount, 'f', -1, 64)}
	}

	return orderBook, nil
}





var wsDialer = websocket.DefaultDialer

func ConnectDepthStream(ob *OrderBook, symbol string, handleMessage func([]byte)) error {
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

func (ob *OrderBook) UpdateFromBookTickerStreamMessage(message []byte) error {
	var bookTickerStreamMessage struct {
		BestBidPrice  string `json:"b"`
		BestBidAmount string `json:"B"`
		BestAskPrice  string `json:"a"`
		BestAskAmount string `json:"A"`
	}

	if err := json.Unmarshal(message, &bookTickerStreamMessage); err != nil {
		return err
	}

	bestBidPrice, err := strconv.ParseFloat(bookTickerStreamMessage.BestBidPrice, 64)
	if err != nil {
		return err
	}
	bestBidAmount, err := strconv.ParseFloat(bookTickerStreamMessage.BestBidAmount, 64)
	if err != nil {
		return err
	}
	ob.UpdateBid(bestBidPrice, bestBidAmount)

	bestAskPrice, err := strconv.ParseFloat(bookTickerStreamMessage.BestAskPrice, 64)
	if err != nil {
		return err
	}
	bestAskAmount, err := strconv.ParseFloat(bookTickerStreamMessage.BestAskAmount, 64)
	if err != nil {
		return err
	}
	ob.UpdateAsk(bestAskPrice, bestAskAmount)

	return nil
}


func (ob *OrderBook) UpdateFromDepthStreamMessage(message []byte) error {
	var depthStreamMessage struct {
		Bids [][2]string `json:"b"`
		Asks [][2]string `json:"a"`
	}

	if err := json.Unmarshal(message, &depthStreamMessage); err != nil {
		return err
	}

	for _, bid := range depthStreamMessage.Bids {
		price, err := strconv.ParseFloat(bid[0], 64)
		if err != nil {
			return err
		}
		amount, err := strconv.ParseFloat(bid[1], 64)
		if err != nil {
			return err
		}
		ob.UpdateBid(price, amount)
	}

	for _, ask := range depthStreamMessage.Asks {
		price, err := strconv.ParseFloat(ask[0], 64)
		if err != nil {
			return err
		}
		amount, err := strconv.ParseFloat(ask[1], 64)
		if err != nil {
			return err
		}
		ob.UpdateAsk(price, amount)
	}

	return nil
}


func (ob *OrderBook) UpdateBid(price, amount float64) {
	order := [2]string{strconv.FormatFloat(price, 'f', -1, 64), strconv.FormatFloat(amount, 'f', -1, 64)}
	ob.Bids = append(ob.Bids, order)
}

func (ob *OrderBook) UpdateAsk(price, amount float64) {
	order := [2]string{strconv.FormatFloat(price, 'f', -1, 64), strconv.FormatFloat(amount, 'f', -1, 64)}
	ob.Asks = append(ob.Asks, order)
}

func (ob *OrderBook) Midpoint() float64 {
	bestBid, _ := strconv.ParseFloat(ob.Bids[0][0], 64)
	bestAsk, _ := strconv.ParseFloat(ob.Asks[0][0], 64)
	return (bestBid + bestAsk) / 2
}
