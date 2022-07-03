// query_and_transfer.go
//
// This script queries the MaltCoin token balance of an account
// and transfer some tokens to another account
//
// Usage:
//
//  $ go run query_and_transfer.go $CONTRACT_ADDRESS $SENDER_PRIVKEY $RECIPIENT_ADDRESS $AMOUNT
//
package main

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	maltcoin "github.com/MalteHerrmann/GoSmartContract/contracts/build"
	"github.com/MalteHerrmann/GoSmartContract/scripts/util"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	// Process input
	contractAddress := common.HexToAddress(os.Args[1])
	senderPrivateKey := os.Args[2]
	recipientAddress := common.HexToAddress(os.Args[3])
	amount := os.Args[4]

	// Convert private key to ECDSA format
	ecdsaPrivateKey, err := crypto.HexToECDSA(senderPrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	// Derive sender address from public key
	senderAddress := crypto.PubkeyToAddress(ecdsaPrivateKey.PublicKey)

	// Convert amount to big integer
	amountBig, ok := new(big.Int).SetString(amount, 10)
	if !ok {
		log.Fatalf("Failed to convert amount to big.Int: %v\n", err)
	}

	// Connect to local EVM and return the client and transaction signer
	client, auth, err := util.GetClientAndTransactionSigner(ecdsaPrivateKey)
	if err != nil {
		log.Fatalf("Error while connecting to the local node and getting the transaction signer: %v", err)
	}

	// Get the necessary call data byte array, that contains the
	// method name and its arguments.
	callData, err := util.GetCallData("transfer", recipientAddress, amountBig)
	if err != nil {
		log.Fatalf("Error while getting the call data: %v", err)
	}

	// Define the ethereum call message, which contains necessary information
	// to estimate gas consumption in order to fill all transaction signer
	// fields.
	callMsg := ethereum.CallMsg{
		From: senderAddress,
		To:   &contractAddress,
		Data: callData,
	}

	// Using the data in the call message struct, the transaction signer
	// can be configured for the transaction.
	auth, err = util.FillTransactionSignerFields(auth, client, callMsg)
	if err != nil {
		log.Fatalf("Error while filling the transaction signer fields: %v", err)
	}

	// Create an instance of the Maltcoin smart contract
	contract, err := maltcoin.NewMaltcoin(contractAddress, client)
	if err != nil {
		log.Fatalf("Failed to load token contract: %v\n", err)
	}

	// Get name of token
	name, err := contract.Name(nil)
	if err != nil {
		log.Fatalf("Failed to retrieve token name: %v\n", err)
	}

	// Get token symbol
	symbol, err := contract.Symbol(nil)
	if err != nil {
		log.Fatalf("Failed to retrieve token name: %v\n", err)
	}

	// Query balance of Maltcoin tokens for deployer address
	senderBalance, err := contract.BalanceOf(nil, senderAddress)
	if err != nil {
		log.Fatalf("Failed to retrieve balance: %v\n", err)
	}

	// Query balance of Maltcoin tokens for recipient address
	recipientBalance, err := contract.BalanceOf(nil, recipientAddress)
	if err != nil {
		log.Fatalf("Failed to retrieve balance: %v\n", err)
	}

	// Transfer tokens from deployer address to recipient address
	tx, err := contract.Transfer(auth, recipientAddress, amountBig)
	if err != nil {
		log.Fatalf("Failed to transfer tokens: %v\n", err)
	}

	// Wait some time for transaction to be included in a block
	time.Sleep(5 * time.Second)

	// Query balance of Maltcoin tokens for deployer address
	senderBalancePost, err := contract.BalanceOf(nil, senderAddress)
	if err != nil {
		log.Fatalf("Failed to retrieve balance: %v\n", err)
	}

	// Query balance of Maltcoin tokens for recipient address
	recipientBalancePost, err := contract.BalanceOf(nil, recipientAddress)
	if err != nil {
		log.Fatalf("Failed to retrieve balance: %v\n", err)
	}

	// Print output to terminal
	fmt.Println("\nquery_and_transfer.go\n-----------------------------------------------------")
	fmt.Printf("This script loads a Maltcoin smart contract, that's deployed to a \nlocal Evmos node, queries token balances and transfers tokens between users.\n\n")
	fmt.Println("Maltcoin contract loaded at address: ", contractAddress)
	fmt.Println("Token name: ", name)
	fmt.Println("Token symbol: ", symbol)
	fmt.Printf("\n\nAccount balances pre transaction (in a%v):\n", symbol)
	fmt.Printf("                  ADDRESS                    |               BALANCE           \n")
	fmt.Printf("---------------------------------------------|----------------------------------\n")
	fmt.Printf("%v   | %v\n", senderAddress, senderBalance)
	fmt.Printf("%v   | %v\n", recipientAddress, recipientBalance)
	fmt.Printf("\n\n%v tokens transferred in tx %v\n\n", amount, tx.Hash().Hex())
	fmt.Printf("\nAccount balances post transaction (in a%v):\n", symbol)
	fmt.Printf("                  ADDRESS                    |               BALANCE           \n")
	fmt.Printf("---------------------------------------------|----------------------------------\n")
	fmt.Printf("%v   | %v\n", senderAddress, senderBalancePost)
	fmt.Printf("%v   | %v\n\n", recipientAddress, recipientBalancePost)
}
