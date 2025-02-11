// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

interface JackalInterface {
    function postFileFrom(address from, string memory merkle, uint64 filesize, string memory note, uint64 expires)
        external
        payable;
    function postFile(string memory merkle, uint64 filesize, string memory note, uint64 expires) external payable;
    function buyStorageFrom(
        address from,
        string memory for_address,
        uint64 duration_days,
        uint64 size_bytes,
        string memory referral
    ) external payable;
    function buyStorage(string memory for_address, uint64 duration_days, uint64 size_bytes, string memory referral)
        external
        payable;
    function deleteFileFrom(address from, string memory merkle, uint64 start) external payable;
    function deleteFile(string memory merkle, uint64 start) external payable;
    function requestReportFormFrom(
        address from,
        string memory prover,
        string memory merkle,
        string memory owner,
        uint64 start
    ) external payable;
    function requestReportForm(string memory prover, string memory merkle, string memory owner, uint64 start)
        external
        payable;
    function postKey(string memory key) external payable;
    function postKeyFrom(address from, string memory key) external payable;
    function deleteFileTree(string memory hash_path, string memory account) external payable;
    function deleteFileTreeFrom(address from, string memory hash_path, string memory account) external payable;
}
