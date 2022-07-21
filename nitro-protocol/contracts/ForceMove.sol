// SPDX-License-Identifier: MIT
pragma solidity 0.7.6;
pragma experimental ABIEncoderV2;

import {ExitFormat as Outcome} from '@statechannels/exit-format/contracts/ExitFormat.sol';
import {NitroUtils} from './libraries/NitroUtils.sol';
import './interfaces/IForceMove.sol';
import './interfaces/IForceMoveApp.sol';
import './StatusManager.sol';

/**
 * @dev An implementation of ForceMove protocol, which allows state channels to be adjudicated and finalized.
 */
contract ForceMove is IForceMove, StatusManager {
    // *****************
    // External methods:
    // *****************

    /**
     * @notice Unpacks turnNumRecord, finalizesAt and fingerprint from the status of a particular channel.
     * @dev Unpacks turnNumRecord, finalizesAt and fingerprint from the status of a particular channel.
     * @param channelId Unique identifier for a state channel.
     * @return turnNumRecord A turnNum that (the adjudicator knows) is supported by a signature from each participant.
     * @return finalizesAt The unix timestamp when `channelId` will finalize.
     * @return fingerprint The last 160 bits of kecca256(stateHash, outcomeHash)
     */
    function unpackStatus(bytes32 channelId)
        external
        view
        returns (
            uint48 turnNumRecord,
            uint48 finalizesAt,
            uint160 fingerprint
        )
    {
        (turnNumRecord, finalizesAt, fingerprint) = _unpackStatus(channelId);
    }

    /**
     * @notice Registers a challenge against a state channel. A challenge will either prompt another participant into clearing the challenge (via one of the other methods), or cause the channel to finalize at a specific time.
     * @dev Registers a challenge against a state channel. A challenge will either prompt another participant into clearing the challenge (via one of the other methods), or cause the channel to finalize at a specific time.
     * @param fixedPart Data describing properties of the state channel that do not change with state updates.
     * @param signedVariableParts An ordered array of structs, that can be signed by any number of participants, each struct describing the properties of the state channel that may change with each state update.
     * @param challengerSig The signature of a participant on the keccak256 of the abi.encode of (supportedStateHash, 'forceMove').
     */
    function challenge(
        FixedPart memory fixedPart,
        SignedVariablePart[] memory signedVariableParts,
        Signature memory challengerSig
    ) external override {
        bytes32 channelId = NitroUtils.getChannelId(fixedPart);
        uint48 largestTurnNum = _lastVariablePart(signedVariableParts).turnNum;

        if (_mode(channelId) == ChannelMode.Open) {
            _requireNonDecreasedTurnNumber(channelId, largestTurnNum);
        } else if (_mode(channelId) == ChannelMode.Challenge) {
            _requireIncreasedTurnNumber(channelId, largestTurnNum);
        } else {
            // This should revert.
            _requireChannelNotFinalized(channelId);
        }
        bytes32 supportedStateHash = _requireStateSupportedBy(
            fixedPart,
            signedVariableParts,
            channelId
        );

        _requireChallengerIsParticipant(supportedStateHash, fixedPart.participants, challengerSig);

        // effects

        emit ChallengeRegistered(
            channelId,
            largestTurnNum,
            uint48(block.timestamp) + fixedPart.challengeDuration, //solhint-disable-line not-rely-on-time
            // This could overflow, so don't join a channel with a huge challengeDuration
            _lastVariablePart(signedVariableParts).isFinal,
            fixedPart,
            signedVariableParts
        );

        statusOf[channelId] = _generateStatus(
            ChannelData(
                largestTurnNum,
                uint48(block.timestamp) + fixedPart.challengeDuration, //solhint-disable-line not-rely-on-time
                supportedStateHash,
                NitroUtils.hashOutcome(_lastVariablePart(signedVariableParts).outcome)
            )
        );
    }

    /**
     * @notice Overwrites the `turnNumRecord` stored against a channel by providing a state with higher turn number, supported by a signature from each participant.
     * @dev Overwrites the `turnNumRecord` stored against a channel by providing a state with higher turn number, supported by a signature from each participant.
     * @param fixedPart Data describing properties of the state channel that do not change with state updates.
     * @param signedVariableParts An ordered array of structs, that can be signed by any number of participants, each struct describing the properties of the state channel that may change with each state update.
     */
    function checkpoint(FixedPart memory fixedPart, SignedVariablePart[] memory signedVariableParts)
        external
        override
    {
        bytes32 channelId = NitroUtils.getChannelId(fixedPart);
        uint48 largestTurnNum = _lastVariablePart(signedVariableParts).turnNum;

        // checks
        _requireChannelNotFinalized(channelId);
        _requireIncreasedTurnNumber(channelId, largestTurnNum);
        _requireStateSupportedBy(fixedPart, signedVariableParts, channelId);

        // effects
        _clearChallenge(channelId, largestTurnNum);
    }

    /**
     * @notice Finalizes a channel by providing a finalization proof. External wrapper for _conclude.
     * @dev Finalizes a channel by providing a finalization proof. External wrapper for _conclude.
     * @param fixedPart Data describing properties of the state channel that do not change with state updates.
     * @param signedVariableParts An array of signed variable parts. All variable parts have to be marked `final`.
     */
    function conclude(FixedPart memory fixedPart, SignedVariablePart[] memory signedVariableParts)
        external
        override
    {
        _conclude(fixedPart, signedVariableParts);
    }

    /**
     * @notice Finalizes a channel by providing a finalization proof. Internal method.
     * @dev Finalizes a channel by providing a finalization proof. Internal method.
     * @param fixedPart Data describing properties of the state channel that do not change with state updates.
     * @param signedVariableParts An array of signed variable parts. All variable parts have to be marked `final`.
     */
    function _conclude(FixedPart memory fixedPart, SignedVariablePart[] memory signedVariableParts)
        internal
        returns (bytes32 channelId)
    {
        channelId = NitroUtils.getChannelId(fixedPart);
        _requireChannelNotFinalized(channelId);

        // checks
        _requireStateSupportedBy(fixedPart, signedVariableParts, channelId);

        // effects
        statusOf[channelId] = _generateStatus(
            ChannelData(
                0,
                uint48(block.timestamp), //solhint-disable-line not-rely-on-time
                bytes32(0),
                NitroUtils.hashOutcome(_lastVariablePart(signedVariableParts).outcome)
            )
        );
        emit Concluded(channelId, uint48(block.timestamp)); //solhint-disable-line not-rely-on-time
    }

    function getChainID() public pure returns (uint256) {
        return NitroUtils.getChainID();
    }

    // *****************
    // Internal methods:
    // *****************

    /**
     * @notice Checks that the challengerSignature was created by one of the supplied participants.
     * @dev Checks that the challengerSignature was created by one of the supplied participants.
     * @param supportedStateHash Forms part of the digest to be signed, along with the string 'forceMove'.
     * @param participants A list of addresses representing the participants of a channel.
     * @param challengerSignature The signature of a participant on the keccak256 of the abi.encode of (supportedStateHash, 'forceMove').
     */
    function _requireChallengerIsParticipant(
        bytes32 supportedStateHash,
        address[] memory participants,
        Signature memory challengerSignature
    ) internal pure {
        address challenger = NitroUtils.recoverSigner(
            keccak256(abi.encode(supportedStateHash, 'forceMove')),
            challengerSignature
        );
        require(_isAddressInArray(challenger, participants), 'Challenger is not a participant');
    }

    /**
     * @notice Tests whether a given address is in a given array of addresses.
     * @dev Tests whether a given address is in a given array of addresses.
     * @param suspect A single address of interest.
     * @param addresses A line-up of possible perpetrators.
     * @return true if the address is in the array, false otherwise
     */
    function _isAddressInArray(address suspect, address[] memory addresses)
        internal
        pure
        returns (bool)
    {
        for (uint256 i = 0; i < addresses.length; i++) {
            if (suspect == addresses[i]) {
                return true;
            }
        }
        return false;
    }

    /**
     * @notice Check that the submitted data constitute a support proof.
     * @dev Check that the submitted data constitute a support proof.
     * @param fixedPart Fixed Part of the states in the support proof.
     * @param signedVariableParts Signed variable parts of the states in the support proof.
     * @param channelId Unique identifier for a channel.
     * @return The hash of the latest state in the proof, if supported, else reverts.
     */
    function _requireStateSupportedBy(
        FixedPart memory fixedPart,
        SignedVariablePart[] memory signedVariableParts,
        bytes32 channelId
    ) internal pure returns (bytes32) {
        VariablePart memory latestVariablePart = IForceMoveApp(fixedPart.appDefinition)
            .latestSupportedState(fixedPart, recoverVariableParts(fixedPart, signedVariableParts));

        // enforcing the latest supported state being in the last slot of the array
        _requireVariablePartIsLast(latestVariablePart, signedVariableParts);

        return
            NitroUtils.hashState(
                channelId,
                latestVariablePart.appData,
                latestVariablePart.outcome,
                latestVariablePart.turnNum,
                latestVariablePart.isFinal
            );
    }

    /**
     * @notice Recover signatures for each variable part in the supplied array.
     * @dev Recover signatures for each variable part in the supplied array.
     * @param fixedPart Fixed Part of the states in the support proof.
     * @param signedVariableParts Signed variable parts of the states in the support proof.
     * @return An array of recoveredVariableParts, identical to the supplied signedVariableParts array but with the signatures replaced with a signedBy bitmask.
     */
    function recoverVariableParts(
        FixedPart memory fixedPart,
        SignedVariablePart[] memory signedVariableParts
    ) internal pure returns (RecoveredVariablePart[] memory) {
        RecoveredVariablePart[] memory recoveredVariableParts = new RecoveredVariablePart[](
            signedVariableParts.length
        );
        for (uint256 i = 0; i < signedVariableParts.length; i++) {
            recoveredVariableParts[i] = recoverVariablePart(fixedPart, signedVariableParts[i]);
        }
        return recoveredVariableParts;
    }

    /**
     * @notice Recover signatures for a variable part.
     * @dev Recover signatures for a variable part.
     * @param fixedPart Fixed Part of the states in the support proof.
     * @param signedVariablePart A signed variable part.
     * @return An RecoveredVariableParts, identical to the supplied signedVariablePart  but with the signatures replaced with a signedBy bitmask.
     */
    function recoverVariablePart(
        FixedPart memory fixedPart,
        SignedVariablePart memory signedVariablePart
    ) internal pure returns (RecoveredVariablePart memory) {
        RecoveredVariablePart memory rvp = RecoveredVariablePart({
            variablePart: signedVariablePart.variablePart,
            signedBy: 0
        });
        //  For each signature
        for (uint256 j = 0; j < signedVariablePart.sigs.length; j++) {
            address signer = NitroUtils.recoverSigner(
                NitroUtils.hashState(fixedPart, signedVariablePart.variablePart),
                signedVariablePart.sigs[j]
            );
            // Check each participant to see if they signed it
            for (uint256 i = 0; i < fixedPart.participants.length; i++) {
                if (signer == fixedPart.participants[i]) {
                    rvp.signedBy += 2**i;
                    break; // Once we have found a match, assuming distinct participants, no-one else signed it
                }
            }
        }
        return rvp;
    }

    /**
     * @notice Check whether supplied variablePart is in the last slot if variableParts.
     * @dev Check whether supplied variablePart is in the last slot if variableParts.
     * @param variablePart VariablePart to be in the last slot.
     * @param signedVariableParts SignedVariableParts the last slot of to check.
     */
    function _requireVariablePartIsLast(
        VariablePart memory variablePart,
        SignedVariablePart[] memory signedVariableParts
    ) internal pure {
        require(
            NitroUtils.bytesEqual(
                abi.encode(signedVariableParts[signedVariableParts.length - 1].variablePart),
                abi.encode(variablePart)
            ),
            'variablePart not the last.'
        );
    }

    /**
     * @notice Clears a challenge by updating the turnNumRecord and resetting the remaining channel storage fields, and emits a ChallengeCleared event.
     * @dev Clears a challenge by updating the turnNumRecord and resetting the remaining channel storage fields, and emits a ChallengeCleared event.
     * @param channelId Unique identifier for a channel.
     * @param newTurnNumRecord New turnNumRecord to overwrite existing value
     */
    function _clearChallenge(bytes32 channelId, uint48 newTurnNumRecord) internal {
        statusOf[channelId] = _generateStatus(
            ChannelData(newTurnNumRecord, 0, bytes32(0), bytes32(0))
        );
        emit ChallengeCleared(channelId, newTurnNumRecord);
    }

    /**
     * @notice Checks that the submitted turnNumRecord is strictly greater than the turnNumRecord stored on chain.
     * @dev Checks that the submitted turnNumRecord is strictly greater than the turnNumRecord stored on chain.
     * @param channelId Unique identifier for a channel.
     * @param newTurnNumRecord New turnNumRecord intended to overwrite existing value
     */
    function _requireIncreasedTurnNumber(bytes32 channelId, uint48 newTurnNumRecord) internal view {
        (uint48 turnNumRecord, , ) = _unpackStatus(channelId);
        require(newTurnNumRecord > turnNumRecord, 'turnNumRecord not increased.');
    }

    /**
     * @notice Checks that the submitted turnNumRecord is greater than or equal to the turnNumRecord stored on chain.
     * @dev Checks that the submitted turnNumRecord is greater than or equal to the turnNumRecord stored on chain.
     * @param channelId Unique identifier for a channel.
     * @param newTurnNumRecord New turnNumRecord intended to overwrite existing value
     */
    function _requireNonDecreasedTurnNumber(bytes32 channelId, uint48 newTurnNumRecord)
        internal
        view
    {
        (uint48 turnNumRecord, , ) = _unpackStatus(channelId);
        require(newTurnNumRecord >= turnNumRecord, 'turnNumRecord decreased.');
    }

    /**
     * @notice Checks that a given channel is NOT in the Finalized mode.
     * @dev Checks that a given channel is in the Challenge mode.
     * @param channelId Unique identifier for a channel.
     */
    function _requireChannelNotFinalized(bytes32 channelId) internal view {
        require(_mode(channelId) != ChannelMode.Finalized, 'Channel finalized.');
    }

    /**
     * @notice Checks that a given channel is in the Open mode.
     * @dev Checks that a given channel is in the Challenge mode.
     * @param channelId Unique identifier for a channel.
     */
    function _requireChannelOpen(bytes32 channelId) internal view {
        require(_mode(channelId) == ChannelMode.Open, 'Channel not open.');
    }

    /**
     * @notice Checks that a given ChannelData struct matches a supplied bytes32 when formatted for storage.
     * @dev Checks that a given ChannelData struct matches a supplied bytes32 when formatted for storage.
     * @param data A given ChannelData data structure.
     * @param s Some data in on-chain storage format.
     */
    function _matchesStatus(ChannelData memory data, bytes32 s) internal pure returns (bool) {
        return _generateStatus(data) == s;
    }

    /**
     * @notice Returns the last VariablePart from array of SignedVariableParts.
     * @dev Returns the last VariablePart from array of SignedVariableParts.
     * @param signedVariableParts Array of SignedVariableParts.
     * @return VariablePart Last VariablePart from array.
     */
    function _lastVariablePart(SignedVariablePart[] memory signedVariableParts)
        internal
        pure
        returns (VariablePart memory)
    {
        return signedVariableParts[signedVariableParts.length - 1].variablePart;
    }
}
