// SPDX-License-Identifier: MIT
pragma solidity 0.7.6;
pragma experimental ABIEncoderV2;

import {ExitFormat as Outcome} from '@statechannels/exit-format/contracts/ExitFormat.sol';
import '../../interfaces/INitroTypes.sol';

library TurnTaking {
    /**
     * @notice Require supplied arguments to comply with turn taking logic, i.e. each participant signed the one state, they were mover for.
     * @dev Require supplied arguments to comply with turn taking logic, i.e. each participant signed the one state, they were mover for.
     * @param fixedPart Data describing properties of the state channel that do not change with state updates.
     * @param signedVariableParts An ordered array of structs, each struct describing dynamic properties of the state channel and must be signed by corresponding moving participant.
     */
    function requireValidTurnTaking(
        INitroTypes.FixedPart memory fixedPart,
        INitroTypes.SignedVariablePart[] memory signedVariableParts
    ) internal pure {
        require(fixedPart.participants.length == signedVariableParts.length, 'Invalid amount of variable parts');
        
        uint48 turnNum = signedVariableParts[0].variablePart.turnNum;

        for (uint i = 0; i < signedVariableParts.length; i++) {
            requireSignedByMover(fixedPart, signedVariableParts[i]);
            requireHasTurnNum(signedVariableParts[i].variablePart, turnNum);
            turnNum++;
        }
    }

    /**
     * @notice Require supplied state is signed by its corresponding moving participant.
     * @dev Require supplied state is signed by its corresponding moving participant.
     * @param fixedPart Data describing properties of the state channel that do not change with state updates.
     * @param signedVariablePart A struct describing dynamic properties of the state channel, that must be signed by moving participant.
     */
    function requireSignedByMover(
        INitroTypes.FixedPart memory fixedPart,
        INitroTypes.SignedVariablePart memory signedVariablePart
    ) internal pure {
        require(signedVariablePart.sigs.length == 1, 'sigs.length != 1');
        _requireSignedBy(
            _hashState(
                _getChannelId(fixedPart),
                signedVariablePart.variablePart.appData,
                signedVariablePart.variablePart.outcome,
                signedVariablePart.variablePart.turnNum,
                signedVariablePart.variablePart.isFinal
            ),
            signedVariablePart.sigs[0],
            _moverAddress(fixedPart.participants, signedVariablePart.variablePart.turnNum)
        );
    }

    /**
     * @notice Find moving participant address based on state turn number.
     * @dev Find moving participant address based on state turn number.
     * @param participants Array of participant addresses.
     * @param turnNum State turn number.
     * @return address Moving partitipant address.
     */
    function _moverAddress(
        address[] memory participants,
        uint48 turnNum
    ) internal pure returns (address) {
        return participants[turnNum % participants.length];
    }

    /**
     * @notice Require supplied stateHash is signed by signer.
     * @dev Require supplied stateHash is signed by signer.
     * @param stateHash State hash to check.
     * @param sig Signed state signature.
     * @param signer Address which must have signed the state.
     */
    function _requireSignedBy(
        bytes32 stateHash,
        INitroTypes.Signature memory sig,
        address signer
    ) internal pure {
        address recovered = _recoverSigner(stateHash, sig);
        require(signer == recovered, 'Invalid signer');
    }

    /**
     * @notice Require supplied variable part has specified turn number.    
     * @dev Require supplied variable part has specified turn number.
     * @param variablePart Variable part to check turn number of.
     * @param turnNum Turn number to compare with.
     */
    function requireHasTurnNum(
        INitroTypes.VariablePart memory variablePart,
        uint48 turnNum
    ) internal pure {
        require(
            variablePart.turnNum == turnNum,
            'Wrong variablePart.turnNum'
        );
    }

    /**
     * @notice Given a digest and ethereum digital signature, recover the signer
     * @dev Given a digest and digital signature, recover the signer
     * @param _d message digest
     * @param sig ethereum digital signature
     * @return signer
     */
    function _recoverSigner(bytes32 _d, INitroTypes.Signature memory sig) internal pure returns (address) {
        bytes32 prefixedHash = keccak256(abi.encodePacked('\x19Ethereum Signed Message:\n32', _d));
        address a = ecrecover(prefixedHash, sig.v, sig.r, sig.s);
        require(a != address(0), 'Invalid signature');
        return (a);
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
     * @notice Computes the unique id of a channel.
     * @dev Computes the unique id of a channel.
     * @param fixedPart Part of the state that does not change
     * @return channelId
     */
    function _getChannelId(INitroTypes.FixedPart memory fixedPart) internal pure returns (bytes32 channelId) {
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
}
