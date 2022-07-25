// SPDX-License-Identifier: MIT
pragma solidity 0.7.6;
pragma experimental ABIEncoderV2;

import './INitroTypes.sol';

/**
 * @dev The IForceMove interface defines the interface that an implementation of ForceMove should implement. ForceMove protocol allows state channels to be adjudicated and finalized.
 */
interface IForceMove is INitroTypes {
    /**
     * @notice Registers a challenge against a state channel. A challenge will either prompt another participant into clearing the challenge (via one of the other methods), or cause the channel to finalize at a specific time.
     * @dev Registers a challenge against a state channel. A challenge will either prompt another participant into clearing the challenge (via one of the other methods), or cause the channel to finalize at a specific time.
     * @param fixedPart Data describing properties of the state channel that do not change with state updates.
     * @param proof An ordered array of structs, that can be signed by any number of participants, each struct describing the properties of the state channel that may change with each state update. The proof is a validation for the supplied candidate.
     * @param candidate A struct, that can be signed by any number of participants, describing the properties of the state channel to change to. The candidate state is supported by proof states.
     * @param challengerSig The signature of a participant on the keccak256 of the abi.encode of (supportedStateHash, 'forceMove').
     */
    function challenge(
        FixedPart memory fixedPart,
        SignedVariablePart[] memory proof,
        SignedVariablePart memory candidate,
        Signature memory challengerSig
    ) external;

    /**
     * @notice Overwrites the `turnNumRecord` stored against a channel by providing a proof with higher turn number.
     * @dev Overwrites the `turnNumRecord` stored against a channel by providing a proof with higher turn number.
     * @param fixedPart Data describing properties of the state channel that do not change with state updates.
     * @param proof An ordered array of structs, that can be signed by any number of participants, each struct describing the properties of the state channel that may change with each state update. The proof is a validation for the supplied candidate.
     * @param candidate A struct, that can be signed by any number of participants, describing the properties of the state channel to change to. The candidate state is supported by proof states.
     */
    function checkpoint(
        FixedPart memory fixedPart,
        SignedVariablePart[] memory proof,
        SignedVariablePart memory candidate
    ) external;

    /**
     * @notice Finalizes a channel by providing a finalization proof. External wrapper for _conclude.
     * @dev Finalizes a channel by providing a finalization proof. External wrapper for _conclude.
     * @param fixedPart Data describing properties of the state channel that do not change with state updates.
     * @param proof An ordered array of structs, that can be signed by any number of participants, each struct describing the properties of the state channel that may change with each state update. The proof is a validation for the supplied candidate.
     * @param candidate A struct, that can be signed by any number of participants, describing the properties of the state channel to change to. The candidate state is supported by proof states.
     */
    function conclude(
        FixedPart memory fixedPart,
        SignedVariablePart[] memory proof,
        SignedVariablePart memory candidate
    ) external;

    // events

    /**
     * @dev Indicates that a challenge has been registered against `channelId`.
     * @param channelId Unique identifier for a state channel.
     * @param turnNumRecord A turnNum that (the adjudicator knows) is supported by a signature from each participant.
     * @param finalizesAt The unix timestamp when `channelId` will finalize.
     * @param isFinal Boolean denoting whether the challenge state is final.
     * @param fixedPart Data describing properties of the state channel that do not change with state updates.
     * @param proof An ordered array of structs, that can be signed by any number of participants, each struct describing the properties of the state channel that may change with each state update. The proof is a validation for the supplied candidate.
     * @param candidate A struct, that can be signed by any number of participants, describing the properties of the state channel to change to. The candidate state is supported by proof states.
     */
    event ChallengeRegistered(
        bytes32 indexed channelId,
        uint48 turnNumRecord,
        uint48 finalizesAt,
        bool isFinal,
        FixedPart fixedPart,
        SignedVariablePart[] memory proof,
        SignedVariablePart memory candidate
    );

    /**
     * @dev Indicates that a challenge, previously registered against `channelId`, has been cleared.
     * @param channelId Unique identifier for a state channel.
     * @param newTurnNumRecord A turnNum that (the adjudicator knows) is supported by a signature from each participant.
     */
    event ChallengeCleared(bytes32 indexed channelId, uint48 newTurnNumRecord);

    /**
     * @dev Indicates that a challenge has been registered against `channelId`.
     * @param channelId Unique identifier for a state channel.
     * @param finalizesAt The unix timestamp when `channelId` finalized.
     */
    event Concluded(bytes32 indexed channelId, uint48 finalizesAt);
}
