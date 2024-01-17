// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

contract CalleeContract {
    uint256 x = 0;
    event ErrorMessage(address indexed sender, string indexed errorMessage);
    event LogMessage(address indexed sender, uint256);

    function EmitWhenErrorOccurred() public {
        emit ErrorMessage(msg.sender, "Error Occurred!");
        revert();
    }

    function EmitWithRead() public {
        emit LogMessage(msg.sender, x);
    }

    function EmitWithWrite() public returns(uint256){
        x += 1;
        emit LogMessage(msg.sender, x);
        return x;
    }
}