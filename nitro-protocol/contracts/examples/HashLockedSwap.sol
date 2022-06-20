// SPDX-License-Identifier: MIT
pragma solidity 0.7.6;
pragma experimental ABIEncoderV2;

import {ExitFormat as Outcome} from '@statechannels/exit-format/contracts/ExitFormat.sol';
import {StrictTurnTaking} from '../libraries/signature-logic/StrictTurnTaking.sol';
import '../interfaces/IForceMoveApp.sol';

/**
 * @dev The HashLockedSwap contract complies with the ForceMoveApp interfaces and strict turn taking logic and implements a HashLockedSwaped payment.
 */
contract HashLockedSwap is IForceMoveApp {
    struct AppData {
        bytes32 h;
        bytes preImage;
    }

    /**
     * @notice Decodes the appData.
     * @dev Decodes the appData.
     * @param appDataBytes The abi.encode of a AppData struct describing the application-specific data.
     * @return AppData struct containing the application-specific data.
     */
    function appData(bytes memory appDataBytes) internal pure returns (AppData memory) {
        bytes memory decodedAppData = abi.decode(appDataBytes, (bytes));
        return abi.decode(decodedAppData, (AppData));
    }

    function latestSupportedState(
        FixedPart calldata fixedPart,
        SignedVariablePart[] calldata signedVariableParts
    ) external pure override returns (VariablePart memory) {
        VariablePart memory from = signedVariableParts[0].variablePart;
        VariablePart memory to = signedVariableParts[1].variablePart;

        // is this the first and only swap?
        require(signedVariableParts.length == 2, 'signedVariableParts.length != 2');
        require(to.turnNum == 4, 'latest turn number != 4');

        StrictTurnTaking.requireValidTurnTaking(fixedPart, signedVariableParts);

        // Decode variables.
        // Assumptions:
        //  - single asset in this channel
        //  - two parties in this channel
        Outcome.Allocation[] memory allocationsA = decode2PartyAllocation(from.outcome);
        Outcome.Allocation[] memory allocationsB = decode2PartyAllocation(to.outcome);
        bytes32 h = appData(from.appData).h;
        bytes memory preImage = appData(to.appData).preImage;

        // is the preimage correct?
        require(sha256(preImage) == h, 'incorrect preimage');
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

        // Throws unless there are exactly 2 allocations
        require(allocations.length == 2, 'allocation.length != 2');
    }
}
