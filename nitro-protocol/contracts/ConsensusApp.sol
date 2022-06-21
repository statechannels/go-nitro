// SPDX-License-Identifier: MIT
pragma solidity 0.7.6;
pragma experimental ABIEncoderV2;

import './interfaces/IForceMoveApp.sol';
import './libraries/NitroUtils.sol';

/**
 * @dev The ConsensusApp contracts complies with the ForceMoveApp interface and requires all participants to hace signed a single state.
 */
contract ConsensusApp is IForceMoveApp {
    /**
     * @notice Encodes consensus app rules.
     * @dev Encodes consensus app rules.
     * @return last variable part.
     */
    function latestSupportedState(
        FixedPart calldata fixedPart,
        SignedVariablePart[] calldata signedVariableParts
    ) external pure override returns (VariablePart memory) {
        require(signedVariableParts.length == 1, '|signedVariableParts|!=1');
        require(NitroUtils.getSignersAmount(signedVariableParts[0].signedBy) == fixedPart.participants.length, 'require everyone to sign');
        return signedVariableParts[0].variablePart;
    }
}
