package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/consensus/misc/eip4844"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
)

func main() {
	url := flag.String("url", "", "Ethereum RPC URL")
	flag.Parse()

	if *url == "" {
		log.Fatal("Please provide an Ethereum RPC URL using the -url flag")
	}

	client, err := ethclient.Dial(*url)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to get latest block header: %v", err)
	}

	latestBlock := header.Number.Uint64()
	// latestBlock = 20

	file, err := os.Create("blob_fees.csv")
	if err != nil {
		log.Fatalf("Failed to create CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{"Block Number", "Excess Blob Gas", "Gas Used", "Blob Fee"})
	if err != nil {
		log.Fatalf("Failed to write CSV header: %v", err)
	}

	for blockNum := uint64(0); blockNum <= latestBlock; blockNum++ {
		block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(blockNum)))
		if err != nil {
			log.Printf("Failed to get block %d: %v", blockNum, err)
			continue
		}

		block.BlobGasUsed()
		excessBlobGas := block.ExcessBlobGas()
		var blobFee *big.Int
		if excessBlobGas != nil {
			blobFee = eip4844.CalcBlobFee(*excessBlobGas)
		} else {
			blobFee = big.NewInt(0)
		}

		// blobNumber := excessBlobGasToBlobNumber(excessBlobGas)
		blobNum := *block.BlobGasUsed() / uint64(params.BlobTxBlobGasPerBlob)
		err = writer.Write([]string{
			strconv.FormatUint(blockNum, 10),
			formatUint64Ptr(excessBlobGas),
			formatUint64Ptr(block.BlobGasUsed()),
			formatUint64Ptr(&blobNum),
			blobFee.String(),
		})
		if err != nil {
			log.Printf("Failed to write row for block %d: %v", blockNum, err)
		}

		if blockNum%10 == 0 {
			fmt.Printf("Processed block %d\n", blockNum)
		}
	}

	writer.Flush()
	fmt.Println("Blob fee data has been written to blob_fees.csv")
}

func formatUint64Ptr(u *uint64) string {
	if u == nil {
		return "nil"
	}
	return strconv.FormatUint(*u, 10)
}
