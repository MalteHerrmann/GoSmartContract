// deployContract.go is a script to deploy an ERC20 token
// contract a local Evmos node.
// It uses the go implementation of a Solidity contract, that
// was generated using the Solidity compiler and abigen.
package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	maltcoin "github.com/MalteHerrmann/GoSmartContract/contracts/build"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// Print a description
	fmt.Println("\ndeploy_contract.go\n-----------------------------------------------------")
	fmt.Printf("This script deploys a contract to a local Evmos node.\n\n")

	// Use ethclient to connect to local Evmos node on port 8545
	client, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to local Evmos node on Port 8545.")

	// Get chain id from client
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Chain ID:", chainID)

	// Get gas price suggestion from client
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Suggested gas price:", gasPrice)

	// Define private key, which is needed to sign transactions
	privateKey := os.Args[1]

	// Get ecdsa representation of private key
	ecdsaPrivateKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// Get address from ecdsa key
	deployerAddress := crypto.PubkeyToAddress(ecdsaPrivateKey.PublicKey)
	fmt.Println("Deployer address:", deployerAddress)

	// Get current nonce for deployer address
	nonce, err := client.PendingNonceAt(context.Background(), deployerAddress)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Current nonce: ", nonce)

	// Define ethereum call message for contract deployment
	callMsg := ethereum.CallMsg{
		From:     deployerAddress,
		To:       nil,
		GasPrice: gasPrice,
		Data:     common.FromHex(maltcoin.MaltcoinMetaData.Bin),
	}

	// Estimate gas usage
	gasLimit, err := client.EstimateGas(context.Background(), callMsg)
	if err != nil {
		// log.Fatalf("Failed to estimate gas: %v\n", err)
		fmt.Printf("Failed to estimate gas: %v\n", err)
	} else {
		fmt.Println("Estimated gas:", gasLimit)
	}

	// Define transaction signer from private key and chain id and configure
	// the transaction options
	auth, err := bind.NewKeyedTransactorWithChainID(ecdsaPrivateKey, chainID)
	if err != nil {
		log.Fatal(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	// auth.GasLimit = 18446744073709551615 // max value for uint64
	// auth.GasLimit = 10000000000
	// auth.GasLimit = 11903790 // value worked
	auth.GasLimit = gasLimit
	auth.GasPrice = gasPrice
	auth.Value = big.NewInt(0)

	// Deploy the contract
	contractAddress, tx, _, err := maltcoin.DeployMaltcoin(auth, client)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n*********** Success ***********")
	fmt.Println("The token contract was deployed in transaction ", tx.Hash().Hex())
	fmt.Println("The contract address is ", contractAddress)
}
