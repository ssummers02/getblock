package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const api = ""

func main() {
	// Get the latest block number
	blockNumber := getLatestBlockNumber()

	// Get the 100 blocks before the latest one
	hundredBlocksBefore := subHex(blockNumber, "0x64")

	// Get the blocks between hundredBlocksBefore and blockNumber
	blocks := getBlocksInRange(hundredBlocksBefore, blockNumber)
	// Calculate the balance changes for each address in each block
	balanceChanges := calculateBalanceChanges(blocks)

	// Find the address with the largest absolute balance change
	maxAddress, maxChange := findMaxBalanceChange(balanceChanges)

	// Display the address and its balance change
	fmt.Printf("Address %s had a balance change of %s\n", maxAddress, maxChange)
}

// Gets the latest block number using the getBlockNumber API endpoint
func getLatestBlockNumber() string {
	reqBody := []byte(`{
		"jsonrpc": "2.0",
		"method": "eth_blockNumber",
		"params": [],
		"id": "getblock.io"
	}`)
	resp, err := responseGet(reqBody)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var blockNumberResponse struct {
		Result string `json:"result"`
	}

	err = json.Unmarshal(bodyBytes, &blockNumberResponse)
	if err != nil {
		panic(err)
	}

	return blockNumberResponse.Result
}

func responseGet(reqBody []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", "https://eth.getblock.io/mainnet/", bytes.NewReader(reqBody))
	if err != nil {
		panic(err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-api-key", api)

	// Create a new HTTP client with the custom headers
	client := &http.Client{
		Transport: &http.Transport{},
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	return resp, err
}

// Gets the blocks in the range between startBlock and endBlock using the getBlockByNumber API endpoint
func getBlocksInRange(startBlock string, endBlock string) []Block {
	blocks := make([]Block, 0, 100)

	for i := startBlock; i <= endBlock; i = addHex(i, "0x1") {
		reqBody := []byte(fmt.Sprintf(`{
			"jsonrpc": "2.0",
			"method": "eth_getBlockByNumber",
			"params": ["%s", true],
			"id": "getblock.io"
		}`, i))

		resp, err := responseGet(reqBody)
		if err != nil {
			return nil
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		var blockResponse struct {
			Result Block `json:"result"`
		}

		err = json.Unmarshal(bodyBytes, &blockResponse)
		if err != nil {
			panic(err)
		}

		blocks = append(blocks, blockResponse.Result)
	}

	return blocks
}
