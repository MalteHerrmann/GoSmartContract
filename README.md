# GoSmartContract
In this repo, an ERC20 token smart contract is deployed to a local Evmos node.

## Pre-Requisites
The following software has to be installed on your machine in order to use the 
latest version of Evmos (currently v5.0.0):
- Go v1.18+ (https://go.dev/)
- Solidity compiler (https://docs.soliditylang.org/en/v0.8.15/)
- Evmos Daemon (https://docs.evmos.org/validators/quickstart/installation.html)

## Evmos Node
### Configuration
After a fresh installation of the `evmosd` CLI, the node has to be configured. 
This can either be done (manually)[https://docs.evmos.org/validators/quickstart/run_node.html#manual-deployment] 
or using the `init.sh` shell script, that is contained in the GitHub repository. 
Upon inspection of said script, one can see, that this pre-configures a local
node for testing purposes, using the `test` keyring-backend, creating an initial
account, with an initial supply of tokens, and more.

### Running the node
You can start your configured Evmos node using `evmosd start` and should see blocks 
being produced.\\
Now it's possible to interact with the node through the CLI. For example, one can
list the available accounts using `evmosd keys list`. Upon the first execution, 
you will only see the genesis account(s).
In order to add more accounts or interact with the node in another way, refer to the
(docs)[docs.evmos.org].

## ERC20 Smart Contract
A 
