// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.26;

import {Test, console} from "forge-std/Test.sol";
import {JackalBridge} from "../src/JackalV1.sol";

contract CounterTest is Test {
    JackalBridge public bridge;

    function setUp() public {
        address[] memory t = new address[](1);
        t[0] = 0x9443A8C2aa7788EEE05f9734Ad4174a6C5CA0A25;

        address priceFeed = 0x9326BFA02ADD2366b30bacB125260Af641031331;

        bridge = new JackalBridge(t, priceFeed);
    }

}
