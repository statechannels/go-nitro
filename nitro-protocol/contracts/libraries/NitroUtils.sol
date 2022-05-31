// SPDX-License-Identifier: MIT
pragma solidity 0.7.6;
pragma experimental ABIEncoderV2;

import '../interfaces/INitroTypes.sol';

library NitroUtils {
    // *****************
    // Signature methods:
    // *****************

    function isSignedBy(uint256 signedBy, uint8 participantIndex) internal pure returns (bool) {
        return ((signedBy >> participantIndex) % 2 == 1);
    }

    // This function can be used inside validTransition
    // To recover the legacy "turn taking" semantics
    function isRoundRobin(
        uint256 nParticipants,
        uint48 turnNumB,
        uint256 signedByFrom,
        uint256 signedByTo
    ) internal pure returns (bool) {
        require(
            turnNumB > 0 &&
                isSignedBy(signedByFrom, uint8((turnNumB - 1) % nParticipants)) &&
                isSignedBy(signedByTo, uint8(turnNumB % nParticipants)),
            'roundRobin violation'
        );
        return true;
    }

    /**
     * @notice Given a digest and ethereum digital signature, recover the signer
     * @dev Given a digest and digital signature, recover the signer
     * @param _d message digest
     * @param sig ethereum digital signature
     * @return signer
     */
    function recoverSigner(bytes32 _d, INitroTypes.Signature memory sig)
        internal
        pure
        returns (address)
    {
        bytes32 prefixedHash = keccak256(abi.encodePacked('\x19Ethereum Signed Message:\n32', _d));
        address a = ecrecover(prefixedHash, sig.v, sig.r, sig.s);
        require(a != address(0), 'Invalid signature');
        return (a);
    }

    // *****************
    // ID methods:
    // *****************

    /**
     * @notice Computes the unique id of a channel.
     * @dev Computes the unique id of a channel.
     * @param fixedPart Part of the state that does not change
     * @return channelId
     */
    function getChannelId(INitroTypes.FixedPart memory fixedPart) internal pure returns (bytes32 channelId) {
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

    function getChainID() public pure returns (uint256) {
        uint256 id;
        /* solhint-disable no-inline-assembly */
        assembly {
            id := chainid()
        }
        /* solhint-disable no-inline-assembly */
        return id;
    }

    // *****************
    // Hash methods:
    // *****************

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
    function hashState(
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
    function hashOutcome(Outcome.SingleAssetExit[] memory outcome)
        internal
        pure
        returns (bytes32)
    {
        return keccak256(Outcome.encodeExit(outcome));
    }

    // *****************
    // Equality methods:
    // *****************

    /**
     * @notice Check for equality of two byte strings
     * @dev Check for equality of two byte strings
     * @param _preBytes One bytes string
     * @param _postBytes The other bytes string
     * @return true if the bytes are identical, false otherwise.
     */
    function bytesEqual(bytes memory _preBytes, bytes memory _postBytes)
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
}
