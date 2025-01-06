// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

interface JackalInterface {
    function postFile(string memory merkle, uint64 filesize) external payable;
}

contract StorageDrawer {

    JackalInterface internal jackalBridge;

    mapping(address => string[]) public cabinet;

    constructor(address _jackalAddress){
        jackalBridge = JackalInterface(_jackalAddress);
    }

    function upload(string memory merkle, uint64 filesize) public payable {
        jackalBridge.postFile{value: msg.value}(merkle, filesize);
        cabinet[msg.sender].push(merkle);
    }

    function fileAddress(address _addr, uint _index) public view returns (string memory){
        require(_index < cabinet[_addr].length);
        return cabinet[_addr][_index];
    }

    function fileCount(address _addr) public view returns (uint){
        return cabinet[_addr].length;
    }
}