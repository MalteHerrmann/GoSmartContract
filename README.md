# GoSmartContract
In this repo, an ERC20 token smart contract is deployed to a local Evmos node.

## Pre-Requisites
The following software has to be installed on your machine in order to use the 
latest version of Evmos (currently v5.0.0):
- Go v1.18+ (https://go.dev/)
- Node JS (https://nodejs.org/en/)
- Solidity compiler (https://docs.soliditylang.org/en/v0.8.15/installing-solidity.html#macos-packages)
- Evmos Daemon (https://docs.evmos.org/validators/quickstart/installation.html)

## Short Summary
When you have installed the required software, configured and ran a local Evmos
node, you can clone this repository, and install the OpenZeppelin contracts using 
`npm install`. 

When executing 
```shell
 $ ./init.sh
```
an instance of the Maltcoin ERC20 token contract is deployed
to the running localnet, the contents of transaction receipt printed, as well as 
a simple transfer of Maltcoin tokens between two accounts executed.

## Evmos Node
### Configuration
After a fresh installation of the `evmosd` CLI, the node has to be configured. 
This can either be done [manually](https://docs.evmos.org/validators/quickstart/run_node.html#manual-deployment) 
or using the `init.sh` shell script, that is contained in the 
[Evmos GitHub repository](https://github.com/evmos/evmos). <br>
Upon inspection of said script, one can see, that this pre-configures a local
node for testing purposes, using the `test` keyring-backend, creating an initial
account, with an initial supply of tokens, and more.

### Running the node
You can start your configured Evmos node using `evmosd start` and should see blocks 
being produced. <br>
Now it's possible to interact with the node through the CLI. For example, one can
list the available accounts using `evmosd keys list`. Upon first execution, you will
only see the genesis account(s). <br>
In order to add more accounts or interact with the node in another way, refer to the
[docs](docs.evmos.org).

## ERC20 Smart Contract
For this exercise, a basic ERC20 token contract is deployed to the Evmos node. 
[ERC20](https://ethereum.org/en/developers/docs/standards/tokens/erc-20/) is a widely 
used standard to design fungible tokens, that offer a specific set of methods and
events.
In order to create a custom token to deploy, the basic `ERC20` contract from 
the OpenZeppelin library of smart contracts. <br>
These can be installed with NPM using `npm install @openzeppelin/contracts`.

### Compilation
In order to deploy the smart contract using go, it first must be compiled using
the Solidity compiler. We create the `.abi` as well as `.bin` files, which are
necessary to deploy and interact with the smart contracts.

```shell
 $ solc --abi contracts/Maltcoin.sol -o contracts/build
 $ solc --bin contracts/Maltcoin.sol -o contracts/build
```

These commands create the mentioned files in the subfolder `build`. Next, the
Go implementation of the contract is generated using `abigen`, which comes with
the installation of `solidity`.

```shell
 $ abigen --bin=contracts/build/Maltcoin.bin --abi=contracts/build/Maltcoin.abi --pkg=maltcoin --out=contracts/build/Maltcoin.go
```

The output of this contains the function `DeployMaltcoin`, which deploys the smart contract to the 
blockchain.

## Deployment using go-ethereum

To deploy the token contract, an account is needed. During the 
initialization of the Evmos node with the `./init.sh` script, 
an initial account was created and supplied with a specific amount 
of tokens. 
Manually, an account can be added using the `keys` command.

```shell
 $  evmosd keys add $KEYNAME
```

The available accounts can be queried with the following instructions:

```shell
 $ evmosd keys list
```

This will print the account list to the terminal output and display 
information like the account name, address and public key.

In order to be able to sign the transaction, which deploys the smart
contract, the private key is needed. It can be exported for a given `$KEYNAME` 
using: 

```shell
 $ evmosd keys unsafe-export-eth-key $KEYNAME --keyring-backend test
693F03A42E6F377D2305CB036EAE9BACCC09B230041CC786252A3BD5C34ED0FA
```

This private key `$PRIVKEY` can then be used to call the deployment 
script, which is part of this repository. It uses the `go-ethereum`
package in combination with the Go bindings, that were generated with
`abigen` to deploy an instance of the Maltcoin token contract on the
local Evmos node.

```shell
 $ go run github.com/MalteHerrmann/scripts/deploy $PRIVKEY

deploy_contract.go
-----------------------------------------------------
This script deploys a contract to a local Evmos node.

Connected to local Evmos node on Port 8545.
Current nonce:  62
Estimated gas: 1190369
Suggested gas price: 7

*********** Success ***********
The token contract was deployed in transaction  0xa5022fff3e376700a3a05a1f48d77b25718e001e2917f44c11c7b81b4e50d2ec
The contract address is  0x7AE756f54C887c1384e5f085d1B060d109E1B80e
```

Execute `receipt.go` from the `scripts` subfolder to print the 
contents of the transaction receipt. This is useful to check, if 
there is any valid contract code at the contract address. For example,
if too little gas is provided for the transaction, the code at the
address is `[]` and the receipt status is `0`. The transaction hash
has to be given as the first call argument.

```shell
 $ go run github.com/MalteHerrmann/scripts/receipt $TXHASH
Connected to local Evmos node at http://localhost:8545.

-------------
Transaction:
0xa5022fff3e376700a3a05a1f48d77b25718e001e2917f44c11c7b81b4e50d2ec

Blocknumber:       83917
Contract address:  0x7AE756f54C887c1384e5f085d1B060d109E1B80e
Status:            1
Gas used:          1190369
Logs:              [0x14000138840]
Length of code at contract address:  4707
```

Another script is provided, which can be used to query the token name
and symbol, and account balances, as well as transfer Maltcoin tokens
between two accounts.
In order to execute these contract calls, the script has to be called 
with the `$CONTRACT` address of the ERC20 token contract, the signer's
private key `$PRIVKEY`, the `$RECIPIENT` address, and a token `$AMOUNT`,
which should be transferred.

```shell
 $ go run github.com/MalteHerrmann/scripts/query_and_transfer $CONTRACT $PRIVKEY $RECIPIENT $AMOUNT

query_and_transfer.go
-----------------------------------------------------
This script loads a Maltcoin smart contract, that's deployed to a
local Evmos node, queries token balances and transfers tokens between users.


Amount to transfer:  10000000
Connected to local Evmos node on Port 8545.
Chain ID: 9000
Maltcoin contract loaded at address:  0x4122088F07d5f505caA67c125Bd89E793FB22274
Token name:  Maltcoin
Token symbol:  MALT

Before the transaction, the account balances are as follows (in aMALT):

0x193bf98e7999646b74A139DBF2fB3e74d380767A: 10000000000000000000000
0xcbAe3855CeDB30ce2Dd5766B82A12a1Ff6c32D25: 0
Tokens transferred in tx with hash:  0x06e38a22f9c1cbf76e51b449cdb33577ebb3bec5c10ee2f0397b696247d89e56
```

## Testing

There are two files for testing purposes:

- Unit testing for utility functions in Go:
    ```shell
    $ go test github.com/MalteHerrmann/scripts/util
    ```

- Testing the ERC20 token
    ```shell
    $ go test github.com/MalteHerrmann/tests
    ```
