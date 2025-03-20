// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import {JackalInterface} from "./JackalInterface.sol";

contract StorageDrawer {
    JackalInterface internal jackalBridge;

    mapping(address => string[]) public cabinet;

    constructor(address _jackalAddress) {
        jackalBridge = JackalInterface(_jackalAddress);
    }

    function upload(string memory merkle, uint64 filesize) public payable {
        jackalBridge.postFileFrom{value: msg.value}(msg.sender, merkle, filesize, "", 30);
        cabinet[msg.sender].push(merkle);
    }

    function fileAddress(address _addr, uint256 _index) public view returns (string memory) {
        require(_index < cabinet[_addr].length);
        return cabinet[_addr][_index];
    }

    function fileCount(address _addr) public view returns (uint256) {
        return cabinet[_addr].length;
    }
}
