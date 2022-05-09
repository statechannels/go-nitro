// SPDX-License-Identifier: MIT
pragma solidity 0.7.6;
pragma experimental ABIEncoderV2;

import {ExitFormat as Outcome} from '@statechannels/exit-format/contracts/ExitFormat.sol';
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
     * @param signedVariableParts An ordered array of structs, that can be signed by any number of participants, each struct decribing the properties of the state channel that may change with each state update. Length is from 1 to the number of participants (inclusive).
     * @param challengerSig The signature of a participant on the keccak256 of the abi.encode of (supportedStateHash, 'forceMove').
     */
    function challenge(
        FixedPart memory fixedPart,
        SignedVariablePart[] memory signedVariableParts,
        Signature memory challengerSig
    ) external override {
        // input type validation
        requireValidInput(
            fixedPart.participants.length,
            signedVariableParts.length
        );

        bytes32 channelId = _getChannelId(fixedPart);
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
                _hashOutcome(_lastVariablePart(signedVariableParts).outcome)
            )
        );
    }

    /**
     * @notice Overwrites the `turnNumRecord` stored against a channel by providing a state with higher turn number, supported by a signature from each participant.
     * @dev Overwrites the `turnNumRecord` stored against a channel by providing a state with higher turn number, supported by a signature from each participant.
     * @param fixedPart Data describing properties of the state channel that do not change with state updates.
     * @param signedVariableParts An ordered array of structs, that can be signed by any number of participants, each struct decribing the properties of the state channel that may change with each state update. Length is from 1 to the number of participants (inclusive).
     */
    function checkpoint(
        FixedPart memory fixedPart,
        SignedVariablePart[] memory signedVariableParts
    ) external override {
        // input type validation
        requireValidInput(
            fixedPart.participants.length,
            signedVariableParts.length
        );

        bytes32 channelId = _getChannelId(fixedPart);
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
    function conclude(
        FixedPart memory fixedPart,
        SignedVariablePart[] memory signedVariableParts
    ) external override {
        _conclude(fixedPart, signedVariableParts);
    }

    /**
     * @notice Finalizes a channel by providing a finalization proof. Internal method.
     * @dev Finalizes a channel by providing a finalization proof. Internal method.
     * @param fixedPart Data describing properties of the state channel that do not change with state updates.
     * @param signedVariableParts An array of signed variable parts. All variable parts have to be marked `final`.
     */
    function _conclude(
        FixedPart memory fixedPart,
        SignedVariablePart[] memory signedVariableParts
    ) internal returns (bytes32 channelId) {
        channelId = _getChannelId(fixedPart);
        _requireChannelNotFinalized(channelId);

        // input type validation
        requireValidInput(
            fixedPart.participants.length,
            signedVariableParts.length
        );

        // checks
        _requireStateSupportedBy(fixedPart, signedVariableParts, channelId);

        // effects
        statusOf[channelId] = _generateStatus(
            ChannelData(
                0,
                uint48(block.timestamp), //solhint-disable-line not-rely-on-time
                bytes32(0),
                _hashOutcome(_lastVariablePart(signedVariableParts).outcome)
            )
        );
        emit Concluded(channelId, uint48(block.timestamp)); //solhint-disable-line not-rely-on-time
    }

    function getChainID() public pure returns (uint256) {
        uint256 id;
        /* solhint-disable no-inline-assembly */
        assembly {
            id := chainid()
        }
        /* solhint-disable no-inline-assembly */
        return id;
    }

    /**
     * @notice Validates related to participants and states input for several external methods.
     * @dev Validates related to participants and states input for several external methods.
     * @param numParticipants Length of the participants array.
     * @param numStates Number of states submitted.
     */
    function requireValidInput(
        uint256 numParticipants,
        uint256 numStates
    ) public pure returns (bool) {
        require((numParticipants >= numStates) && (numStates > 0), 'Insufficient or excess states');
        require(numParticipants <= type(uint8).max, 'Too many participants!'); // type(uint8).max = 2**8 - 1 = 255
        // no more than 255 participants
        // max index for participants is 254
        return true;
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
        address challenger = _recoverSigner(
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
     * @notice Given a digest and ethereum digital signature, recover the signer
     * @dev Given a digest and digital signature, recover the signer
     * @param _d message digest
     * @param sig ethereum digital signature
     * @return signer
     */
    function _recoverSigner(bytes32 _d, Signature memory sig) internal pure returns (address) {
        bytes32 prefixedHash = keccak256(abi.encodePacked('\x19Ethereum Signed Message:\n32', _d));
        address a = ecrecover(prefixedHash, sig.v, sig.r, sig.s);
        require(a != address(0), 'Invalid signature');
        return (a);
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
            .latestSupportedState(fixedPart, signedVariableParts);

        // enforcing the latest supported state being in the last slot of the array
        _requireVariablePartIsLast(latestVariablePart, signedVariableParts);

        return _hashState(
            channelId,
            latestVariablePart.appData,
            latestVariablePart.outcome,
            latestVariablePart.turnNum,
            latestVariablePart.isFinal
        );
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
            _bytesEqual(
                abi.encode(signedVariableParts[signedVariableParts.length - 1].variablePart),
                abi.encode(variablePart)
            ),
            'variablePart not the last.'
        );
    }

    /**
     * @notice Check for equality of two byte strings
     * @dev Check for equality of two byte strings
     * @param _preBytes One bytes string
     * @param _postBytes The other bytes string
     * @return true if the bytes are identical, false otherwise.
     */
    function _bytesEqual(bytes memory _preBytes, bytes memory _postBytes)
        internal
        pure
        returns (bool)
    {
        // copied from https://www.npmjs.com/package/solidity-bytes-utils/v/0.1.1
        bool success = true;

        /* solhint-disable no-inline-assembly */
        assembly {
            let length := mload(_preBytes)

            // if lengths don't match the arrays are not equal
            switch eq(length, mload(_postBytes))
            case 1 {
                // cb is a circuit breaker in the for loop since there's
                //  no said feature for inline assembly loops
                // cb = 1 - don't breaker
                // cb = 0 - break
                let cb := 1

                let mc := add(_preBytes, 0x20)
                let end := add(mc, length)

                for {
                    let cc := add(_postBytes, 0x20)
                    // the next line is the loop condition:
                    // while(uint256(mc < end) + cb == 2)
                } eq(add(lt(mc, end), cb), 2) {
                    mc := add(mc, 0x20)
                    cc := add(cc, 0x20)
                } {
                    // if any of these checks fails then arrays are not equal
                    if iszero(eq(mload(mc), mload(cc))) {
                        // unsuccess:
                        success := 0
                        cb := 0
                    }
                }
            }
            default {
                // unsuccess:
                success := 0
            }
        }
        /* solhint-disable no-inline-assembly */

        return success;
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
     * @notice Checks that a given ChannelData struct matches the challenge stored on chain, and that the channel is in Challenge mode.
     * @dev Checks that a given ChannelData struct matches the challenge stored on chain, and that the channel is in Challenge mode.
     * @param data A given ChannelData data structure.
     * @param channelId Unique identifier for a channel.
     */
    function _requireSpecificChallenge(ChannelData memory data, bytes32 channelId) internal view {
        _requireMatchingStorage(data, channelId);
        _requireOngoingChallenge(channelId);
    }

    /**
     * @notice Checks that a given channel is in the Challenge mode.
     * @dev Checks that a given channel is in the Challenge mode.
     * @param channelId Unique identifier for a channel.
     */
    function _requireOngoingChallenge(bytes32 channelId) internal view {
        require(_mode(channelId) == ChannelMode.Challenge, 'No ongoing challenge.');
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
     * @notice Checks that a given ChannelData struct matches the challenge stored on chain.
     * @dev Checks that a given ChannelData struct matches the challenge stored on chain.
     * @param data A given ChannelData data structure.
     * @param channelId Unique identifier for a channel.
     */
    function _requireMatchingStorage(ChannelData memory data, bytes32 channelId) internal view {
        require(_matchesStatus(data, statusOf[channelId]), 'status(ChannelData)!=storage');
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
     * @notice Computes the hash of the state corresponding to the input data.
     * @dev Computes the hash of the state corresponding to the input data.
     * @param turnNum Turn number
     * @param isFinal Is the state final?
     * @param channelId Unique identifier for the channel
     * @param appData Application specific data.
     * @param outcome Outcome structure.
     * @return The stateHash
     */
    function _hashState(
        bytes32 channelId,
        bytes memory appData,
        Outcome.SingleAssetExit[] memory outcome,
        uint48 turnNum,
        bool isFinal
    ) internal pure returns (bytes32) {
        return keccak256(abi.encode(channelId, appData, outcome, turnNum, isFinal));
    }

    /**
     * @notice Hashes the outcome structure. Internal helper.
     * @dev Hashes the outcome structure. Internal helper.
     * @param outcome Outcome structure to encode hash.
     * @return bytes32 Hash of encoded outcome structure.
     */
    function _hashOutcome(Outcome.SingleAssetExit[] memory outcome)
        internal
        pure
        returns (bytes32)
    {
        return keccak256(Outcome.encodeExit(outcome));
    }

    /**
     * @notice Computes the unique id of a channel.
     * @dev Computes the unique id of a channel.
     * @param fixedPart Part of the state that does not change
     * @return channelId
     */
    function _getChannelId(FixedPart memory fixedPart) internal pure returns (bytes32 channelId) {
        require(fixedPart.chainId == getChainID(), 'Incorrect chainId');
        channelId = keccak256(
            abi.encode(
                getChainID(),
                fixedPart.participants,
                fixedPart.channelNonce,
                fixedPart.appDefinition,
                fixedPart.challengeDuration
            )
        );
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
