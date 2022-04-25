// SPDX-License-Identifier: MIT
pragma solidity 0.7.6;
pragma experimental ABIEncoderV2;

import './INitroTypes.sol';

/**
 * @dev The IForceMoveApp interface calls for its children to implement an application-specific validTransition function, defining the state machine of a ForceMove state channel DApp.
 */
interface IForceMoveApp is INitroTypes {
    
    /**
     * @notice Encodes application-specific rules for a particular ForceMove-compliant state channel.
     * @dev Encodes application-specific rules for a particular ForceMove-compliant state channel.
     * @param a State being transitioned from.
     * @param b State being transitioned to.
     * @param nParticipants Number of participants in this state channel.
     * @return true if the transition conforms to this application's rules, false otherwise
     */
    function validTransition(
        VariablePart calldata a,
        VariablePart calldata b,
        uint256 nParticipants
    ) external pure returns (bool);
}
