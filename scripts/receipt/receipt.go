// receipt.go prints information from the transaction receipt
// given a valid transaction hash in hex format.
//
// Usage:
//
//  $ go run receipt.go $TXHASH
//
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/MalteHerrmann/GoSmartContract/scripts/util"
	"github.com/ethereum/go-ethereum/common"
)

func main() {
	// Get script arguments
	txHashHex := os.Args[1]

	// Connect to local evmos node
	client, err := util.GetClient()
	if err != nil {
		log.Fatalf("Failed to connect to local Evmos node: %v\n", err)
	}

	// Get transaction receipt using the client and transaction hash,
	// which must be given as the first command line argument.
	receipt, err := util.GetReceipt(client, txHashHex)
	if err != nil {
		log.Fatalf("Failed to retrieve receipt: %v\n", err)
	}
	// Print information to terminal output
	fmt.Printf("\n-------------\nTransaction:\n%s\n\n", txHashHex)
	fmt.Println("Blocknumber:      ", receipt.BlockNumber)
	fmt.Println("Contract address: ", receipt.ContractAddress)
	fmt.Println("Status:           ", receipt.Status)
	fmt.Println("Gas used:         ", receipt.GasUsed)
	fmt.Println("Logs:             ", receipt.Logs)

	// Get the code stored at the contract address
	if (receipt.ContractAddress != common.Address{}) {
		code, err := client.CodeAt(context.Background(), receipt.ContractAddress, nil)
		if err != nil {
			log.Fatalf("Failed to retrieve code: %v\n", err)
		}
		fmt.Println("Length of code at contract address: ", len(code))
	}
}
