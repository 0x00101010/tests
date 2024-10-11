package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// Connect to an Ethereum node
	client, err := ethclient.Dial("http://98.82.123.63:6000")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	// Start from block 16364870
	blockNumber := big.NewInt(16364870)

	for {
		// Get the block
		block, err := client.BlockByNumber(context.Background(), blockNumber)
		if err != nil {
			log.Fatalf("Failed to fetch block %d: %v", blockNumber, err)
		}

		fmt.Printf("Checking block %d\n", blockNumber)

		// Check each transaction in the block
		for _, tx := range block.Transactions() {
			if tx.To() == nil {
				fmt.Printf("Found a contract creation transaction in block %d\n", blockNumber)
				fmt.Printf("Transaction hash: %s\n", tx.Hash().Hex())
				printTxDetails(tx)
				return
			}
		}

		// Move to the previous block
		blockNumber.Sub(blockNumber, big.NewInt(1))

		// Safety check to avoid going too far back
		if blockNumber.Cmp(big.NewInt(0)) <= 0 {
			fmt.Println("Reached genesis block without finding a contract creation transaction")
			return
		}
	}
}

func printTxDetails(tx *types.Transaction) {
	fmt.Printf("From: %s\n", getFromAddress(tx))
	fmt.Printf("Value: %s\n", tx.Value().String())
	fmt.Printf("Gas: %d\n", tx.Gas())
	fmt.Printf("Gas Price: %s\n", tx.GasPrice().String())
	fmt.Printf("Nonce: %d\n", tx.Nonce())
	fmt.Printf("Data: 0x%x\n", tx.Data())
}

func getFromAddress(tx *types.Transaction) string {
	from, err := types.Sender(types.NewEIP155Signer(tx.ChainId()), tx)
	if err != nil {
		return "Unknown"
	}
	return from.Hex()
}
