# GoSmartContract
In this repo, an ERC20 token smart contract is deployed to a local Evmos node.

- [Pre-Requisites](#pre-requisites)
- [Short Summary](#short-summary)
- [Evmos Node](#evmos-node)
  - [Configuration](#configuration)
  - [Running the Node](#running-the-node)
- [ERC20 Smart Contract](#erc20-smart-contract)
  - [Compilation](#compilation)
- [Deployment Using Go-Ethereum](#deployment-using-go-ethereum)
- [Further Scope](#further-scope)

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
the OpenZeppelin library of smart contracts is used. <br>
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
of tokens. <br>
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
 $ go run github.com/MalteHerrmann/GoSmartContract/scripts/deploy $PRIVKEY
```

```
deploy_contract.go
-----------------------------------------------------
This script deploys a contract to a local Evmos node.

Current nonce:  81
Estimated gas: 1190381
Suggested gas price: 7

*********** Success ***********
The token contract was deployed in transaction  0xfcc62270b21c303ddfd39967ee956985906da4ee83af9b343a64c02696375e4a
The contract address is  0x089e91Aae4Bb044DD1477cCf43499e4E4758dEBD
```

Execute `receipt.go` from the `scripts` subfolder to print the 
contents of the transaction receipt. This is useful to check, if 
there is any valid contract code at the contract address. For example,
if too little gas is provided for the transaction, the code at the
address is `[]` and the receipt status is `0`. The transaction hash
has to be given as the first call argument.

```shell
 $ go run github.com/MalteHerrmann/GoSmartContract/scripts/receipt $TXHASH
```

````
receipt.go
-----------------------------------------------------
This script prints values from the transaction receipt, given a valid tx hash.


-------------
Transaction:
0xfcc62270b21c303ddfd39967ee956985906da4ee83af9b343a64c02696375e4a

Blocknumber:       102785
Contract address:  0x089e91Aae4Bb044DD1477cCf43499e4E4758dEBD
Status:            1
Gas used:          1190381
Logs:              [0x1400013a840]
Length of code at contract address:  4707
````

Another script is provided, which can be used to query the token name
and symbol, and account balances, as well as transfer Maltcoin tokens
between two accounts.
In order to execute these contract calls, the script has to be called 
with the `$CONTRACT` address of the ERC20 token contract, the signer's
private key `$PRIVKEY`, the `$RECIPIENT` address, and a token `$AMOUNT`,
which should be transferred.

```shell
 $ go run github.com/MalteHerrmann/GoSmartContract/scripts/query_and_transfer $CONTRACT $PRIVKEY $RECIPIENT $AMOUNT
```
```
query_and_transfer.go
-----------------------------------------------------
This script loads a Maltcoin smart contract, that's deployed to a
local Evmos node, queries token balances and transfers tokens between users.

Maltcoin contract loaded at address:  0xFdCa4BBB8040A59A7C2f1eF5b59BDa338791fe78
Token name:  Maltcoin
Token symbol:  MALT


Account balances pre transaction (in aMALT):
                  ADDRESS                    |               BALANCE
---------------------------------------------|----------------------------------
0x193bf98e7999646b74A139DBF2fB3e74d380767A   | 9999999999999999880000
0xcbAe3855CeDB30ce2Dd5766B82A12a1Ff6c32D25   | 120000


10000 tokens transferred in tx 0xa9f7d8cb3a5a84c8740cd106c5334bdb13d09d4b81087a681fbc3ad2860dc557


Account balances post transaction (in aMALT):
                  ADDRESS                    |               BALANCE
---------------------------------------------|----------------------------------
0x193bf98e7999646b74A139DBF2fB3e74d380767A   | 9999999999999999870000
0xcbAe3855CeDB30ce2Dd5766B82A12a1Ff6c32D25   | 130000

```

The three scripts all access utility functions, which are defined 
in `scripts/util/util.go`. 
This additional file was created, to have a central library of functions
and variables readily available in order to write further 
client scripts. 

## Testing

There are two commands for testing purposes:

- Unit testing for utility functions in Go:
    ```shell
    $ go test github.com/MalteHerrmann/GoSmartContract/scripts/util
    ```

- Testing the ERC20 token
    ```shell
    $ go test github.com/MalteHerrmann/GoSmartContract/tests
    ```

Please bear in mind, that the Solidity contract has 
to be compiled **before** the tests are run, because 
they depend on the generated ABI. 
Also, for some of the tests it is necessary 
to have a local Evmos node running and 
to adjust the value of the transaction hash (`testTxHashHex` in `util_test.go`) 
for testing purposes to a valid one.

Within the test files, there are two distinct approaches to testing to be mentioned:

- `util_test.go` contains [table-driven tests](https://dev.to/boncheff/table-driven-unit-tests-in-go-407b)
- `maltcoin_bdd_test.go` contains [BDD](https://www.bddtesting.com/what-is-bdd-testing/)-style tests

## Further scope

Additional things, that may be done for the further development of this basic repository:

- Customize the ERC20 token contract, which is just out of the box for now
- Currently, some ERC20 methods are untested, like `increaseAllowance` or `decreaseAllowance`, so tests for these can be added.
- Build an interactive command prompt for interactions with a Maltcoin token contract
- Use Go generics to reduce separate functions for simulated backend and actual client
- Use events to determine whether a transaction was included in a block instead of waiting some time.

Some remarks, that have occured to me during work on this task, are documented in the [Remarks](./docs/remarks.md) file.
