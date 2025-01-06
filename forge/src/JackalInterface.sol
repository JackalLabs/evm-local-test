// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

interface JackalInterface {
    function postFile(string memory merkle, uint64 filesize) external payable;
    function postFileFrom(address from, string memory merkle, uint64 filesize) external payable;
}
