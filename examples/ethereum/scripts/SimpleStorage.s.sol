// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

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

// Testing Contract
contract SimpleStorageTest {
    SimpleStorage public simpleStorage;

    // Constructor to set up the test
    constructor() {
        simpleStorage = new SimpleStorage();
    }

    // Test the set and get functions
    function testSetAndGet() public returns (bool) {
        simpleStorage.set(42);
        require(simpleStorage.get() == 42, "Storage mismatch");
        return true;
    }
}