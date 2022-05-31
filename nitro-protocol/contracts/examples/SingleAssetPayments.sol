// SPDX-License-Identifier: MIT
pragma solidity 0.7.6;
pragma experimental ABIEncoderV2;

import {ExitFormat as Outcome} from '@statechannels/exit-format/contracts/ExitFormat.sol';
import '../interfaces/IForceMoveApp.sol';
import {TurnTaking} from '../examples/signature-logic/TurnTaking.sol';

/**
 * @dev The SingleAssetPayments contract complies with the ForceMoveApp and TurnTaking interfaces and implements a simple payment channel with a single asset type only.
 */
contract SingleAssetPayments is IForceMoveApp {
    /**
     * @notice Encodes application-specific rules for a particular ForceMove-compliant state channel.
     * @dev Encodes application-specific rules for a particular ForceMove-compliant state channel.
     * @param fixedPart Fixed part of the state channel.
     * @param signedVariableParts Array of variable parts to find the latest of.
     * @return VariablePart Latest supported by application variable part from supplied array.
     */    
    function latestSupportedState(
        FixedPart calldata fixedPart,
        SignedVariablePart[] calldata signedVariableParts
    ) external pure override returns (VariablePart memory) {
        TurnTaking.requireValidTurnTaking(fixedPart, signedVariableParts);

        for (uint256 i = 0; i < signedVariableParts.length; i++) {
            _requireValidOutcome(fixedPart.participants.length, signedVariableParts[i].variablePart.outcome);

            if (i > 0) {
                _requireValidTransition(
                    fixedPart.participants.length,
                    signedVariableParts[i-1].variablePart,
                    signedVariableParts[i].variablePart
                );
            }
        }

        return signedVariableParts[signedVariableParts.length - 1].variablePart;
    }

    /**
     * @notice Require specific rules in outcome are followed.
     * @dev Require specific rules in outcome are followed.
     * @param nParticipants Number of participants in a channel.
     * @param outcome Outcome to check.
     */
    function _requireValidOutcome(
        uint256 nParticipants,
        Outcome.SingleAssetExit[] memory outcome
    ) internal pure {
        // Throws if more than one asset
        require(outcome.length == 1, 'outcome: Only one asset allowed');

        // Throws unless that allocation has exactly n outcomes
        Outcome.Allocation[] memory allocations = outcome[0].allocations;

        require(allocations.length == nParticipants, '|Allocation|!=|participants|');

        for (uint256 i = 0; i < nParticipants; i++) {
            require(
                allocations[i].allocationType == uint8(Outcome.AllocationType.simple),
                'not a simple allocation'
            );
        }
    }

    /**
     * @notice Require specific rules in variable parts are followed when progressing state.
     * @dev Require specific rules in variable parts are followed when progressing state.
     * @param nParticipants Number of participants in a channel.
     * @param a Variable part to progress from.
     * @param b Variable part to progress to.
     */
    function _requireValidTransition(
        uint256 nParticipants,
        VariablePart memory a,
        VariablePart memory b
    ) internal pure {
        Outcome.Allocation[] memory allocationsA = a.outcome[0].allocations;
        Outcome.Allocation[] memory allocationsB = b.outcome[0].allocations;

        // Interprets the nth outcome as benefiting participant n
        // checks the destinations have not changed
        // Checks that the sum of assets hasn't changed
        // And that for all non-movers
        // the balance hasn't decreased
        uint256 allocationSumA;
        uint256 allocationSumB;
        for (uint256 i = 0; i < nParticipants; i++) {
            require(
                allocationsB[i].destination == allocationsA[i].destination,
                'Destinations may not change'
            );
            allocationSumA += allocationsA[i].amount;
            allocationSumB += allocationsB[i].amount;
            if (i != b.turnNum % nParticipants) {
                require(
                    allocationsB[i].amount >= allocationsA[i].amount,
                    'Nonmover balance decreased'
                );
            }
        }
        require(allocationSumA == allocationSumB, 'Total allocated cannot change');
    }
}
