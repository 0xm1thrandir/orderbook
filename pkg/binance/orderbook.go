package binance

import (
	"strconv"
)


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

