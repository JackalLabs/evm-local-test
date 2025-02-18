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
    function deleteFileFrom(address from, string memory merkle, uint64 start) external;
    function deleteFile(string memory merkle, uint64 start) external;
    function requestReportFormFrom(
        address from,
        string memory prover,
        string memory merkle,
        string memory owner,
        uint64 start
    ) external;
    function requestReportForm(string memory prover, string memory merkle, string memory owner, uint64 start)
        external;
    function postKey(string memory key) external;
    function postKeyFrom(address from, string memory key) external;
    function deleteFileTree(string memory hash_path, string memory account) external;
    function deleteFileTreeFrom(address from, string memory hash_path, string memory account) external;
    function provisionFileTree(string memory editors, string memory viewers, string memory tracking_number) external;
    function provisionFileTreeFrom(
        address from,
        string memory editors,
        string memory viewers,
        string memory tracking_number
    ) external;
    function postFileTree(
        string memory account,
        string memory hash_parent,
        string memory hash_child,
        string memory contents,
        string memory viewers,
        string memory editors,
        string memory tracking_number
    ) external;
    function postFileTreeFrom(
        address from,
        string memory account,
        string memory hash_parent,
        string memory hash_child,
        string memory contents,
        string memory viewers,
        string memory editors,
        string memory tracking_number
    ) external;
    function addViewers(
        string memory viewer_ids,
        string memory viewer_keys,
        string memory for_address,
        string memory file_owner
    ) external;
    function addViewersFrom(
        address from,
        string memory viewer_ids,
        string memory viewer_keys,
        string memory for_address,
        string memory file_owner
    ) external;
    function removeViewers(string memory viewer_ids, string memory for_address, string memory file_owner) external;
    function removeViewersFrom(
        address from,
        string memory viewer_ids,
        string memory for_address,
        string memory file_owner
    ) external;
    function resetViewers(string memory for_address, string memory file_owner) external;
    function resetViewersFrom(address from, string memory for_address, string memory file_owner) external;
    function changeOwner(string memory for_address, string memory file_owner, string memory new_owner) external;
    function changeOwnerFrom(address from, string memory for_address, string memory file_owner, string memory new_owner)
        external;
    function addEditors(
        string memory editor_ids,
        string memory editor_keys,
        string memory for_address,
        string memory file_owner
    ) external;
    function addEditorsFrom(
        address from,
        string memory editor_ids,
        string memory editor_keys,
        string memory for_address,
        string memory file_owner
    ) external;
    function removeEditors(string memory editor_ids, string memory for_address, string memory file_owner) external;
    function removeEditorsFrom(
        address from,
        string memory editor_ids,
        string memory for_address,
        string memory file_owner
    ) external;
    function resetEditors(string memory for_address, string memory file_owner) external;
    function resetEditorsFrom(address from, string memory for_address, string memory file_owner) external;
}
