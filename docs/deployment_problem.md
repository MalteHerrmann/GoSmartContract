# Problem Summary
## TLDR
The deployment of the ERC20 token contract cannot be completed, because 
the gas limit for the transaction is reached, no matter what the value is set 
to. 

1. Clone repo
2. npm install
3. Reproduce ERC20 problem with `'./reproduce_problem.sh'` script
4. Run test suite `'tests/maltcoin_test.go'` -> fails
5. Successfully deploy simpler contract with `'./init.sh'` script
6. Run test suite `'tests/maltcoin_test.go'` -> passes

## Description
I have stripped the ERC20 contract (`Maltcoin.sol`) down to the bare minimum implementation 
of an OpenZeppelin contract, that can be generated with https://wizard.openzeppelin.com/, and tested it with [Remix](https://remix.ethereum.org/). The compilation and deployment works fine in Remix and fails using `go-ethereum`. However, no error is thrown by the deployment function. The problem shows itself, when trying to call the contract functions, e.g. `maltcoin.Name(nil)`. 

On the other hand, a simple Solidity contract (`Maltcoin_temp.sol`) can be deployed with the same Go functions, `solc` and `abigen` commands. Therefore, the problem is most likely in importing the ERC20 base contract from OpenZeppelin.

There are encounters of this error in some ressources online, but some were either solved with a higher gas limit, which didn't help in this case, or were solved after commiting the transaction, which is not the problem either:
- https://github.com/ethereum/go-ethereum/issues/20636
- https://github.com/ethereum/go-ethereum/issues/15930
- https://ethereum.stackexchange.com/questions/68706/no-contract-code-at-give-address
- https://ethereum.stackexchange.com/questions/19725/contract-creation-code-storage-out-of-gas

## Testing
When testing with a simulated backend (`tests/maltcoin_test.go`), an error trace shows this too: 
```shell
 > ./reproduce_problem.sh
...

 > go test tests/maltcoin_test.go 
--- FAIL: TestMaltcoin (0.00s)
    maltcoin_test.go:49: 
                Error Trace:    /Users/malte/dev/go/GoSmartContract/tests/maltcoin_test.go:49
                Error:          Received unexpected error:
                                contract creation code storage out of gas
                Test:           TestMaltcoin
                Messages:       Could not deploy contract
FAIL
```

However, when compiling the simpler `Maltcoin_temp.sol` by executing `./init.sh`, and running the test suite, all tests pass:
```shell
 > ./init.sh
...

 > go test tests/maltcoin_test.go
ok      command-line-arguments  0.185s
```

This shows, that the test suite itself is set up correctly.

