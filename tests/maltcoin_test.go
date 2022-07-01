// maltcoin_test.go
//
// Testing suite for the Maltcoin smart contract
//
package maltcoin_test

import (
	"crypto/ecdsa"
	"log"
	"math/big"
	"testing"

	maltcoin "github.com/MalteHerrmann/GoSmartContract/contracts/build"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

func TestTransfer(t *testing.T) {
	// TODO: blockGasLimit in suite hinterlegen
	// Test parameters
	blockGasLimit := uint64(4712388)
	chainID := big.NewInt(1337)

	// Generate private keys
	var privKeys []*ecdsa.PrivateKey
	var addresses []common.Address
	for i := 0; i < 2; i++ {
		privKey, err := crypto.GenerateKey()
		if err != nil {
			log.Fatalf("Error generating private key: %v\n", err)
		}
		privKeys = append(privKeys, privKey)
		addresses = append(addresses, crypto.PubkeyToAddress(privKey.PublicKey))
	}

	// Get simulated backend and transaction signer for testing
	client, auth, err := GetClientAndTransactionSigner(privKeys[0], blockGasLimit, chainID)
	require.NoError(t, err, "Error getting client and transaction signer")

	// Deploy contract
	_, _, contract, err := DeployContract(auth, client)
	require.NoError(t, err, "Could not deploy contract")

	// Get maltcoin token balance of account1
	//
	// Since account1 deployed the contract, it is assigned the deployer token balance.
	// This is 10000 * 10^18 in this case.
	big_10_18 := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	initialDeployerBalance := new(big.Int).Mul(big.NewInt(10000), big_10_18)

	// Check maltcoin token balance of account1
	// This account deployed the contract, so it should have the initial deployer token balance.
	balance, err := contract.BalanceOf(nil, addresses[0])
	require.NoError(t, err, "Could not retrieve balance of sender account before transfer")
	require.Equal(t, initialDeployerBalance, balance, "Wrong balance of sender account before transfer")

	// Get balance of account2
	balance, err = contract.BalanceOf(nil, addresses[1])
	require.NoError(t, err, "Could not retrieve balance of account2 before transfer")
	require.Equal(t, int64(0), balance.Int64(), "Wrong balance of account2 before transfer")

	// Transfer tokens from account1 to account2
	amount := new(big.Int).Mul(big.NewInt(1), big_10_18)
	_, err = contract.Transfer(auth, addresses[1], amount)
	require.NoError(t, err, "Could not transfer tokens")
	client.Commit()

	// Check if transfer was successful
	balance, err = contract.BalanceOf(nil, addresses[0])
	require.NoError(t, err, "Could not retrieve balance of sender account after transfer")
	require.Equal(t, new(big.Int).Sub(initialDeployerBalance, amount), balance, "Wrong balance of sender account after transfer")

	balance, err = contract.BalanceOf(nil, addresses[1])
	require.NoError(t, err, "Could not retrieve balance of recipient account after transfer")
	require.Equal(t, amount, balance, "Wrong balance of recipient account after transfer")
}

// This function tests if the token settings are set up correctly.
// It checks the name, symbol, decimals and total supply.
func TestTokenSettings(t *testing.T) {
	// Test parameters
	blockGasLimit := uint64(4712388)
	chainID := big.NewInt(1337)

	// Generate private key
	privKey, err := crypto.GenerateKey()
	require.NoError(t, err, "Error generating private key")

	// Get simulated backend and transaction signer for testing
	client, auth, err := GetClientAndTransactionSigner(privKey, blockGasLimit, chainID)
	require.NoError(t, err, "Error getting client and transaction signer")

	// Deploy contract
	_, _, contract, err := DeployContract(auth, client)
	require.NoError(t, err, "Could not deploy contract")

	// Get name of token
	name, err := contract.Name(nil)
	require.NoError(t, err, "Could not retrieve token name")
	require.Equal(t, "Maltcoin", name, "Token name should be Maltcoin")

	// Get token symbol
	symbol, err := contract.Symbol(nil)
	require.NoError(t, err, "Could not retrieve token symbol")
	require.Equal(t, "MALT", symbol, "Token symbol should be MALT")

	// Get token decimals
	decimals, err := contract.Decimals(nil)
	require.NoError(t, err, "Could not retrieve token decimals")
	require.Equal(t, uint8(18), decimals, "Token decimals should be 18")

	// Get total supply
	totalSupply, err := contract.TotalSupply(nil)
	require.NoError(t, err, "Could not retrieve total supply")
	require.Equal(t, "10000000000000000000000", totalSupply.String(), "Wrong total supply")
}

func GetClientAndTransactionSigner(privKey *ecdsa.PrivateKey, blockGasLimit uint64, chainID *big.Int) (*backends.SimulatedBackend, *bind.TransactOpts, error) {
	// Define initial aevmos balance for deployer account
	balance1 := big.NewInt(1000000000000000000)

	// Define genesis state for simulated backend
	address1 := crypto.PubkeyToAddress(privKey.PublicKey)
	genesisAlloc := map[common.Address]core.GenesisAccount{
		address1: {
			Balance: balance1,
		},
	}

	// Get simulated backend as client
	// client := backends.NewSimulatedBackend(genesisAlloc, uint64(4712388))
	client := backends.NewSimulatedBackend(genesisAlloc, blockGasLimit)

	// Define transaction signer
	//
	// As described in the go-ethereum documentation, the chain ID for
	// the simulated backend must be 1337.
	auth, err := bind.NewKeyedTransactorWithChainID(privKey, chainID)
	if err != nil {
		return nil, nil, err
	}

	return client, auth, nil
}

func DeployContract(auth *bind.TransactOpts, client *backends.SimulatedBackend) (common.Address, *types.Transaction, *maltcoin.Maltcoin, error) {
	// Deploy contract
	contractAddress, tx, contract, err := maltcoin.DeployMaltcoin(auth, client)
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	// Commit transaction
	client.Commit()

	return contractAddress, tx, contract, nil
}
