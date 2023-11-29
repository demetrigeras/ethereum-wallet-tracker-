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
	"github.com/flosch/pongo2/v4"
	"github.com/gin-gonic/gin"
)

// Transaction struct used in your other file.
// type Transactions struct {
// 	Transactionsinfura string `json:"transaction-hash"`
// 	// Define your transaction fields here.
// }

func Startweb() {
	r := gin.Default()

	// Serve HTML templates
	r.LoadHTMLGlob("templates/*")

	// Route for the home page
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", pongo2.Context{})
	})

	// Route for submitting a wallet address

	r.POST("/fetch_transactions", func(c *gin.Context) {
		walletAddress := c.PostForm("walletAddress")
		transactions, err := FetchTransactions(walletAddress)
		if err != nil {
			// Handle the error
			c.HTML(http.StatusInternalServerError, "error.html", pongo2.Context{"message": err.Error()})
			return
		}

		c.HTML(http.StatusOK, "transactions.html", pongo2.Context{"transactions": transactions})

	})

	r.Run() // listen and serve on 0.0.0.0:8080
}

// FetchTransactions is a placeholder for your actual transaction fetching function.
func FetchTransactions(walletAddress string) ([]Transaction, error) {
	etherscanAPIKey := "J48JTZN2IZVYUUA7DME5IVPBKMM4YAMF68"
	infuraProjectID := "f1bfabaa66614342a34543701a76b373"

	// Fetch transactions from Etherscan
	etherscanURL := fmt.Sprintf("https://api.etherscan.io/api?module=account&action=txlist&address=%s&startblock=0&endblock=99999999&sort=asc&apikey=%s", walletAddress, etherscanAPIKey)
	resp, err := http.Get(etherscanURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var etherscanResponse EtherscanResponse
	if err := json.Unmarshal(body, &etherscanResponse); err != nil {
		return nil, err
	}

	var transactions []Transaction

	// Setup Infura client
	client, err := ethclient.Dial(fmt.Sprintf("https://mainnet.infura.io/v3/%s", infuraProjectID))
	if err != nil {
		return nil, err
	}

	// Fetch and add transaction details using Infura
	for _, tx := range etherscanResponse.Result {
		txHash := common.HexToHash(tx.Hash)
		transaction, isPending, err := client.TransactionByHash(context.Background(), txHash)
		if err != nil {
			log.Printf("Failed to fetch transaction %s: %v\n", tx.Hash, err)
			continue
		}
		if !isPending {
			// Create a Transaction struct and append it to the transactions slice
			transactions = append(transactions, Transaction{
				Hash:  tx.Hash,
				From:  tx.From,
				To:    tx.To,
				Value: transaction.Value().String(), // Example value
				// Add other details as needed
			})
		}
	}
	return transactions, nil
}