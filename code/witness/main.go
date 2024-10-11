package main

import (
	"context"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	rethURL := os.Getenv("RETH_ARCHIVAL_NODE_URL")
	gethURL := os.Getenv("GETH_ARCHIVAL_NODE_URL")

	ctx := context.Background()
	client, err := ethclient.Dial(gethURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	startBlock := os.Getenv("BLOCK_TO_WALK_BACK_FROM")
	if startBlock == "0" || startBlock == "" {
		block, err := client.BlockByNumber(ctx, big.NewInt(int64(rpc.LatestBlockNumber)))
		if err != nil {
			log.Fatalf("Failed to get latest block: %v", err)
		}
		startBlock = hexutil.EncodeBig(block.Number())
		log.Printf("Starting to walk back from latest block number: %s", startBlock)
	}

	log.Printf("RETH_ARCHIVAL_NODE_URL: %s\n", rethURL)
	log.Printf("GETH_ARCHIVAL_NODE_URL: %s\n", gethURL)

	// client.
}
