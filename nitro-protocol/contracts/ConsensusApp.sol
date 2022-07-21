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
        RecoveredVariablePart[] calldata recoveredVariableParts
    ) external pure override returns (VariablePart memory) {
        require(recoveredVariableParts.length == 1, '|signedVariableParts|!=1');
        require(NitroUtils.getClaimedSignersNum(recoveredVariableParts[0].signedBy) == fixedPart.participants.length, '!unanimous');
        return recoveredVariableParts[0].variablePart;
    }
}
