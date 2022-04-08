// SPDX-License-Identifier: MIT
pragma solidity 0.7.4;
pragma experimental ABIEncoderV2;

import '../interfaces/IForceMoveApp.sol';
import {ExitFormat as Outcome} from '@statechannels/exit-format/contracts/ExitFormat.sol';

/**
 * @dev The SingleAssetPayments contract complies with the ForceMoveApp interface and implements a simple payment channel with a single asset type only.
 */
contract SingleAssetPayments is IForceMoveApp {
    /**
     * @notice Encodes the payment channel update rules.
     * @dev Encodes the payment channel update rules.
     * @param a State being transitioned from.
     * @param b State being transitioned to.
     * @param nParticipants Number of participants in this state channel.
     * @return true if the transition conforms to the rules, false otherwise.
     */
    function validTransition(
        VariablePart memory a,
        VariablePart memory b,
        uint256 nParticipants
    ) public override pure returns (bool) {
        // Throws if more than one asset
        require(a.outcome.length == 1, 'outcomeA: Only one asset allowed');
        require(b.outcome.length == 1, 'outcomeB: Only one asset allowed');

        Outcome.SingleAssetExit memory assetOutcomeA = a.outcome[0];
        Outcome.SingleAssetExit memory assetOutcomeB = b.outcome[0];

        // Throws unless that allocation has exactly n outcomes
        Outcome.Allocation[] memory allocationsA = assetOutcomeA.allocations;
        Outcome.Allocation[] memory allocationsB = assetOutcomeB.allocations;

        require(allocationsA.length == nParticipants, '|AllocationA|!=|participants|');
        require(allocationsB.length == nParticipants, '|AllocationB|!=|participants|');

        for (uint256 i = 0; i < nParticipants; i++) {
            require(
                allocationsA[i].allocationType == uint8(Outcome.AllocationType.simple),
                'not a simple allocation'
            );
            require(
                allocationsB[i].allocationType == uint8(Outcome.AllocationType.simple),
                'not a simple allocation'
            );
        }

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

        return true;
    }
}
