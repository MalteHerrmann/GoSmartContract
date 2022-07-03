// deployContract.go is a script to deploy an ERC20 token
// contract a local Evmos node.
// It uses the go implementation of a Solidity contract, that
// was generated using the Solidity compiler and abigen.
//
// It must be called with the private key in hex format, that
// which will be used to deploy the contract.
//
// Usage:
//
//  $ go run deploy_contract.go $PRIVKEY
//
package main

import (
	"fmt"
	"log"
	"os"

	maltcoin "github.com/MalteHerrmann/GoSmartContract/contracts/build"
	"github.com/MalteHerrmann/GoSmartContract/scripts/util"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	// Get ecdsa representation of private key, which is given as the first
	// command line argument.
	privKey, err := crypto.HexToECDSA(os.Args[1])
	if err != nil {
		log.Fatalf("Error while converting the private key to ecdsa: %v", err)
	}

	// Connect to local EVM and return the client plus a transaction signer,
	// that can be used to deploy the contract.
	client, auth, err := util.GetClientAndTransactionSigner(privKey)
	if err != nil {
		log.Fatalf("Error while connecting to the local node and getting the transaction signer: %v", err)
	}

	// Define data that should be executed on the contract (in this case deployment)
	callData := common.FromHex(maltcoin.MaltcoinMetaData.Bin)

	// Define the ethereum call message, which contains necessary information
	// to estimate gas consumption in order to fill all transaction signer
	// fields.
	callMsg := ethereum.CallMsg{
		From: auth.From,
		To:   nil,
		Data: callData,
	}

	// Fill transaction signer fields for this specific transaction
	auth, err = util.FillTransactionSignerFields(auth, client, callMsg)
	if err != nil {
		log.Fatalf("Error while filling transaction signer fields: %v", err)
	}

	// Deploy the contract
	contractAddress, tx, _, err := maltcoin.DeployMaltcoin(auth, client)
	if err != nil {
		log.Fatalf("Error while deploying the token contract: %v", err)
	}

	// Print information into terminal output
	fmt.Println("\ndeploy_contract.go\n-----------------------------------------------------")
	fmt.Printf("This script deploys a contract to a local Evmos node.\n\n")
	fmt.Println("Connected to local Evmos node on Port 8545.")
	fmt.Println("Current nonce: ", auth.Nonce)
	fmt.Println("Estimated gas:", auth.GasLimit)
	fmt.Println("Suggested gas price:", auth.GasPrice)
	fmt.Println("\n*********** Success ***********")
	fmt.Println("The token contract was deployed in transaction ", tx.Hash().Hex())
	fmt.Println("The contract address is ", contractAddress)
}
