# GoSmartContract
In this repo, an ERC20 token smart contract is deployed to a local Evmos node.

## Pre-Requisites
The following software has to be installed on your machine in order to use the 
latest version of Evmos (currently v5.0.0):
- Go v1.18+ (https://go.dev/)
- Solidity compiler (https://docs.soliditylang.org/en/v0.8.15/installing-solidity.html#macos-packages)
- Evmos Daemon (https://docs.evmos.org/validators/quickstart/installation.html)

## Evmos Node
### Configuration
After a fresh installation of the `evmosd` CLI, the node has to be configured. 
This can either be done [manually](https://docs.evmos.org/validators/quickstart/run_node.html#manual-deployment) 
or using the `init.sh` shell script, that is contained in the GitHub repository. <br>
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

### Deployment using go-ethereum

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


