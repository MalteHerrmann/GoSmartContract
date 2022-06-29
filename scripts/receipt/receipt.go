package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// Print a description
	fmt.Println("\n-----------------------------------------------------")
	fmt.Printf("This script prints information from the transaction receipt\ngiven a valid transaction hash.\n\n")

	// Use ethclient to connect to local Evmos node on port 8545
	client, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		log.Fatalf("Failed to connect to local Evmos node: %v\n", err)
	}
	fmt.Println("Connected to local Evmos node on Port 8545.")

	// Define transaction hash, for which the receipt should be returned
	txHash := common.HexToHash("0xab693c5de3f831adcfe81032004b1f07ee5b13ac0e847f18010a7b0907750569")

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

	// Define contract address
	contractAddress := common.HexToAddress("0xFc3e94E429Bf36099d235b421ec770eB9AFb3b7F")

	// Get the code stored at the contract address
	code, err := client.CodeAt(context.Background(), contractAddress, nil)
	if err != nil {
		log.Fatalf("Failed to retrieve code: %v\n", err)
	}
	fmt.Println("Code at address: ", code)
}
