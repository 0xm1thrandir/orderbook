
module github.com/0xm1thrandir/orderbook

go 1.20

require github.com/gorilla/websocket v1.5.0 // indirect
require github.com/0xm1thdandir/orderbook/pkg/binance v0.0.0

replace github.com/0xm1thdandir/orderbook/pkg/binance => ./pkg/binance

