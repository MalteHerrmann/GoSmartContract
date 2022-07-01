// QueryAndTransfer.go
//
// This script queries the MaltCoin token balance of an account
// and transfer some tokens to another account
//
// The script has to be called with the following arguments
//
// go run query_and_transfer.go $CONTRACT_ADDRESS $SENDER_PRIVKEY $RECIPIENT_ADDRESS $AMOUNT
package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	maltcoin "github.com/MalteHerrmann/GoSmartContract/contracts/build"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// Print a description
	fmt.Println("\nquery_and_transfer.go\n-----------------------------------------------------")
	fmt.Printf("This script loads a Maltcoin smart contract, that's deployed to a \nlocal Evmos node, queries token balances and transfers tokens between users.\n\n")

	// Process input
	contractAddress := common.HexToAddress(os.Args[1])
	senderPrivateKey := os.Args[2]
	recipientAddress := common.HexToAddress(os.Args[3])
	amount := os.Args[4]

	// // Initialize call data
	// var callData []byte

	// Convert private key to ECDSA format
	ecdsaPrivateKey, err := crypto.HexToECDSA(senderPrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	// Derive sender address from public key
	senderAddress := crypto.PubkeyToAddress(ecdsaPrivateKey.PublicKey)

	// Convert amount to big integer
	tempBig := big.NewInt(0)
	amountBig, ok := tempBig.SetString(amount, 10)
	if !ok {
		log.Fatalf("Failed to convert amount to big.Int: %v\n", err)
	}
	fmt.Println("\nAmount to transfer: ", amountBig)

	// Use ethclient to connect to local Evmos node on port 8545
	client, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		log.Fatalf("Failed to connect to local Evmos node: %v\n", err)
	}
	fmt.Println("Connected to local Evmos node on Port 8545.")

	// Get chain ID from client
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatalf("Failed to return chain id: %v\n", err)
	}
	fmt.Println("Chain ID:", chainID)

	// Suggest gas price from client
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Failed to get gas price: %v\n", err)
	}

	// Get pending nonce for sender account
	nonce, err := client.PendingNonceAt(context.Background(), senderAddress)
	if err != nil {
		log.Fatalf("Failed to get nonce: %v\n", err)
	}

	// // In order to estimate the gas cost of the transaction, we need to
	// // fill an ethereum.CallMsg struct.
	// //
	// // According to https://goethereumbook.org/en/transfer-tokens/, the
	// // data field has to be a byte array filled with:
	// //   - Method ID of the called method
	// //   - Padded recipient address
	// //   - Padded transferred token amount
	// //
	// // Generate method ID from the function signature of the called method.
	// transactionFnSignature := []byte("transfer(address,uint256)")
	// signatureHash := sha3.NewLegacyKeccak256()
	// signatureHash.Write(transactionFnSignature)
	// methodID := signatureHash.Sum(nil)[:4]
	// fmt.Println("Method ID:", hexutil.Encode(methodID))

	// // Generate byte array representation of padded recipient address
	// paddedAddress := common.LeftPadBytes(recipientAddress.Bytes(), 32)

	// // Generate byte array representation of padded recipient address
	// paddedAmount := common.LeftPadBytes(amountBig.Bytes(), 32)

	// // Fill data byte array
	// callData = append(callData, methodID...)
	// callData = append(callData, paddedAddress...)
	// callData = append(callData, paddedAmount...)

	// // Estimate gas usage to call the transfer function
	// gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
	// 	From: senderAddress,
	// 	Data: callData,
	// })
	// if err != nil {
	// 	log.Fatalf("Failed to estimate gas: %v\n", err)
	// }

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

	// Get token symbol
	symbol, err := contract.Symbol(nil)
	if err != nil {
		log.Fatalf("Failed to retrieve token name: %v\n", err)
	}
	fmt.Println("Token symbol: ", symbol)

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

	// Print balances
	fmt.Printf("\nBefore the transaction, the account balances are as follows (in a%v):\n", symbol)
	fmt.Printf("\n%v: %v", senderAddress, senderBalance)
	fmt.Printf("\n%v: %v", recipientAddress, recipientBalance)

	// Define transaction signer
	auth, err := bind.NewKeyedTransactorWithChainID(ecdsaPrivateKey, chainID)
	if err != nil {
		log.Fatalf("Failed to create transaction signer: %v\n", err)
	}
	auth.From = senderAddress
	auth.Nonce = big.NewInt(int64(nonce))
	// auth.GasLimit = gasLimit
	auth.GasPrice = gasPrice

	// Transfer tokens from deployer address to recipient address
	tx, err := contract.Transfer(auth, recipientAddress, amountBig)
	if err != nil {
		log.Fatalf("Failed to transfer tokens: %v\n", err)
	}
	fmt.Println("\nTokens transferred in tx with hash: ", tx.Hash().Hex())
}
