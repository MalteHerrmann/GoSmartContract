// maltcoin_test.go
//
// Testing suite for the Maltcoin smart contract
//
package maltcoin_test

import (
	"math/big"
	"testing"

	maltcoin "github.com/MalteHerrmann/GoSmartContract/contracts/build"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

func TestMaltcoin(t *testing.T) {
	// Define gas limit
	gasLimit := uint64(1000000)

	// Generate private keys
	privateKeySender, err := crypto.GenerateKey()
	require.NoError(t, err, "Could not generate sender's private key")

	// Define transaction signer
	auth, err := bind.NewKeyedTransactorWithChainID(privateKeySender, big.NewInt(1337))
	require.NoError(t, err, "Failed to create transaction signer")

	// Define balance for sender account
	balance := big.NewInt(1000000000000000000)

	// Define genesis state for simulated backend
	address := auth.From
	genesisAlloc := map[common.Address]core.GenesisAccount{
		address: {
			Balance: balance,
		},
	}

	// Use simulated backend for testing
	client := backends.NewSimulatedBackend(genesisAlloc, gasLimit)

	// Deploy contract
	// contractAddress, tx, contract, err := maltcoin.DeployMaltcoin(auth, client)
	_, _, contract, err := maltcoin.DeployMaltcoin(auth, client)
	require.NoError(t, err, "Could not deploy contract")

	// Commit transaction
	client.Commit()

	// Get name of token
	name, err := contract.Name(nil)
	require.NoError(t, err, "Could not retrieve token name")
	require.Equal(t, "Maltcoin", name, "Token name should be Maltcoin")
}
