// SPDX-License-Identifier: MIT
pragma solidity 0.8.17;
pragma experimental ABIEncoderV2;

import './interfaces/IForceMoveApp.sol';
import './libraries/NitroUtils.sol';
import './interfaces/INitroTypes.sol';
import {ExitFormat as Outcome} from '@statechannels/exit-format/contracts/ExitFormat.sol';

/**
 * @dev The CurrencySwapApp contract complies with the ForceMoveApp interface and allows Alice and Bob (participants[0] to participants[n+1], where n is the number of intermediaries) to trade two assets atomically.
 */
contract CurrencySwapApp is IForceMoveApp {
    struct Order {
        address asset;
        uint256 amount;
    }

    struct Voucher {
        Order takerOrder;
        Order makerOrder;
        // TODO: I assume both side need to confirm the amount
        INitroTypes.Signature takerSignature;
        INitroTypes.Signature makerSignature;
        // TODO: to differentiate between Vouchers
        uint256 number;
    }

    //TODO: Remove once logic is tested
    struct VoucherAmountAndSignature {
        uint256 amount;
        INitroTypes.Signature signature; // signature on abi.encode(channelId,amount)
    }

    // TODO: To confirm indices are always 0 (Taker) and 1 (Maker)
    enum AllocationIndices {
        Taker, // Aggressor order
        Maker // Passive order was found in broker orderbook
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
        // TODO:
        // This channel has several states which can be supported:
        // 0    prefund
        // 1    postfund, both Maker and Taker have funded their side
        // 2    order created, both Taker and Maker agree on an order
        // 3a   order change
        // 3b   margin call, both Maker and Taker agree on outcome change due to some token / currency (from order) price change
        // 3c   order executed, both Maker and Taker agree on execution, margin changes cleared
        // 4    -> step 2
        //
        //
        // states 0,1,2 require only themselves to prove supported
        //
        // if (proof.length == 0) {
        //     require(
        //         NitroUtils.getClaimedSignersNum(candidate.signedBy) ==
        //             fixedPart.participants.length,
        //         '!unanimous; |proof|=0'
        //     );
        //     if (candidate.variablePart.turnNum == 0) return; // prefund
        //     if (candidate.variablePart.turnNum == 1) return; // postfund
        // }
        //
        // states
        //
        //
        // // states 0,1,3 can be supported via unanimous consensus:
        // if (proof.length == 0) {
        //     require(
        //         NitroUtils.getClaimedSignersNum(candidate.signedBy) ==
        //             fixedPart.participants.length,
        //         '!unanimous; |proof|=0'
        //     );
        //     if (candidate.variablePart.turnNum == 0) return; // prefund
        //     if (candidate.variablePart.turnNum == 1) return; // postfund
        //     if (candidate.variablePart.turnNum == 3) {
        //         // final (note: there is a core protocol escape hatch for this, too, so it could be removed)
        //         require(candidate.variablePart.isFinal, '!final; turnNum=3 && |proof|=0');
        //         return;
        //     }
        //     revert('bad candidate turnNum; |proof|=0');
        // }
        // // State 2 can be supported via a forced transition from state 1:
        // //
        // //      (2)_B     [redemption state signed by Bob, includes a voucher signed by Alice. The outcome may be updated in favour of Bob]
        // //      ^
        // //      |
        // //      (1)_AIB   [fully signed postfund]
        // if (proof.length == 1) {
        //     requireProofOfUnanimousConsensusOnPostFund(proof[0], fixedPart.participants.length);
        //     require(candidate.variablePart.turnNum == 2, 'bad candidate turnNum; |proof|=1');
        //     require(
        //         NitroUtils.isClaimedSignedBy(candidate.signedBy, 2),
        //         'redemption not signed by Bob'
        //     );
        //     uint256 voucherAmount = requireValidVoucher(candidate.variablePart.appData, fixedPart);
        //     requireCorrectAdjustments(
        //         proof[0].variablePart.outcome,
        //         candidate.variablePart.outcome,
        //         voucherAmount
        //     );
        //     return;
        // }
        // revert('bad proof length');
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

    function requireValidVoucher(
        bytes memory appData,
        FixedPart memory fixedPart
    ) internal pure returns (uint256) {
        VoucherAmountAndSignature memory voucher = abi.decode(appData, (VoucherAmountAndSignature));

        address signer = NitroUtils.recoverSigner(
            keccak256(abi.encode(NitroUtils.getChannelId(fixedPart), voucher.amount)),
            voucher.signature
        );
        //TODO: Verify amount, signature of Maker
        require(signer == fixedPart.participants[0], 'invalid signature for voucher'); // could be incorrect channelId or incorrect signature
        return voucher.amount;
    }

    function requireCorrectAdjustments(
        Outcome.SingleAssetExit[] memory oldOutcome,
        Outcome.SingleAssetExit[] memory newOutcome,
        uint256 voucherAmount
    ) internal pure {
        // TODO: require Outcome.length == 2;
        // && Outcome[0].asset == TakerAsset; Outcome[1].asset == MakerAsset
        require(
            oldOutcome.length == 1 &&
                newOutcome.length == 1 &&
                oldOutcome[0].asset == address(0) &&
                newOutcome[0].asset == address(0),
            'only native asset allowed'
        );

        // TODO: Adjustments as follows
        // Outcome[0].allocations[Taker].amount - voucherTakerAmount
        // Outcome[1].allocations[Taker].amount + voucherMakerAmount
        // Outcome[0].allocations[Maker].amount - voucherMakerAmount
        // Outcome[1].allocations[Maker].amount + voucherTakerAmount
        require(
            newOutcome[0].allocations[uint256(AllocationIndices.Taker)].amount ==
                oldOutcome[0].allocations[uint256(AllocationIndices.Taker)].amount - voucherAmount,
            'Alice not adjusted correctly'
        );
        require(
            newOutcome[0].allocations[uint256(AllocationIndices.Maker)].amount == voucherAmount,
            'Bob not adjusted correctly'
        );
    }
}
