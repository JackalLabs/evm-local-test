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
    event AddedViewers(address from, string viewer_ids, string viewer_keys, string for_address, string file_owner);
    event RemovedViewers(address from, string viewer_ids, string for_address, string file_owner);
    event ResetViewers(address from, string for_address, string file_owner);
    event ChangedOwner(address from, string for_address, string file_owner, string new_owner);
    event AddedEditors(address from, string editor_ids, string editor_keys, string for_address, string file_owner);
    event RemovedEditors(address from, string editor_ids, string for_address, string file_owner);
    event ResetEditors(address from, string for_address, string file_owner);
    event CreatedNotification(address from, string to, string contents, string private_contents);
    event DeletedNotification(address from, string notification_from, uint64 time);
    event BlockedSenders(address from, string[] to_block);

    struct JackalMessage {
        string id;
        address sender;
        uint256 height;
        uint256 value;
    }

    JackalMessage[] public messages;

    function newMessage(address sender, string memory messageType) private returns (string memory) {
        uint256 height = block.number;

        // concat messageType and address
        string memory s = string.concat(messageType, Strings.toHexString(uint160(sender), 20));
        s = string.concat(s, Strings.toString(height));

        JackalMessage memory message = JackalMessage(s, sender, height, msg.value);

        messages.push(message);

        return s;
    }

    function _remove(uint256 index) internal {
        require(index < messages.length);
        messages[index] = messages[messages.length - 1];
        messages.pop();
    }

    function refund(string memory id) public {
        refundFrom(payable(msg.sender), id);
    }

    function refundFrom(address payable from, string memory id) public hasAllowance(from) {
        for (uint256 i = 0; i < messages.length; i++) {
            JackalMessage memory m = messages[i];
            if (Strings.equal(m.id, id)) {
                // if we found the item we're looking for
                if (m.height + 60 < block.number) {
                    // if 60 blocks have passed and this message is still here
                    from.transfer(m.value);
                    _remove(i);
                }
            }
        }
    }

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
        // requires chainlink oracle, will not work on localnet
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

    function postFile(string memory merkle, uint64 filesize, string memory note, uint64 expires) public payable {
        postFileFrom(msg.sender, merkle, filesize, note, expires);
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
        newMessage(from, "PostedFile");
        emit PostedFile(from, merkle, filesize, note, expires);
    }

    function buyStorage(string memory for_address, uint64 duration_days, uint64 size_bytes, string memory referral)
        public
        payable
    {
        buyStorageFrom(msg.sender, for_address, duration_days, size_bytes, referral);
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
        newMessage(from, "BoughtStorage");
        emit BoughtStorage(from, for_address, duration_days, size_bytes, referral);
    }

    function deleteFile(string memory merkle, uint64 start) public {
        deleteFileFrom(msg.sender, merkle, start);
    }

    function deleteFileFrom(address from, string memory merkle, uint64 start) public validAddress hasAllowance(from) {
        newMessage(from, "DeletedFile");
        emit DeletedFile(from, merkle, start); // file deletion is free
    }

    function requestReportForm(string memory prover, string memory merkle, string memory owner, uint64 start) public {
        requestReportFormFrom(msg.sender, prover, merkle, owner, start);
    }

    function requestReportFormFrom(
        address from,
        string memory prover,
        string memory merkle,
        string memory owner,
        uint64 start
    ) public validAddress hasAllowance(from) {
        newMessage(from, "RequestedReportForm");
        emit RequestedReportForm(from, prover, merkle, owner, start);
    }

    function postKey(string memory key) public {
        postKeyFrom(msg.sender, key);
    }

    function postKeyFrom(address from, string memory key) public validAddress hasAllowance(from) {
        newMessage(from, "PostedKey");
        emit PostedKey(from, key);
    }

    function deleteFileTree(string memory hash_path, string memory account) public {
        deleteFileTreeFrom(msg.sender, hash_path, account);
    }

    function deleteFileTreeFrom(address from, string memory hash_path, string memory account)
        public
        validAddress
        hasAllowance(from)
    {
        newMessage(from, "DeletedFileTree");
        emit DeletedFileTree(from, hash_path, account);
    }

    function provisionFileTree(string memory editors, string memory viewers, string memory tracking_number) public {
        provisionFileTreeFrom(msg.sender, editors, viewers, tracking_number);
    }

    function provisionFileTreeFrom(
        address from,
        string memory editors,
        string memory viewers,
        string memory tracking_number
    ) public validAddress hasAllowance(from) {
        newMessage(from, "ProvisionedFileTree");
        emit ProvisionedFileTree(from, editors, viewers, tracking_number);
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
        newMessage(from, "PostedFileTree");
        emit PostedFileTree(from, account, hash_parent, hash_child, contents, viewers, editors, tracking_number);
    }

    function addViewers(
        string memory viewer_ids,
        string memory viewer_keys,
        string memory for_address,
        string memory file_owner
    ) public {
        addViewersFrom(msg.sender, viewer_ids, viewer_keys, for_address, file_owner);
    }

    function addViewersFrom(
        address from,
        string memory viewer_ids,
        string memory viewer_keys,
        string memory for_address,
        string memory file_owner
    ) public validAddress hasAllowance(from) {
        newMessage(from, "AddedViewers");
        emit AddedViewers(from, viewer_ids, viewer_keys, for_address, file_owner);
    }

    function removeViewers(string memory viewer_ids, string memory for_address, string memory file_owner) public {
        removeViewersFrom(msg.sender, viewer_ids, for_address, file_owner);
    }

    function removeViewersFrom(
        address from,
        string memory viewer_ids,
        string memory for_address,
        string memory file_owner
    ) public validAddress hasAllowance(from) {
        newMessage(from, "RemovedViewers");
        emit RemovedViewers(from, viewer_ids, for_address, file_owner);
    }

    function resetViewers(string memory for_address, string memory file_owner) public {
        resetViewersFrom(msg.sender, for_address, file_owner);
    }

    function resetViewersFrom(address from, string memory for_address, string memory file_owner)
        public
        validAddress
        hasAllowance(from)
    {
        newMessage(from, "ResetViewers");
        emit ResetViewers(from, for_address, file_owner);
    }

    function changeOwner(string memory for_address, string memory file_owner, string memory new_owner) public {
        changeOwnerFrom(msg.sender, for_address, file_owner, new_owner);
    }

    function changeOwnerFrom(address from, string memory for_address, string memory file_owner, string memory new_owner)
        public
        validAddress
        hasAllowance(from)
    {
        newMessage(from, "ChangedOwner");
        emit ChangedOwner(from, for_address, file_owner, new_owner);
    }

    function addEditors(
        string memory editor_ids,
        string memory editor_keys,
        string memory for_address,
        string memory file_owner
    ) public {
        addEditorsFrom(msg.sender, editor_ids, editor_keys, for_address, file_owner);
    }

    function addEditorsFrom(
        address from,
        string memory editor_ids,
        string memory editor_keys,
        string memory for_address,
        string memory file_owner
    ) public validAddress hasAllowance(from) {
        newMessage(from, "AddedEditors");
        emit AddedEditors(from, editor_ids, editor_keys, for_address, file_owner);
    }

    function removeEditors(string memory editor_ids, string memory for_address, string memory file_owner) public {
        removeEditorsFrom(msg.sender, editor_ids, for_address, file_owner);
    }

    function removeEditorsFrom(
        address from,
        string memory editor_ids,
        string memory for_address,
        string memory file_owner
    ) public validAddress hasAllowance(from) {
        newMessage(from, "RemovedEditors");
        emit RemovedEditors(from, editor_ids, for_address, file_owner);
    }

    function resetEditors(string memory for_address, string memory file_owner) public {
        resetEditorsFrom(msg.sender, for_address, file_owner);
    }

    function resetEditorsFrom(address from, string memory for_address, string memory file_owner)
        public
        validAddress
        hasAllowance(from)
    {
        newMessage(from, "ResetEditors");
        emit ResetEditors(from, for_address, file_owner);
    }

    function createNotification(string memory to, string memory contents, string memory private_contents) public {
        createNotificationFrom(msg.sender, to, contents, private_contents);
    }

    function createNotificationFrom(
        address from,
        string memory to,
        string memory contents,
        string memory private_contents
    ) public validAddress hasAllowance(from) {
        newMessage(from, "CreatedNotification");
        emit CreatedNotification(from, to, contents, private_contents);
    }

    function deleteNotification(string memory notification_from, uint64 time) public {
        deleteNotificationFrom(msg.sender, notification_from, time);
    }

    function deleteNotificationFrom(address from, string memory notification_from, uint64 time)
        public
        validAddress
        hasAllowance(from)
    {
        newMessage(from, "DeletedNotification");
        emit DeletedNotification(from, notification_from, time);
    }

    function blockSenders(string[] memory to_block) public {
        blockSendersFrom(msg.sender, to_block);
    }

    function blockSendersFrom(address from, string[] memory to_block) public validAddress hasAllowance(from) {
        newMessage(from, "BlockedSenders");
        emit BlockedSenders(from, to_block);
    }
}
