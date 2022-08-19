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
     * @notice Encodes application-specific rules for a particular ForceMove-compliant state channel.
     * @dev Encodes application-specific rules for a particular ForceMove-compliant state channel.
     * @param fixedPart Fixed part of the state channel.
     * @param proof Array of recovered variable parts which constitutes a support proof for the candidate.
     * @param candidate Recovered variable part the proof was supplied for.
     */
    function requireStateSupported(
        FixedPart calldata fixedPart,
        RecoveredVariablePart[] calldata proof,
        RecoveredVariablePart calldata candidate
    ) external pure override {
        require(proof.length == 0, '|proof|!=0');
        require(
            NitroUtils.getClaimedSignersNum(candidate.signedBy) == fixedPart.participants.length,
            '!unanimous'
        );
    }
}
