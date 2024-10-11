package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	overwrite := os.Getenv("OVERWRITE")
	outputPath := os.Getenv("OUTPUT_PATH")
	absolutePath, exists := dirExistsRelative(outputPath)
	if exists {
		log.Printf("Output path: %s\n", absolutePath)
	} else {
		log.Fatalf("Directory %s does not exist", outputPath)
	}

	rethURL := os.Getenv("RETH_ARCHIVAL_NODE_URL")
	gethURL := os.Getenv("GETH_ARCHIVAL_NODE_URL")

	ctx := context.Background()
	geth, err := ethclient.Dial(gethURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	log.Printf("GETH_ARCHIVAL_NODE_URL: %s\n", gethURL)

	reth, err := ethclient.Dial(rethURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	log.Printf("RETH_ARCHIVAL_NODE_URL: %s\n", rethURL)

	startBlock := os.Getenv("BLOCK_TO_WALK_BACK_FROM")
	if startBlock == "0" || startBlock == "" {
		block, err := geth.BlockByNumber(ctx, nil)
		if err != nil {
			log.Fatalf("Failed to get latest block: %v", err)
		}
		startBlock = hexutil.EncodeBig(block.Number())
		log.Printf("Starting to walk back from latest block number: %s", startBlock)
	}

	startBlockNumber := hexutil.MustDecodeBig(startBlock).Uint64()
	endBlockNumber := startBlockNumber - (24 * 60 * 60 / 2 * 30) // 30 days
	for i := startBlockNumber; i >= endBlockNumber; i-- {
		proceed, err := checkPrerequisites(absolutePath, i, overwrite == "true")
		if err != nil || !proceed {
			log.Printf("Skipping block %d, %v, %w", i, proceed, err)
			continue
		}

		log.Printf("Generating execution witness for block %d / %s", i, hexutil.EncodeUint64(i))
		generateExecutionWitness(ctx, reth, geth, big.NewInt(int64(i)), absolutePath)
	}
}

func checkPrerequisites(parentPath string, blockNumber uint64, overwrite bool) (bool, error) {
	path := filepath.Join(parentPath, hexutil.EncodeUint64(blockNumber))
	info, err := os.Stat(path)

	// if folder exists
	if err == nil && info.IsDir() {
		// do we want to overwrite?
		if overwrite {
			if err = os.RemoveAll(path); err != nil {
				log.Printf("Failed to remove directory %s: %v", path, err)
				return false, err
			}
		} else {
			log.Printf("Directory %s already exists", path)
			return false, nil
		}
	}

	// create the directory
	if err = os.MkdirAll(path, 0755); err != nil {
		log.Printf("Failed to create directory %s: %v", path, err)
		return false, err
	}

	return true, nil
}

func dirExistsRelative(path string) (string, bool) {
	_, currentFile, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(currentFile)

	absolutePath := filepath.Join(currentDir, path)

	info, err := os.Stat(absolutePath)
	return absolutePath, err == nil && info.IsDir()
}

func generateExecutionWitness(ctx context.Context, reth *ethclient.Client, geth *ethclient.Client, blockNumber *big.Int, outputPath string) {
	folder := filepath.Join(outputPath, hexutil.EncodeBig(blockNumber))
	if err := callGenerateWitness(ctx, reth, blockNumber, folder, "reth"); err != nil {
		log.Fatalf("reth: Failed to generate execution witness for block %d: %w", blockNumber, err)
	}

	if err := callGenerateWitness(ctx, geth, blockNumber, folder, "geth"); err != nil {
		log.Fatalf("geth: Failed to generate execution witness for block %d: %w", blockNumber, err)
	}

	time.Sleep(500 * time.Millisecond)
}

func callGenerateWitness(ctx context.Context, client *ethclient.Client, blockNumber *big.Int, folder, filename string) error {
	start := time.Now()
	var witness ExecutionWitness
	input := hexutil.EncodeBig(blockNumber)
	if err := client.Client().CallContext(ctx, &witness, "debug_executionWitness", input); err != nil {
		return err
	}
	duration := time.Since(start)
	exectime, err := os.Create(filepath.Join(folder, fmt.Sprintf("%s_%s", filename, duration)))
	if err != nil {
		return err
	}
	defer exectime.Close()

	bytes, _ := json.Marshal(witness)
	file := filepath.Join(folder, fmt.Sprintf("%s.json", filename))
	if err := os.WriteFile(file, bytes, 0755); err != nil {
		return err
	}

	return nil
}

// ExecutionWitness is a witness json encoding for transferring across clients
// in the future, we'll probably consider using the extWitness format instead for less overhead.
// currently we're using this format for compatibility with the reth and also for simplicify in terms of parsing.
type ExecutionWitness struct {
	Headers []*types.Header   `json:"headers"`
	Codes   map[string]string `json:"codes"`
	State   map[string]string `json:"state"`
}
