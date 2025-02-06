// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

interface JackalInterface {
    function postFile(string memory merkle, uint64 filesize, string memory note, uint64 expires) external payable;
    function postFileFrom(address from, string memory merkle, uint64 filesize, string memory note, uint64 expires)
        external
        payable;
    function buyStorage(string memory for_address, uint64 duration_days, uint64 size_bytes, string memory referral)
        external
        payable;
    function buyStorageFrom(
        address from,
        string memory for_address,
        uint64 duration_days,
        uint64 size_bytes,
        string memory referral
    ) external payable;
    function deleteFile(string memory merkle, uint64 start) external payable;
    function deleteFileFrom(address from, string memory merkle, uint64 start) external payable;
}
