// SPDX-License-Identifier: MIT
pragma solidity 0.8.17;

contract SimpleCoin {
    mapping(address => uint256) internal balances;

    event Transfer(address indexed _from, address indexed _to, uint256 _value);

    constructor() {
        balances[msg.sender] = 10000;
    }

    function sendCoin(address receiver, uint256 amount) public returns (bool sufficient) {
        if (balances[msg.sender] < amount) return false;
        balances[msg.sender] -= amount;
        balances[receiver] += amount;
        return true;
    }

    function getBalanceInEth(address addr) public view returns (uint256) {
        return getBalance(addr) * 2;
    }

    function getBalance(address addr) public view returns (uint256) {
        return balances[addr];
    }
}
