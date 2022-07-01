# ------------------------
# Script to successfully deploy a Solidity 
# smart contract to a local Evmos node.
# 
# After deployment, the transaction receipt is read
# and a simple token transfer between two accounts
# is executed.
# ------------------------
# User settings
SENDER_KEYNAME=mykey
RECIPIENT_KEYNAME=testKey
AMOUNT=10000000

# File paths
DEPLOY=scripts/deploy/deploy_contract.go
RECEIPT=scripts/receipt/receipt.go
QUERY=scripts/transfer/query_and_transfer.go

# Derive account information from evmosd CLI 
SENDER_PRIVKEY=$(evmosd keys unsafe-export-eth-key $SENDER_KEYNAME --keyring-backend=test)
SENDER_BECH32=$(evmosd keys show $SENDER_KEYNAME| grep 'address' | grep -o 'evmos[0-9a-z]*')
SENDER_HEX=0x$(evmosd keys parse $SENDER_BECH32 | grep 'bytes' | grep -o '[0-9A-Z]*')
RECIPIENT_BECH32=$(evmosd keys show $RECIPIENT_KEYNAME| grep 'address' | grep -o 'evmos[0-9a-z]*')
RECIPIENT_HEX=0x$(evmosd keys parse $RECIPIENT_BECH32 | grep 'bytes' | grep -o '[0-9A-Z]*')

# Remove previous build files
rm -rf contracts/build*

# Compile contract
solc --abi contracts/Maltcoin.sol -o contracts/build
solc --bin contracts/Maltcoin.sol -o contracts/build

# Generate go bindings
abigen --abi=contracts/build/Maltcoin.abi --bin=contracts/build/Maltcoin.bin --pkg=maltcoin --out=contracts/build/Maltcoin.go

# Run deployment function
go run scripts/deploy/deploy_contract.go $SENDER_PRIVKEY > tmp.txt
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
go run $QUERY $CONTRACT $SENDER_PRIVKEY $RECIPIENT_HEX $AMOUNT

