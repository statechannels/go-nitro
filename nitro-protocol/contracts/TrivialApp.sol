// SPDX-License-Identifier: MIT
pragma solidity 0.7.6;
pragma experimental ABIEncoderV2;

import './interfaces/IForceMoveApp.sol';

/**
 * @dev The TrivialApp contracts complies with the ForceMoveApp interface and allows all transitions, regardless of the data. Used for testing purposes.
 */
contract TrivialApp is IForceMoveApp {
    /**
     * @notice Encodes trivial rules.
     * @dev Encodes trivial rules.
     * @return last variable part.
     */
    function latestSupportedState(
        FixedPart calldata, // fixedPart, unused
        SignedVariablePart[] calldata signedVariableParts
    ) external pure override returns (uint256) {
        return signedVariableParts.length - 1;
    }
}
