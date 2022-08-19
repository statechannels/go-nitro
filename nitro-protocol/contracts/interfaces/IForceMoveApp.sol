// SPDX-License-Identifier: MIT
pragma solidity 0.7.6;
pragma experimental ABIEncoderV2;

import './INitroTypes.sol';

/**
 * @dev The IForceMoveApp interface calls for its children to implement an application-specific requireStateSupported function, defining the state machine of a ForceMove state channel DApp.
 */
interface IForceMoveApp is INitroTypes {
    /**
     * @notice Encodes application-specific rules for a particular ForceMove-compliant state channel. Must revert when invalid support proof and a candidate are supplied.
     * @dev Encodes application-specific rules for a particular ForceMove-compliant state channel. Must revert when invalid support proof and a candidate are supplied.
     * @param fixedPart Fixed part of the state channel.
     * @param proof Array of recovered variable parts which constitutes a support proof for the candidate. May be omitted when `candidate` constitutes a support proof itself.
     * @param candidate Recovered variable part the proof was supplied for. Also may constitute a support proof itself.
     */
    function requireStateSupported(
        FixedPart calldata fixedPart,
        RecoveredVariablePart[] calldata proof,
        RecoveredVariablePart calldata candidate
    ) external pure;
}
