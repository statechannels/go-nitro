// SPDX-License-Identifier: MIT
pragma solidity 0.7.6;
pragma experimental ABIEncoderV2;

import '../interfaces/IForceMoveApp.sol';
import './examples/signature-logic/TurnTaking.sol';
import {ExitFormat as Outcome} from '@statechannels/exit-format/contracts/ExitFormat.sol';

/**
 * @dev The HashLockedSwap contract complies with the ForceMoveApp interface and implements a HashLockedSwaped payment
 */
contract HashLockedSwap is IForceMoveApp, TurnTaking {
    struct AppData {
        bytes32 h;
        bytes preImage;
    }

    function latestSupportedState(
        FixedPart calldata fixedPart,
        SignedVariablePart[] calldata signedVariableParts
    ) external pure override returns (VariablePart memory) {
        // is this the first and only swap?
        require(signedVariableParts.length == 2, 'excess states');
        require(signedVariableParts[1].variablePart.turnNum == 4, 'latest turn number != 4');

        _requireValidTurnTaking(fixedPart, signedVariableParts);

        // Decode variables.
        // Assumptions:
        //  - single asset in this channel
        //  - two parties in this channel
        //  - not a "guarantee" channel (c.f. Nitro paper)
        Outcome.Allocation[] memory allocationsA = decode2PartyAllocation(signedVariableParts[0].variablePart.outcome);
        Outcome.Allocation[] memory allocationsB = decode2PartyAllocation(signedVariableParts[1].variablePart.outcome);
        bytes memory preImage = abi.decode(signedVariableParts[1].variablePart.appData, (AppData)).preImage;
        bytes32 h = abi.decode(signedVariableParts[0].variablePart.appData, (AppData)).h;

        // is the preimage correct?
        require(sha256(preImage) == h, 'Incorrect preimage');
        // NOTE ON GAS COSTS
        // The gas cost of hashing depends on the choice of hash function
        // and the length of the the preImage.
        // sha256 is twice as expensive as keccak256
        // https://ethereum.stackexchange.com/a/3200
        // But is compatible with bitcoin.

        // slots for each participant unchanged
        require(
            allocationsA[0].destination == allocationsB[0].destination &&
                allocationsA[1].destination == allocationsB[1].destination,
            'destinations may not change'
        );

        // was the payment made?
        require(
            allocationsA[0].amount == allocationsB[1].amount &&
                allocationsA[1].amount == allocationsB[0].amount,
            'amounts must be permuted'
        );

        return signedVariableParts[1].variablePart;
    }

    function decode2PartyAllocation(Outcome.SingleAssetExit[] memory outcome)
        private
        pure
        returns (Outcome.Allocation[] memory allocations)
    {
        Outcome.SingleAssetExit memory assetOutcome = outcome[0];

        allocations = assetOutcome.allocations; // TODO should we check each allocation is a "simple" one?

        // Throws unless there are exactly 3 allocations
        require(allocations.length == 2, 'allocation.length != 3');
    }
}
