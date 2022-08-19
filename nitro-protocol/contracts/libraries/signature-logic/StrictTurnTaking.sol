// SPDX-License-Identifier: MIT
pragma solidity 0.7.6;
pragma experimental ABIEncoderV2;

import {ExitFormat as Outcome} from '@statechannels/exit-format/contracts/ExitFormat.sol';
import {NitroUtils} from '../NitroUtils.sol';
import '../../interfaces/INitroTypes.sol';

library StrictTurnTaking {
    /**
     * @notice Require supplied arguments to comply with turn taking logic, i.e. each participant signed the one state, they were mover for.
     * @dev Require supplied arguments to comply with turn taking logic, i.e. each participant signed the one state, they were mover for.
     * @param fixedPart Data describing properties of the state channel that do not change with state updates.
     * @param proof Array of recovered variable parts which constitutes a support proof for the candidate. The proof is a validation for the supplied candidate.
     * @param candidate Recovered variable part the proof was supplied for. The candidate state is supported by proof states.
     */
    function requireValidTurnTaking(
        INitroTypes.FixedPart memory fixedPart,
        INitroTypes.RecoveredVariablePart[] memory proof,
        INitroTypes.RecoveredVariablePart memory candidate
    ) internal pure {
        _requireValidInput(fixedPart.participants.length, proof.length + 1);

        uint48 turnNum = proof[0].variablePart.turnNum;

        for (uint256 i = 0; i < proof.length; i++) {
            isSignedByMover(fixedPart, (i < proof.length ? proof[i] : candidate));
            requireHasTurnNum(
                i < proof.length ? proof[i].variablePart : candidate.variablePart,
                turnNum
            );
            turnNum++;
        }
    }

    /**
     * @notice Require supplied state is signed by its corresponding moving participant.
     * @dev Require supplied state is signed by its corresponding moving participant.
     * @param fixedPart Data describing properties of the state channel that do not change with state updates.
     * @param recoveredVariablePart A struct describing dynamic properties of the state channel, that must be signed by moving participant.
     */
    function isSignedByMover(
        INitroTypes.FixedPart memory fixedPart,
        INitroTypes.RecoveredVariablePart memory recoveredVariablePart
    ) internal pure {
        require(
            NitroUtils.isClaimedSignedOnlyBy(
                recoveredVariablePart.signedBy,
                uint8(recoveredVariablePart.variablePart.turnNum % fixedPart.participants.length)
            ),
            'Invalid signedBy'
        );
    }

    /**
     * @notice Require supplied variable part has specified turn number.
     * @dev Require supplied variable part has specified turn number.
     * @param variablePart Variable part to check turn number of.
     * @param turnNum Turn number to compare with.
     */
    function requireHasTurnNum(INitroTypes.VariablePart memory variablePart, uint48 turnNum)
        internal
        pure
    {
        require(variablePart.turnNum == turnNum, 'Wrong variablePart.turnNum');
    }

    /**
     * @notice Find moving participant address based on state turn number.
     * @dev Find moving participant address based on state turn number.
     * @param participants Array of participant addresses.
     * @param turnNum State turn number.
     * @return address Moving partitipant address.
     */
    function _moverAddress(address[] memory participants, uint48 turnNum)
        internal
        pure
        returns (address)
    {
        return participants[turnNum % participants.length];
    }

    /**
     * @notice Validate input for turn taking logic.
     * @dev Validate input for turn taking logic.
     * @param numParticipants Number of participants in a channel.
     * @param numStates Number of states submitted.
     */
    function _requireValidInput(uint256 numParticipants, uint256 numStates) internal pure {
        require((numParticipants >= numStates) && (numStates > 1), 'Insufficient or excess states');

        // no more than 255 participants
        require(numParticipants <= type(uint8).max, 'Too many participants'); // type(uint8).max = 2**8 - 1 = 255
    }
}
