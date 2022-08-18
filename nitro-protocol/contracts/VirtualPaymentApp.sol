// SPDX-License-Identifier: MIT
pragma solidity 0.7.6;
pragma experimental ABIEncoderV2;

import './interfaces/IForceMoveApp.sol';
import './libraries/NitroUtils.sol';
import './interfaces/INitroTypes.sol';

/**
 * @dev The VirtualPaymentApp contract complies with the ForceMoveApp interface and allows payments to be made virtually from Alice to Bob (participants[0] to participants[1]).
 */
contract VirtualPaymentApp is IForceMoveApp {
    struct VoucherAmountAndSignature {
        uint256 amount;
        INitroTypes.Signature signature; // signature on abi.encode(channelId,amount)
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
            revert('bad candidate turnNum');
        }

        // State 2 can be supported via a forced transition from state 1:
        //
        //      (2)_B     [redemption state signed by Bob, includes a voucher signed by Alice. The outcome may be updated in favour of Bob]
        //      ^
        //      |
        //      (1)_AIB   [fully signed postfund]

        if (proof.length == 1) {
            requireProofOfUnanimousConsensusOnPostFund(proof[0], fixedPart.participants.length);

            require(candidate.variablePart.turnNum == 2, 'invalid transition; |proof|=1');

            require(
                NitroUtils.isClaimedSignedBy(candidate.signedBy, 2),
                'redemption not signed by Bob'
            );

            requireValidVoucher(candidate.variablePart.appData, fixedPart);

            // TODO remove assumption about single asset, and factor into CheckAliceAndBobOutcomes (don't use magic indices, potentially validate destinations)
            // we want to be sure that the voucher amount wasn't greater than alice's original balance. This should be handled by the underflow protection which is now a part of solidity.
            require(
                candidate.variablePart.outcome[0].allocations[0].amount ==
                    proof[0].variablePart.outcome[0].allocations[0].amount - voucher.amount,
                'Alice not adjusted correctly'
            );
            require(
                candidate.variablePart.outcome[0].allocations[1].amount == voucher.amount,
                'Bob not adjusted correctly'
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

    function requireValidVoucher(bytes memory appData, FixedPart memory fixedPart) internal pure {
        VoucherAmountAndSignature memory voucher = abi.decode(appData, (VoucherAmountAndSignature));

        address signer = NitroUtils.recoverSigner(
            keccak256(abi.encode(NitroUtils.getChannelId(fixedPart), voucher.amount)),
            voucher.signature
        );
        require(signer == fixedPart.participants[0], 'irrelevant voucher'); // could be incorrect channelId or incorrect signature
    }
}
