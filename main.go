package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EtherscanResponse struct {
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Result  []Transaction `json:"result"`
}

type Transaction struct {
	Hash  string `json:"hash"`
	From  string `json:"from"`
	To    string `json:"to"`
	Value string `json:"value"`
}

func main() {

	walletAddress := "0xdE0336765d7549fB555883eB6c85e8862b4fDc41"
	etherscanAPIKey := "J48JTZN2IZVYUUA7DME5IVPBKMM4YAMF68"
	infuraProjectID := "f1bfabaa66614342a34543701a76b373"

	// Fetch transactions from Etherscan
	etherscanURL := fmt.Sprintf("https://api.etherscan.io/api?module=account&action=txlist&address=%s&startblock=0&endblock=99999999&sort=asc&apikey=%s", walletAddress, etherscanAPIKey)
	resp, err := http.Get(etherscanURL)
	if err != nil {
		log.Fatalf("Error fetching transactions from Etherscan: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading Etherscan response: %v", err)
	}

	var etherscanResponse EtherscanResponse
	if err := json.Unmarshal(body, &etherscanResponse); err != nil {
		log.Fatalf("Error parsing Etherscan JSON: %v", err)
	}

	// Setup Infura client
	client, err := ethclient.Dial(fmt.Sprintf("https://mainnet.infura.io/v3/%s", infuraProjectID))
	if err != nil {
		log.Fatalf("Failed to connect to Infura: %v", err)
	}

	// Fetch and print transaction details using Infura
	for _, tx := range etherscanResponse.Result {
		txHash := common.HexToHash(tx.Hash)
		transaction, isPending, err := client.TransactionByHash(context.Background(), txHash)
		if err != nil {
			log.Printf("Failed to fetch transaction %s: %v\n", tx.Hash, err)
			continue
		}
		if !isPending {
			fmt.Printf("Transaction Hash: %s, Gas Price: %s\n", tx.Hash, transaction.GasPrice().String())
		}
	}
	Startweb()

}
