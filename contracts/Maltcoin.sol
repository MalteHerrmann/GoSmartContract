// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

// import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "node_modules/@openzeppelin/contracts/token/ERC20/ERC20.sol";

/// @title Maltcoin
/// @author Malte Herrmann
/// @notice This contract defines an ERC20 token called Maltcoin
/** @dev This contract was generated using the OpenZeppelin contract 
wizard: https://wizard.openzeppelin.com/
*/ 
contract Maltcoin is ERC20 {
    /** @notice The constructor function is called upon deployment of the
    contract. It initializes the contract with the name and symbol 
    of the token.
    */ 
    /// @dev 10.000 tokens are minted and assigned to the transaction sender.
    constructor() ERC20("Maltcoin", "MALT") {
        _mint(msg.sender, 10000 * 10 ** decimals());
    }
}