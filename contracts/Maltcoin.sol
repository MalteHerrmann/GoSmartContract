// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

// Import IERC20 interface from OpenZeppelin contracts library
import "node_modules/@openzeppelin/contracts/token/ERC20/ERC20.sol";

// Define the contract
contract Maltcoin is ERC20 {
    // Constructor
    constructor() ERC20("Maltcoin", "Malt") {
        _mint(msg.sender, 10000 * 10 ** decimals());
    }
}