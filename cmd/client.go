package main

import (
	"github.com/0xm1thrandir/orderbook/pkg/grpc"
)

func main() {
	grpc.StartClient("localhost:50051")
}

