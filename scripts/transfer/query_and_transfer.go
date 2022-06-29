// QueryAndTransfer.go
//
// *WIP*: This script currently only retrieves the token name.
//
// This script queries the MaltCoin token balance of an account
// and transfer some tokens to another account
package main

import (
	"fmt"
	"log"
	"os"

	maltcoin "github.com/MalteHerrmann/GoSmartContract/contracts/build"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// Print a description
	fmt.Println("\nquery_and_transfer.go\n-----------------------------------------------------")
	fmt.Printf("This script loads a Maltcoin smart contract, that's deployed to a \nlocal Evmos node, queries token balances and transfers tokens between users.\n\n")

	// Use ethclient to connect to local Evmos node on port 8545
	client, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		log.Fatalf("Failed to connect to local Evmos node: %v\n", err)
	}
	fmt.Println("Connected to local Evmos node on Port 8545.")

	// Define contract address
	// contractAddress := common.HexToAddress("0x0ED5a4E91490DAc67aAE538D0e77680141Fd6e5B")
	contractAddress := common.HexToAddress(os.Args[1])

	// Create an instance of the Maltcoin smart contract
	contract, err := maltcoin.NewMaltcoin(contractAddress, client)
	if err != nil {
		log.Fatalf("Failed to load token contract: %v\n", err)
	}
	fmt.Println("Maltcoin contract loaded at address: ", contractAddress)

	// Get name of token
	name, err := contract.Name(nil)
	if err != nil {
		log.Fatalf("Failed to retrieve token name: %v\n", err)
	}
	fmt.Println("Token name: ", name)

	// // Define sender address
	// senderAddress := common.HexToAddress("0x193bf98e7999646b74A139DBF2fB3e74d380767A")

	// // Query balance of Maltcoin tokens for deployer address
	// balance, err := contract.BalanceOf(nil, senderAddress)
	// if err != nil {
	// 	log.Fatalf("Failed to retrieve balance: %v\n", err)
	// }
	// fmt.Println("Balance of Maltcoin tokens for address ", senderAddress, " is: ", balance)
}
