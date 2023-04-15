package binance

type OrderBook struct {
   Bids []Order
   Asks []Order
}

type Order struct {
   Price  float64
   Amount float64
}

func (ob *OrderBook) UpdateBid(price, amount float64) {
   // Update the order book's bids
}

func (ob *OrderBook) UpdateAsk(price, amount float64) {
   // Update the order book's asks
}

func (ob *OrderBook) Midpoint() float64 {
   bestBid := ob.Bids[0].Price
   bestAsk := ob.Asks[0].Price
   return (bestBid + bestAsk) / 2
}

