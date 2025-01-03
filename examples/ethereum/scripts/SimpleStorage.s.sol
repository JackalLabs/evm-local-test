// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import "forge-std/Script.sol"; // Import Foundry's Script utilities

// Simple Storage Contract
contract SimpleStorage {
    uint256 private storedData;

    // Set the stored data
    function set(uint256 x) public {
        storedData = x;
    }

    // Get the stored data
    function get() public view returns (uint256) {
        return storedData;
    }
}

// Script for deploying and interacting with SimpleStorage
contract SimpleStorageScript is Script {
    function run() public {
        // Start broadcasting transactions
        vm.startBroadcast();

        // Deploy the SimpleStorage contract
        SimpleStorage simpleStorage = new SimpleStorage();

        // Interact with the deployed contract (optional)
        simpleStorage.set(42);
        uint256 retrievedValue = simpleStorage.get();
        require(retrievedValue == 42, "Storage mismatch");

        // Stop broadcasting transactions
        vm.stopBroadcast();

        // Log results
        console.log("SimpleStorage deployed and verified. Stored value:", retrievedValue);
    }
}
