// SPDX-License-Identifier: MIT
pragma solidity 0.7.6;
pragma experimental ABIEncoderV2;

import '../interfaces/IForceMoveApp.sol';
import {ExitFormat as Outcome} from '@statechannels/exit-format/contracts/ExitFormat.sol';

/**
 * @dev The SingleAssetPayments contract complies with the ForceMoveApp interface and implements a simple payment channel with a single asset type only.
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
    ) external pure override returns (uint256) {
        // TODO see https://github.com/statechannels/go-nitro/issues/558
        // // Throws if more than one asset
        // require(a.outcome.length == 1, 'outcomeA: Only one asset allowed');
        // require(b.outcome.length == 1, 'outcomeB: Only one asset allowed');

        // Outcome.SingleAssetExit memory assetOutcomeA = a.outcome[0];
        // Outcome.SingleAssetExit memory assetOutcomeB = b.outcome[0];

        // // Throws unless that allocation has exactly n outcomes
        // Outcome.Allocation[] memory allocationsA = assetOutcomeA.allocations;
        // Outcome.Allocation[] memory allocationsB = assetOutcomeB.allocations;

        // require(allocationsA.length == nParticipants, '|AllocationA|!=|participants|');
        // require(allocationsB.length == nParticipants, '|AllocationB|!=|participants|');

        // for (uint256 i = 0; i < nParticipants; i++) {
        //     require(
        //         allocationsA[i].allocationType == uint8(Outcome.AllocationType.simple),
        //         'not a simple allocation'
        //     );
        //     require(
        //         allocationsB[i].allocationType == uint8(Outcome.AllocationType.simple),
        //         'not a simple allocation'
        //     );
        // }

        // // Interprets the nth outcome as benefiting participant n
        // // checks the destinations have not changed
        // // Checks that the sum of assets hasn't changed
        // // And that for all non-movers
        // // the balance hasn't decreased
        // uint256 allocationSumA;
        // uint256 allocationSumB;
        // for (uint256 i = 0; i < nParticipants; i++) {
        //     require(
        //         allocationsB[i].destination == allocationsA[i].destination,
        //         'Destinations may not change'
        //     );
        //     allocationSumA += allocationsA[i].amount;
        //     allocationSumB += allocationsB[i].amount;
        //     if (i != b.turnNum % nParticipants) {
        //         require(
        //             allocationsB[i].amount >= allocationsA[i].amount,
        //             'Nonmover balance decreased'
        //         );
        //     }
        // }
        // require(allocationSumA == allocationSumB, 'Total allocated cannot change');

        // return true;
    }
}
