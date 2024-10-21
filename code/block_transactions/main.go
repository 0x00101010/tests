package main

import (
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// Connect to an Ethereum node
	client, err := ethclient.Dial("http://98.82.123.63:6000")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
}
