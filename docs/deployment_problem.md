# Problem Summary
## TLDR
The deployment of the ERC20 token contract cannot be completed, because 
the gas limit for the transaction is reached, no matter what the value is set 
to. 

1. Clone repo <br>`$ gh repo clone MalteHerrmann/GoSmartContract`
2. Start local Evmos node <br> `$ evmosd start`
3. Install OpenZeppelin contracts <br> `$ npm install`
4. Reproduce ERC20 problem <br>`$ ./reproduce_problem.sh` 
5. Run test suite <br> `$ go test tests/maltcoin_test.go` -> fails with `contract creation code storage out of gas` error
6. Successfully deploy simpler contract <br> `$ ./init.sh` 
7. Run test suite <br> `$ go test tests/maltcoin_test.go` -> passes

## Description
I have stripped the ERC20 contract (`Maltcoin.sol`) down to the bare minimum implementation 
of an OpenZeppelin contract, that can be generated with https://wizard.openzeppelin.com/, 
and tested it with [Remix](https://remix.ethereum.org/). The compilation and deployment 
works fine in Remix and fails using `go-ethereum`. However, no error is thrown by the 
deployment function. The problem shows itself, when trying to call the contract functions, 
e.g. `maltcoin.Name(nil)`. The transaction status in the transaction receipt is `0` and the 
code at the `receipt.ContractAddress` is empty.

```
Status:                             0
Gas used:                           1000000              // == gasLimit
Length of code at contract address: 0
```

On the other hand, a simple Solidity contract (`Maltcoin_temp.sol`) can be deployed with the 
same Go functions, `solc` and `abigen` commands. Therefore, the problem is most likely in 
importing the ERC20 base contract from OpenZeppelin.

```
Status:                             1
Gas used:                           500000
Length of code at contract address: 557
```

There are encounters of this error in some ressources online, but some were either solved 
with a higher gas limit, which didn't help in this case, or were solved after commiting 
the transaction, which is not the problem either:
- https://github.com/ethereum/go-ethereum/issues/20636
- https://github.com/ethereum/go-ethereum/issues/15930
- https://ethereum.stackexchange.com/questions/68706/no-contract-code-at-give-address
- https://ethereum.stackexchange.com/questions/19725/contract-creation-code-storage-out-of-gas

## Testing
When testing with a simulated backend (`tests/maltcoin_test.go`), an error trace shows this too: 
```shell
 $ ./reproduce_problem.sh
...

 $ go test tests/maltcoin_test.go 
--- FAIL: TestMaltcoin (0.00s)
    maltcoin_test.go:49: 
                Error Trace:    /Users/malte/dev/go/GoSmartContract/tests/maltcoin_test.go:49
                Error:          Received unexpected error:
                                contract creation code storage out of gas
                Test:           TestMaltcoin
                Messages:       Could not deploy contract
FAIL
```

However, when compiling the simpler `Maltcoin_temp.sol` by executing `./init.sh`, and running 
the test suite, all tests pass:
```shell
 $ ./init.sh
...

 $ go test tests/maltcoin_test.go
ok      command-line-arguments  0.185s
```

This shows, that the test suite itself is set up correctly.

## Solution

The solution was to estimate the necessary gas limit using `client.EstimateGas(...)` instead of 
manually setting a value. Interestingly, even when manually setting a value considerably higher than the estimated gas,
the transaction does not deploy the code correctly, stating it runs out of gas.

The following transaction receipts can be generated with the estimated gas (`1190381`): 

```

-------------
Transaction:
0xfcc62270b21c303ddfd39967ee956985906da4ee83af9b343a64c02696375e4a

Blocknumber:       102785
Contract address:  0x089e91Aae4Bb044DD1477cCf43499e4E4758dEBD
Status:            1
Gas used:          1190381
Logs:              [0x1400013a840]
Length of code at contract address:  4707

```

As one can see, the status is `1`, the contract code was correctly deployed, and the gas consumption equals the estimated value.

On the other hand, when manually setting the gas limit to `10000000000` the contract deployment fails, while still all of the gas is consumend.

```

-------------
Transaction:
0x8a9aa6f668828bc0c04d7951db581bf4f80c3c2ab2dbd79fd08035d21fe44246

Blocknumber:       104960
Contract address:  0xd268fddeF1dF461A25D4E836e1094D646C0B705e
Status:            0
Gas used:          10000000000
Logs:              []
Length of code at contract address:  0

```
