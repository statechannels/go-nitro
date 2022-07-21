// SPDX-License-Identifier: MIT
pragma solidity 0.7.6;
pragma experimental ABIEncoderV2;

import {ExitFormat as Outcome} from '@statechannels/exit-format/contracts/ExitFormat.sol';

interface INitroTypes {
    struct Signature {
        uint8 v;
        bytes32 r;
        bytes32 s;
    }

    struct VariablePart {
        Outcome.SingleAssetExit[] outcome;
        bytes appData;
        uint48 turnNum;
        bool isFinal;
    }

    struct SignedVariablePart {
        VariablePart variablePart;
        Signature[] sigs;
    }

    struct RecoveredVariablePart {
        VariablePart variablePart;
        uint256 signedBy; // bitmask
    }

    struct FixedPart {
        uint256 chainId;
        address[] participants;
        uint48 channelNonce;
        address appDefinition;
        uint48 challengeDuration;
    }

    struct State {
        // participants sign the hash of this
        bytes32 channelId; // keccack(FixedPart)
        bytes appData;
        bytes outcome;
        uint48 turnNum;
        bool isFinal;
    }
}
