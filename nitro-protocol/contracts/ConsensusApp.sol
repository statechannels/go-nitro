// SPDX-License-Identifier: MIT
pragma solidity 0.7.6;
pragma experimental ABIEncoderV2;

import {ShortcuttingTurnTaking} from './libraries/signature-logic/ShortcuttingTurnTaking.sol';
import './interfaces/IForceMoveApp.sol';
import './libraries/NitroUtils.sol';

/**
 * @dev The ConsensusApp contracts complies with the ForceMoveApp interface and requires all participants to have signed a single state.
 */
contract ConsensusApp is IForceMoveApp {
    /**
     * @notice Encodes consensus app rules.
     * @dev Encodes consensus app rules.
     * @return lastest supported state.
     */
    function latestSupportedState(
        FixedPart calldata fixedPart,
        SignedVariablePart[] calldata signedVariableParts
    ) external pure override returns (VariablePart memory) {
        require(signedVariableParts.length == 1, '|signedVariableParts|!=1');
        ShortcuttingTurnTaking.requireValidTurnTaking(fixedPart, signedVariableParts);
        return signedVariableParts[0].variablePart;
    }
}
