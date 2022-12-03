// SPDX-License-Identifier: MIT
pragma solidity 0.8.17;
pragma experimental ABIEncoderV2;

import './interfaces/IForceMoveApp.sol';
import './libraries/NitroUtils.sol';
import './interfaces/INitroTypes.sol';
import {ExitFormat as Outcome} from '@statechannels/exit-format/contracts/ExitFormat.sol';

/**
 * @dev The VirtualMarginApp contract complies with the ForceMoveApp interface and allows payments to be made virtually from Alice to Bob (participants[0] to participants[n+1], where n is the number of intermediaries).
 */
contract VirtualMarginApp is IForceMoveApp {

    struct MarginState {
        uint256 leaderAmount;
        uint256 followerAmount;
        uint256 version;                    // Highest number is the valid state
        INitroTypes.Signature leaderSig;    // signature on abi.encode(channelId,amount)
        INitroTypes.Signature followerSig;  // signature on abi.encode(channelId,amount)
    }

    enum AllocationIndices {
        Leader,     // Leader is the virtual-channel initiator
        Follower    // Follower accepted to establish the link
    }

    /**
     * @notice Encodes application-specific rules for a particular ForceMove-compliant state channel.
     * @dev Encodes application-specific rules for a particular ForceMove-compliant state channel.
     * @param fixedPart Fixed part of the state channel.
     * @param proof Array of recovered variable parts which constitutes a support proof for the candidate.
     * @param candidate Recovered variable part the proof was supplied for.
     */
    function requireStateSupported(
        FixedPart calldata fixedPart,
        RecoveredVariablePart[] calldata proof,
        RecoveredVariablePart calldata candidate
    ) external pure override {
        // This channel has only 4 states which can be supported:
        // 0 prefund
        // 1 postfund
        // 2 redemption
        // 3 final

        // states 0,1,3 can be supported via unanimous consensus:

        if (proof.length == 0) {
            require(
                NitroUtils.getClaimedSignersNum(candidate.signedBy) ==
                    fixedPart.participants.length,
                '!unanimous; |proof|=0'
            );
            if (candidate.variablePart.turnNum == 0) return; // prefund
            if (candidate.variablePart.turnNum == 1) return; // postfund
            if (candidate.variablePart.turnNum == 3) {
                // final (note: there is a core protocol escape hatch for this, too, so it could be removed)
                require(candidate.variablePart.isFinal, '!final; turnNum=3 && |proof|=0');
                return;
            }
            revert('bad candidate turnNum; |proof|=0');
        }

        // State 2 can be supported via a forced transition from state 1:
        //
        //      (2)_B     [redemption state signed by Bob, includes a voucher signed by Alice. The outcome may be updated in favour of Bob]
        //      ^
        //      |
        //      (1)_AIB   [fully signed postfund]

        if (proof.length == 1) {
            requireProofOfUnanimousConsensusOnPostFund(proof[0], fixedPart.participants.length);
            require(candidate.variablePart.turnNum == 2, 'bad candidate turnNum; |proof|=1');
            require(
                NitroUtils.isClaimedSignedBy(candidate.signedBy, 2),
                'redemption not signed by Bob'
            );
            MarginState memory receipt = requireValidMarginState(candidate.variablePart.appData, fixedPart);
            requireCorrectAdjustments(
                proof[0].variablePart.outcome,
                candidate.variablePart.outcome,
                receipt
            );
            return;
        }
        revert('bad proof length');
    }

    function requireProofOfUnanimousConsensusOnPostFund(
        RecoveredVariablePart memory rVP,
        uint256 numParticipants
    ) internal pure {
        require(rVP.variablePart.turnNum == 1, 'bad proof[0].turnNum; |proof|=1');
        require(
            NitroUtils.getClaimedSignersNum(rVP.signedBy) == numParticipants,
            'postfund !unanimous; |proof|=1'
        );
    }

    function requireValidMarginState(bytes memory appData, FixedPart memory fixedPart)
        internal
        pure
        returns (MarginState memory)
    {
        MarginState memory receipt = abi.decode(appData, (MarginState));
        bytes32 receiptHash = keccak256(abi.encode(NitroUtils.getChannelId(fixedPart), receipt.leaderAmount, receipt.followerAmount, receipt.version));

        address firstSigner = NitroUtils.recoverSigner(
            receiptHash,
            receipt.leaderSig
        );
        address lastSigner = NitroUtils.recoverSigner(
            receiptHash,
            receipt.followerSig
        );
        require(firstSigner == fixedPart.participants[0], 'invalid signature from Leader');
        require(lastSigner == fixedPart.participants[fixedPart.participants.length - 1], 'invalid signature from Follower');
        return receipt;
    }

    function requireCorrectAdjustments(
        Outcome.SingleAssetExit[] memory oldOutcome,
        Outcome.SingleAssetExit[] memory newOutcome,
        MarginState memory receipt
    ) internal pure {
        // TODO: Validate collateral asset type is valid
        //
        // require(
        //     oldOutcome.length == 1 &&
        //         newOutcome.length == 1 &&
        //         oldOutcome[0].asset == address(0) &&
        //         newOutcome[0].asset == address(0),
        //     'only native asset allowed'
        // );

        require(
            newOutcome[0].allocations[uint256(AllocationIndices.Leader)].amount == receipt.leaderAmount,
            'Leader not adjusted correctly'
        );
        require(
            newOutcome[0].allocations[uint256(AllocationIndices.Follower)].amount == receipt.followerAmount,
            'Follower not adjusted correctly'
        );
        // TODO: Find a way to validate version number
    }
}
