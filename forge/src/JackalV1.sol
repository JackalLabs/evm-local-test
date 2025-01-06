// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {AggregatorV3Interface} from "@chainlink/interfaces/feeds/AggregatorV3Interface.sol";
import {Jackal} from "./Jackal.sol";

contract JackalBridge is Ownable, Jackal {

    AggregatorV3Interface internal priceFeed;

    address[] public relays;

    constructor(address[] memory _relays, address _priceFeed) Ownable(msg.sender){
        require(_relays.length > 0, "must provide relays");

        priceFeed = AggregatorV3Interface(_priceFeed);
        relays = _relays;
    }

    // Modifier to restrict access to owner or relays
    modifier onlyOwnerOrRelay() {
        require(msg.sender == owner() || isRelay(msg.sender), "not owner or relay");
        _;
    }

    function isRelay(address _relay) internal view returns (bool) {
        for (uint i = 0; i < relays.length; i++) {
            if (relays[i] == _relay) {
                return true;
            }
        }
        return false;
    }

    // Function to add a relay, only callable by the owner
    function addRelay(address _relay) public onlyOwner {
        relays.push(_relay);
    }

    // Function to remove a relay, only callable by the owner
    function removeRelay(address _relay) public onlyOwner {

        require(relays.length > 1); // require there to be at least one relay in the list after removal

        for (uint i = 0; i < relays.length; i++) {
            if (relays[i] == _relay) {
                relays[i] = relays[relays.length - 1];
                relays.pop();
                break;
            }
        }

    }



    function distributeBalance() public onlyOwnerOrRelay {
        uint256 balance = address(this).balance;
        require(balance > 2, "not enough wei to distribute");

        uint256 ownerShare = balance / 2;
        payable(owner()).transfer(ownerShare);

        uint256 relayShare = balance - ownerShare; // Remaining 50%
        uint256 perRelay = relayShare / relays.length;

        for (uint i = 0; i < relays.length; i++) {
            payable(relays[i]).transfer(perRelay);
        }

    }

    function getPrice() public override view returns (int) {
        (,int price,,,) = priceFeed.latestRoundData();
        return price;
    }


}