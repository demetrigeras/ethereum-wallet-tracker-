package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// EtherscanResponse struct to hold the JSON response from Etherscan API
type EtherscanResponse struct {
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Result  []Transaction `json:"result"`
}

// Transaction struct to hold individual transaction details
type Transaction struct {
	Hash  string `json:"hash"`
	From  string `json:"from"`
	To    string `json:"to"`
	Value string `json:"value"`
	// Include other fields as required, like From, To, Value, etc.
}

func main() {
	// Replace with your wallet address and Etherscan API key
	walletAddress := "0xdE0336765d7549fB555883eB6c85e8862b4fDc41"
	apiKey := "J48JTZN2IZVYUUA7DME5IVPBKMM4YAMF68"

	// Constructing the Etherscan API URL
	url := fmt.Sprintf("https://api.etherscan.io/api?module=account&action=txlist&address=%s&startblock=0&endblock=99999999&sort=asc&apikey=%s", walletAddress, apiKey)

	// Making the HTTP GET request to the Etherscan API
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error fetching transactions: %v", err)
	}
	defer resp.Body.Close()

	// Reading the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	// Unmarshalling the JSON response into the EtherscanResponse struct
	var response EtherscanResponse
	if err := json.Unmarshal(body, &response); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	// Iterating over and printing each transaction
	for _, tx := range response.Result {
		fmt.Printf("Transaction Hash: %s\n", tx.Hash)
		fmt.Printf("Transaction From: %s\n", tx.From)
		fmt.Printf("Transaction To: %s\n", tx.To)
		fmt.Printf("Transaction Value: %s\n", tx.Value)
		// Print additional transaction details as needed
	}
}
