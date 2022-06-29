# ------------------------
# Script to reproduce the problem with 
# deploying the ERC20 token contract.
# 
# The deployment function itself throws no
# errors, but the transaction status is 0,
# when reading the transaction receipt.
#
# Also, there is no code at the contract 
# address. Calls to contract functions
# fail with 'no contract code at given address'.
# ------------------------
KEYNAME=mykey
PRIVKEY=$(evmosd keys unsafe-export-eth-key $KEYNAME --keyring-backend=test)
DEPLOY=scripts/deploy/deploy_contract.go
RECEIPT=scripts/receipt/receipt.go
QUERY=scripts/transfer/query_and_transfer.go

# Remove previous build files
rm -rf contracts/build*

# Compile contract
solc --abi contracts/Maltcoin.sol -o contracts/build
solc --bin contracts/Maltcoin.sol -o contracts/build

# Generate go bindings
abigen --abi=contracts/build/Maltcoin.abi --bin=contracts/build/Maltcoin.bin --pkg=maltcoin --out=contracts/build/Maltcoin.go

# Run deployment function
go run scripts/deploy/deploy_contract.go $PRIVKEY > tmp.txt
cat tmp.txt
TXHASH=$(cat tmp.txt | grep "transaction" | grep -o "0x[a-z0-9]*")
CONTRACT=$(cat tmp.txt | grep 'contract address' | grep -o '0x[0-9a-zA-Z]*')
rm -f tmp.txt

# Wait for transaction to be included in a block
echo "Waiting for transaction to be included in a block .. "
sleep 5

# Run receipt function
go run $RECEIPT $TXHASH

# Query name of token
go run $QUERY $CONTRACT
