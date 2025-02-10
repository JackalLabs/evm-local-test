// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {Strings} from "@openzeppelin/contracts/utils/Strings.sol";

abstract contract Jackal {
    event PostedFile(address from, string merkle, uint64 size, string note, uint64 expires);
    event BoughtStorage(address from, string for_address, uint64 duration_days, uint64 size_bytes, string referral);
    event DeletedFile(address from, string merkle, uint64 start);
    event RequestedReportForm(address from, string prover, string merkle, string owner, uint64 start);
    event PostedKey(address from, string key);

    function getPrice() public view virtual returns (int256);

    mapping(address => mapping(address => bool)) public allowances;

    modifier validAddress() {
        require(msg.sender != address(0), "Invalid sender address");
        _;
    }

    modifier hasAllowance(address from) {
        require(getAllowance(msg.sender, from), "No allowance set for contract");
        _;
    }

    function getAllowance(address from, address to) public view returns (bool) {
        if (from == to) {
            return true;
        }
        return allowances[to][from];
    }

    function addAllowance(address allow) public {
        allowances[msg.sender][allow] = true;
    }

    function removeAllowance(address allow) public {
        allowances[msg.sender][allow] = false;
    }

    function getStoragePrice(uint64 filesize, uint256 months) public view returns (uint256) {
        uint256 price = uint256(getPrice());
        uint256 storagePrice = 15; // price at 8 decimal places
        uint256 multiplier = 2;

        uint256 fs = filesize;
        if (fs <= 1024 * 1024) {
            fs = 1024 * 1024; // minimum file size of one MB for pricing
        }

        // Calculate the price in wei
        // 1e8 adjusts for the 8 decimals of USD, 1e18 converts ETH to wei
        uint256 BSM = storagePrice * multiplier * months * fs;
        uint256 p = (BSM * 1e8 * 1e18) / (price * 1099511627776);

        if (p == 0) {
            p = 5000 gwei;
        }

        return p;
    }

    function postFileFrom(address from, string memory merkle, uint64 filesize, string memory note, uint64 expires)
        public
        payable
        validAddress
        hasAllowance(from)
    {
        require(expires >= 30 || expires == 0);
        if (expires != 0) {
            uint256 pE = 0; // getStoragePrice(filesize, 2400); // 12 * 200 months
            require(msg.value >= pE, string.concat("Insufficient payment, need ", Strings.toString(pE), " wei"));
        }
        emit PostedFile(from, merkle, filesize, note, expires);
    }

    function postFile(string memory merkle, uint64 filesize, string memory note, uint64 expires) public payable {
        postFileFrom(msg.sender, merkle, filesize, note, expires);
    }

    function buyStorageFrom(
        address from,
        string memory for_address,
        uint64 duration_days,
        uint64 size_bytes,
        string memory referral
    ) public payable validAddress hasAllowance(from) {
        uint256 pE = 0; // getStoragePrice(size_bytes, duration_days / 30); // months
        require(msg.value >= pE, string.concat("Insufficient payment, need ", Strings.toString(pE), " wei"));
        emit BoughtStorage(from, for_address, duration_days, size_bytes, referral);
    }

    function buyStorage(string memory for_address, uint64 duration_days, uint64 size_bytes, string memory referral)
        public
        payable
    {
        buyStorageFrom(msg.sender, for_address, duration_days, size_bytes, referral);
    }

    function deleteFileFrom(address from, string memory merkle, uint64 start)
        public
        payable
        validAddress
        hasAllowance(from)
    {
        // is file deletion free?
        emit DeletedFile(from, merkle, start);
    }

    function deleteFile(string memory merkle, uint64 start) public payable {
        deleteFileFrom(msg.sender, merkle, start);
    }

    function requestReportFormFrom(
        address from,
        string memory prover,
        string memory merkle,
        string memory owner,
        uint64 start
    ) public payable validAddress hasAllowance(from) {
        emit RequestedReportForm(from, prover, merkle, owner, start);
    }

    function requestReportForm(string memory prover, string memory merkle, string memory owner, uint64 start)
        public
        payable
    {
        requestReportFormFrom(msg.sender, prover, merkle, owner, start);
    }

    function postKeyFrom(address from, string memory key) public payable validAddress hasAllowance(from) {
        emit PostedKey(from, key);
    }

    function postKey(string memory key) public payable {
        postKeyFrom(msg.sender, key);
    }
}
