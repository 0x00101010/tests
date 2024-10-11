package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// RPCResponse represents a JSON-RPC response
type RPCResponse struct {
	JsonRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	// Error   *RPCError   `json:"error,omitempty"`
	ID int `json:"id"`
}

type ExecutionWitness struct {
	Codes map[string]string `json:"codes"`
	State map[string]string `json:"state"`
}

func LoadExecutionWitnessFromFile(file string) (*ExecutionWitness, error) {
	_, err := os.Stat(file)
	if err != nil {
		return nil, err
	}

	bytes, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var resp RPCResponse
	if err := json.Unmarshal(bytes, &resp); err != nil {
		return nil, err
	}

	// fmt.Println(resp.Result)
	var ew ExecutionWitness
	if err := json.Unmarshal(resp.Result, &ew); err != nil {
		return nil, err
	}

	return &ew, nil
}

func main() {
	blockHash := os.Args[1]

	geth_witness_file := fmt.Sprintf("results/witness/base-sepolia/%s/geth.json", blockHash)
	reth_witness_file := fmt.Sprintf("results/witness/base-sepolia/%s/reth.json", blockHash)

	fmt.Printf("Witness files: %s, %s\n", geth_witness_file, reth_witness_file)

	geth_witness, err := LoadExecutionWitnessFromFile(geth_witness_file)
	if err != nil {
		panic(err)
	}

	reth_witness, err := LoadExecutionWitnessFromFile(reth_witness_file)
	if err != nil {
		panic(err)
	}

	fmt.Println("Geth witness:")
	fmt.Println(len(geth_witness.Codes))
	fmt.Println(len(geth_witness.State))

	fmt.Println("Reth witness:")
	fmt.Println(len(reth_witness.Codes))
	fmt.Println(len(reth_witness.State))

	count := 0
	for k, v := range geth_witness.State {
		if reth_witness.State[k] != v {
			fmt.Printf("State mismatch %d: %s\n", count, k)
			count++
		}
	}
	fmt.Println(count)

	count = 0
	for k, v := range geth_witness.Codes {
		// if !strings.Contains(reth_witness.Codes[k], v) {
		// 	fmt.Printf("Code mismatch %d: %s\n", count, k)
		// 	count++
		// }
		if reth_witness.Codes[k] != v {
			fmt.Printf("Code mismatch %d: %s\n", count, k)
			count++
		}
	}
	fmt.Println(count)

	// fmt.Println(geth_witness)
	// fmt.Println(reth_witness)

	// rethSet := convertMapToValueSet(reth_witness.Codes)
	// gethSet := convertMapToValueSet(geth_witness.Codes)
	// inSet1NotInSet2, inSet2NotInSet1 := findDifferences(rethSet, gethSet)
	// fmt.Println("In reth not in geth:", inSet1NotInSet2)
	// fmt.Println("In geth not in reth:", inSet2NotInSet1)
}

func convertMapToValueSet(m map[string]string) map[string]struct{} {
	set := make(map[string]struct{})
	for _, v := range m {
		set[v] = struct{}{}
	}
	return set
}

// findDifferences finds values in set1 that are not in set2 and vice versa
func findDifferences(set1, set2 map[string]struct{}) ([]string, []string) {
	inSet1NotInSet2 := []string{}
	inSet2NotInSet1 := []string{}

	for v := range set1 {
		if _, exists := set2[v]; !exists {
			inSet1NotInSet2 = append(inSet1NotInSet2, v)
		}
	}

	for v := range set2 {
		if _, exists := set1[v]; !exists {
			inSet2NotInSet1 = append(inSet2NotInSet1, v)
		}
	}

	return inSet1NotInSet2, inSet2NotInSet1
}
