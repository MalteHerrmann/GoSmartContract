// receipt.go
//
// This script prints information from the transaction receipt.
// Script must be called with a transaction argument as input.
//
// Example:
// go run receipt.go $TXHASH
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// Print a description
	fmt.Println("\nreceipt.go\n-----------------------------------------------------")
	fmt.Printf("This script prints information from the transaction receipt\ngiven a valid transaction hash.\n\n")

	// Use ethclient to connect to local Evmos node on port 8545
	client, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		log.Fatalf("Failed to connect to local Evmos node: %v\n", err)
	}
	fmt.Println("Connected to local Evmos node on Port 8545.")

	// Define transaction hash, for which the receipt should be returned
	// txHashHex := "0x900c0aa59e57327bcf26221b77c6904466c54f0b08918d99a3c838c233d13126"
	txHashHex := os.Args[1]
	txHash := common.HexToHash(txHashHex)

	// Get transaction receipt using the client and transaction hash
	receipt, err := client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		log.Fatalf("Failed to retrieve receipt: %v\n", err)
	}
	fmt.Printf("\n-------------\nTransaction:\n%s\n\n", txHash.Hex())
	fmt.Println("Blocknumber:      ", receipt.BlockNumber)
	fmt.Println("Contract address: ", receipt.ContractAddress)
	fmt.Println("Status:           ", receipt.Status)
	fmt.Println("Gas used:         ", receipt.GasUsed)
	fmt.Println("Logs:             ", receipt.Logs)

	// Get the code stored at the contract address
	code, err := client.CodeAt(context.Background(), receipt.ContractAddress, nil)
	if err != nil {
		log.Fatalf("Failed to retrieve code: %v\n", err)
	}
	fmt.Println("Length of code at contract address: ", len(code))
}
