package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	// Url = "https://mainnet.base.org"
	Url = "https://c3-chainproxy-base-mainnet-archive-dev.cbhq.net:8545"

	Hour       = 60 * 60 / 2
	Day        = 24 * 60 * 60 / 2
	Week       = Day * 7
	StartBlock = 22374553
	EndBlock   = StartBlock - Week
)

var (
	BundlerAddresses = map[common.Address]any{
		common.HexToAddress("0x0cd73c6191906b6f5795efd525f77e65d6aa7561"): nil, // alchemy
		common.HexToAddress("0xfBd85a0c200b286Ef2D7a08306113669013C39dA"): nil, // cb 1
		common.HexToAddress("0x1984c070e64e561631A7E20Ea3c4CbF4eb198Da8"): nil, // cb 2
		common.HexToAddress("0xC4a4e8Ae10B82a954519cA2EcC9EFC8f77819E86"): nil, // cb 3
		common.HexToAddress("0x6d10c567DB15b40Bfb1A162c16cbd7a3e80bB12b"): nil, // cb 4
	}
)

func main() {
	client, err := ethclient.Dial(Url)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	allTxFile, err := os.Create("all_txs.csv")
	if err != nil {
		log.Fatalf("Failed to create all_txs.csv: %v", err)
	}
	defer allTxFile.Close()

	allTxWriter := csv.NewWriter(allTxFile)
	allTxWriter.Write([]string{"block_number", "tx_size", "compressed_size"})
	defer allTxWriter.Flush()

	alchemyBundlerTxFile, err := os.Create("alchemy_bundler_txs.csv")
	if err != nil {
		log.Fatalf("Failed to create alchemy_bundler_txs.csv: %v", err)
	}
	defer alchemyBundlerTxFile.Close()

	alchemyBundlerTxWriter := csv.NewWriter(alchemyBundlerTxFile)
	alchemyBundlerTxWriter.Write([]string{"block_number", "tx_size", "compressed_size"})
	defer alchemyBundlerTxWriter.Flush()

	cbBundlerTxFile, err := os.Create("cb_bundler_txs.csv")
	if err != nil {
		log.Fatalf("Failed to create cb_bundler_txs.csv: %v", err)
	}
	defer cbBundlerTxFile.Close()

	cbBundlerTxWriter := csv.NewWriter(cbBundlerTxFile)
	cbBundlerTxWriter.Write([]string{"block_number", "tx_size", "compressed_size"})
	defer cbBundlerTxWriter.Flush()

	bundlerTxFile, err := os.Create("bundler_txs.csv")
	if err != nil {
		log.Fatalf("Failed to create bundler_txs.csv: %v", err)
	}
	defer bundlerTxFile.Close()

	bundlerTxWriter := csv.NewWriter(bundlerTxFile)
	bundlerTxWriter.Write([]string{"block_number", "tx_size", "compressed_size"})
	defer bundlerTxWriter.Flush()

	signer := types.LatestSignerForChainID(big.NewInt(8453))
	for blockNum := StartBlock; blockNum > EndBlock; blockNum-- {
		if (StartBlock-blockNum)%10 == 0 {
			fmt.Printf("Processing block %6d, remaining: %6d\n", StartBlock-blockNum, blockNum-EndBlock)
		}

		block, err := getBlock(client, blockNum)
		if err != nil {
			log.Printf("Failed to get block %d: %v", blockNum, err)
			continue
		}

		for _, tx := range block.Transactions() {
			txSize := tx.Size()
			daSize := tx.RollupCostData().EstimatedDASize()

			allTxWriter.Write([]string{
				fmt.Sprintf("%d", blockNum),
				fmt.Sprintf("%d", txSize),
				fmt.Sprintf("%d", daSize),
			})

			from, err := types.Sender(signer, tx)
			if err != nil {
				log.Fatalf("Failed to get sender for tx %v: %v", tx.Hash(), err)
			}

			if _, exists := BundlerAddresses[from]; exists {
				bundlerTxWriter.Write([]string{
					fmt.Sprintf("%d", blockNum),
					fmt.Sprintf("%d", txSize),
					fmt.Sprintf("%d", daSize),
				})

				if from == common.HexToAddress("0x0cd73c6191906b6f5795efd525f77e65d6aa7561") {
					alchemyBundlerTxWriter.Write([]string{
						fmt.Sprintf("%d", blockNum),
						fmt.Sprintf("%d", txSize),
						fmt.Sprintf("%d", daSize),
					})
				} else {
					cbBundlerTxWriter.Write([]string{
						fmt.Sprintf("%d", blockNum),
						fmt.Sprintf("%d", txSize),
						fmt.Sprintf("%d", daSize),
					})
				}
			}
		}
	}
}

// retryWithBackoff attempts to execute the given function with exponential backoff
func retryWithBackoff(attempts int, initialDelay time.Duration, maxDelay time.Duration, operation func() error) error {
	var err error
	delay := initialDelay

	for i := 0; i < attempts; i++ {
		err = operation()
		if err == nil {
			return nil
		}

		if i < attempts-1 { // don't sleep after last attempt
			log.Printf("Attempt %d failed, retrying after %v: %v", i+1, delay, err)
			time.Sleep(delay)
			// Double the delay for next attempt, but don't exceed maxDelay
			delay *= 2
			if delay > maxDelay {
				delay = maxDelay
			}
		}
	}

	return fmt.Errorf("failed after %d attempts: %v", attempts, err)
}

func getBlock(client *ethclient.Client, blockNum int) (*types.Block, error) {
	var block *types.Block

	operation := func() error {
		var err error
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		block, err = client.BlockByNumber(ctx, big.NewInt(int64(blockNum)))
		if err != nil {
			return fmt.Errorf("failed to get block %d: %v", blockNum, err)
		}
		return nil
	}

	err := retryWithBackoff(3, 1*time.Second, 5*time.Second, operation)
	if err != nil {
		return nil, err
	}

	return block, nil
}
