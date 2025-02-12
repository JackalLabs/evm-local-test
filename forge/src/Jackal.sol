// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {Strings} from "@openzeppelin/contracts/utils/Strings.sol";

abstract contract Jackal {
    event PostedFile(address from, string merkle, uint64 size, string note, uint64 expires);
    event BoughtStorage(address from, string for_address, uint64 duration_days, uint64 size_bytes, string referral);
    event DeletedFile(address from, string merkle, uint64 start);
    event RequestedReportForm(address from, string prover, string merkle, string owner, uint64 start);
    event PostedKey(address from, string key);
    event DeletedFileTree(address from, string hash_path, string account);
    event ProvisionedFileTree(address from, string editors, string viewers, string tracking_number);
    event PostedFileTree(
        address from,
        string account,
        string hash_parent,
        string hash_child,
        string contents,
        string viewers,
        string editors,
        string tracking_number
    );

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

    function deleteFileFrom(address from, string memory merkle, uint64 start) public validAddress hasAllowance(from) {
        emit DeletedFile(from, merkle, start); // file deletion is free
    }

    function deleteFile(string memory merkle, uint64 start) public {
        deleteFileFrom(msg.sender, merkle, start);
    }

    function requestReportFormFrom(
        address from,
        string memory prover,
        string memory merkle,
        string memory owner,
        uint64 start
    ) public validAddress hasAllowance(from) {
        emit RequestedReportForm(from, prover, merkle, owner, start);
    }

    function requestReportForm(string memory prover, string memory merkle, string memory owner, uint64 start) public {
        requestReportFormFrom(msg.sender, prover, merkle, owner, start);
    }

    function postKeyFrom(address from, string memory key) public validAddress hasAllowance(from) {
        emit PostedKey(from, key);
    }

    function postKey(string memory key) public {
        postKeyFrom(msg.sender, key);
    }

    function deleteFileTreeFrom(address from, string memory hash_path, string memory account)
        public
        validAddress
        hasAllowance(from)
    {
        emit DeletedFileTree(from, hash_path, account);
    }

    function deleteFileTree(string memory hash_path, string memory account) public {
        deleteFileTreeFrom(msg.sender, hash_path, account);
    }

    function provisionFileTreeFrom(
        address from,
        string memory editors,
        string memory viewers,
        string memory tracking_number
    ) public validAddress hasAllowance(from) {
        emit ProvisionedFileTree(from, editors, viewers, tracking_number);
    }

    function provisionFileTree(string memory editors, string memory viewers, string memory tracking_number) public {
        provisionFileTreeFrom(msg.sender, editors, viewers, tracking_number);
    }

    function postFileTreeFrom(
        address from,
        string memory account,
        string memory hash_parent,
        string memory hash_child,
        string memory contents,
        string memory viewers,
        string memory editors,
        string memory tracking_number
    ) public validAddress hasAllowance(from) {
        // confirm postFileTree is free
        emit PostedFileTree(from, account, hash_parent, hash_child, contents, viewers, editors, tracking_number);
    }

    function postFileTree(
        string memory account,
        string memory hash_parent,
        string memory hash_child,
        string memory contents,
        string memory viewers,
        string memory editors,
        string memory tracking_number
    ) public {
        postFileTreeFrom(msg.sender, account, hash_parent, hash_child, contents, viewers, editors, tracking_number);
    }
}
