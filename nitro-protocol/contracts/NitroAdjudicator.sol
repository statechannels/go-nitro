// SPDX-License-Identifier: MIT
pragma solidity 0.7.4;
pragma experimental ABIEncoderV2;

import './ForceMove.sol';
import {ExitFormat as Outcome} from '@statechannels/exit-format/contracts/ExitFormat.sol';
import './MultiAssetHolder.sol';

/**
 * @dev The NitroAdjudicator contract extends MultiAssetHolder and ForceMove
 */
contract NitroAdjudicator is ForceMove, MultiAssetHolder {
    /**
     * @notice Finalizes a channel by providing a finalization proof, and liquidates all assets for the channel.
     * @dev Finalizes a channel by providing a finalization proof, and liquidates all assets for the channel.
     * @param fixedPart Data describing properties of the state channel that do not change with state updates.
     * @param latestVariablePart Latest variable part in finalization proof. Must have the largest turnNum and the same appData and outcome as all other variable parts in finalization proof.
     * @param numStates The number of states in the finalization proof.
     * @param whoSignedWhat An array denoting which participant has signed which state: `participant[i]` signed the state with index `whoSignedWhat[i]`.
     * @param sigs Array of signatures, one for each participant, in participant order (e.g. [sig of participant[0], sig of participant[1], ...]).
     */
    function concludeAndTransferAllAssets(
        FixedPart memory fixedPart,
        IForceMoveApp.VariablePart memory latestVariablePart,
        uint8 numStates,
        uint8[] memory whoSignedWhat,
        Signature[] memory sigs
    ) public {
        bytes32 channelId = _conclude(
            fixedPart,
            latestVariablePart,
            numStates,
            whoSignedWhat,
            sigs
        );

        transferAllAssets(channelId, latestVariablePart.outcome, bytes32(0));
    }

    /**
     * @notice Liquidates all assets for the channel
     * @dev Liquidates all assets for the channel
     * @param channelId Unique identifier for a state channel
     * @param outcome An array of SingleAssetExit[] items.
     * @param stateHash stored state hash for the channel
     */
    function transferAllAssets(
        bytes32 channelId,
        Outcome.SingleAssetExit[] memory outcome,
        bytes32 stateHash
    ) public {
        // checks
        _requireChannelFinalized(channelId);
        _requireMatchingFingerprint(stateHash, _hashOutcome(outcome), channelId);

        // computation
        bool allocatesOnlyZerosForAllAssets = true;
        Outcome.SingleAssetExit[] memory exit = new Outcome.SingleAssetExit[](outcome.length);
        uint256[] memory initialHoldings = new uint256[](outcome.length);
        uint256[] memory totalPayouts = new uint256[](outcome.length);
        for (uint256 assetIndex = 0; assetIndex < outcome.length; assetIndex++) {
            Outcome.SingleAssetExit memory assetOutcome = outcome[assetIndex];
            Outcome.Allocation[] memory allocations = assetOutcome.allocations;
            address asset = outcome[assetIndex].asset;
            initialHoldings[assetIndex] = holdings[asset][channelId];
            (
                Outcome.Allocation[] memory newAllocations,
                bool allocatesOnlyZeros,
                Outcome.Allocation[] memory exitAllocations,
                uint256 totalPayoutsForAsset
            ) = compute_transfer_effects_and_interactions(
                    initialHoldings[assetIndex],
                    allocations,
                    new uint256[](0)
                );
            if (!allocatesOnlyZeros) allocatesOnlyZerosForAllAssets = false;
            totalPayouts[assetIndex] = totalPayoutsForAsset;
            outcome[assetIndex].allocations = newAllocations;
            exit[assetIndex] = Outcome.SingleAssetExit(
                asset,
                assetOutcome.metadata,
                exitAllocations
            );
        }

        // effects
        for (uint256 assetIndex = 0; assetIndex < outcome.length; assetIndex++) {
            address asset = outcome[assetIndex].asset;
            holdings[asset][channelId] -= totalPayouts[assetIndex];
            emit AllocationUpdated(channelId, assetIndex, initialHoldings[assetIndex]);
        }

        if (allocatesOnlyZerosForAllAssets) {
            delete statusOf[channelId];
        } else {
            _updateFingerprint(channelId, stateHash, _hashOutcome(outcome));
        }

        // interactions
        _executeExit(exit);
    }

    /**
    * @notice Check that the submitted pair of states form a valid transition (public wrapper for internal function _requireValidTransition)
    * @dev Check that the submitted pair of states form a valid transition (public wrapper for internal function _requireValidTransition)
    * @param nParticipants Number of participants in the channel.
    transition
    * @param ab Variable parts of each of the pair of states
    * @param appDefinition Address of deployed contract containing application-specific validTransition function.
    * @return true if the later state is a validTransition from its predecessor, reverts otherwise.
    */
    function validTransition(
        uint256 nParticipants,
        IForceMoveApp.VariablePart[2] memory ab, // [a,b]
        address appDefinition
    ) public pure returns (bool) {
        return _requireValidTransition(nParticipants, ab, appDefinition);
    }

    /**
     * @notice Executes an exit by paying out assets and calling external contracts
     * @dev Executes an exit by paying out assets and calling external contracts
     * @param exit The exit to be paid out.
     */
    function _executeExit(Outcome.SingleAssetExit[] memory exit) internal {
        for (uint256 assetIndex = 0; assetIndex < exit.length; assetIndex++) {
            _executeSingleAssetExit(exit[assetIndex]);
        }
    }
}
