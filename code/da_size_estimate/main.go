package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	Url = "https://mainnet.base.org"
	// Address = "0x0cd73c6191906b6f5795efd525f77e65d6aa7561"
	Address = "0xfBd85a0c200b286Ef2D7a08306113669013C39dA"
)

func main() {
	client, err := ethclient.Dial(Url)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	address := common.HexToAddress(Address)

	latestBlock, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Fatalf("Failed to get the latest block: %v", err)
	}

	startBlock := latestBlock - 1000
	txCount := 0
	for blockNum := latestBlock; blockNum > startBlock && txCount < 100; blockNum-- {
		block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(blockNum)))
		if err != nil {
			log.Fatalf("Failed to get block %d: %v", blockNum, err)
		}

		for _, tx := range block.Transactions() {
			from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
			if err != nil {
				log.Fatalf("Failed to get sender: %v", err)
			}

			if from == address {
				fmt.Printf("Block %d: Transaction hash: %s, tx size: %d, compressed size: %s\n", blockNum, tx.Hash().Hex(), tx.Size(), tx.RollupCostData().EstimatedDASize())
				txCount++
				if txCount >= 100 {
					break
				}
			}
		}
	}
}
