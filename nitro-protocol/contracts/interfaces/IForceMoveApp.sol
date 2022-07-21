// SPDX-License-Identifier: MIT
pragma solidity 0.7.6;
pragma experimental ABIEncoderV2;

import './INitroTypes.sol';

/**
 * @dev The IForceMoveApp interface calls for its children to implement an application-specific latestSupportedState function, defining the state machine of a ForceMove state channel DApp.
 */
interface IForceMoveApp is INitroTypes {
    /**
     * @notice Encodes application-specific rules for a particular ForceMove-compliant state channel.
     * @dev Encodes application-specific rules for a particular ForceMove-compliant state channel.
     * @param fixedPart Fixed part of the state channel.
     * @param recoveredVariableParts Array of recovered variable parts to find the latest of.
     * @return VariablePart Latest supported by application variable part from supplied array.
     */
    function latestSupportedState(
        FixedPart calldata fixedPart,
        RecoveredVariablePart[] calldata recoveredVariableParts
    ) external pure returns (VariablePart memory);
}
