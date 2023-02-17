pragma solidity 0.8.17;
pragma experimental ABIEncoderV2;

import './interfaces/IForceMoveApp.sol';
import './libraries/NitroUtils.sol';
import {ExitFormat as Outcome} from '@statechannels/exit-format/contracts/ExitFormat.sol';

// LedgerFinancingApp is a ForceMoveApp that allows a intermediary to earn interest
// on a deposit. It functions as a ConsensusApp with the following additional rule:
//   - the intermediary can unilaterally transition from state n to state n+1,
//     forcing calculated interest into the intermediary's Outcome allocation
contract LedgerFinancingApp is IForceMoveApp {
    struct Funds {
        address[] asset;
        uint256[] amount;
    }

    struct InterestAppData {
        // the per-block interest rate, expressed as the denominator of a fraction
        // eg, 1% per block would be expressed by bpy = 100,
        //     2% per block would be expressed by bpy = 50,
        //     0.1% per block would be expressed by bpy = 1000,
        // etc.
        uint256 bpy;
        // the block number of the latest principal adjustment
        uint256 blocknumber;
        // the current principal. Decreases as the serviceProvider earns via the channel.
        Funds principal;
        // the total interest collected so far. Strictly increasing.
        // The value of the intermediary's allocation can never be less than this:
        // ie, when the intermediary's collectedInterest grows to be equal to their
        // allocation, the channel is effectively spent and can be concluded.
        Funds collectedInterest;
    }

    enum AllocationIndicies {
        serviceProvider, // borrower: recovers service fees from intermediary's deposit
        intermediary // financier: makes initial deposit and earns interest
    }

    function requireStateSupported(
        FixedPart calldata,
        RecoveredVariablePart[] calldata proof,
        RecoveredVariablePart calldata candidate
    ) external view override {
        if (proof.length == 0) {
            // unanimous consensus check
            require(
                NitroUtils.getClaimedSignersNum(candidate.signedBy) == 2,
                '!unanimous; |proof|=0'
            );
            return;
        } else if (proof.length == 1) {
            // check that proof[0] -> candidate respects the stated interest rate.
            // Requires:
            //  - proof state is unanimous
            //  - candidate state immediately follows proof state (by turnNum)
            //  - the intermediary has not taken more funds than owed according
            //    to the interest rate agreement of the channel
            require(
                NitroUtils.getClaimedSignersNum(proof[0].signedBy) == 2,
                '!unanimous proof state'
            );
            require(
                proof[0].variablePart.turnNum + 1 == candidate.variablePart.turnNum,
                'turn(candidate) != turn(proof)+1'
            );

            Funds memory outstandingInterest = computeOutstandingInterest(
                abi.decode(proof[0].variablePart.appData, (InterestAppData))
            );
            requireFairOutcomeAdjustment(
                proof[0].variablePart.outcome,
                candidate.variablePart.outcome,
                outstandingInterest
            );
        } else {
            revert('|proof| > 1');
        }
    }

    // The outstanding interest is calculated based on:
    //  - the latest consensus principal
    //  - the channel's interest rate
    //  - the time elapsed since the last principal adjustment
    function computeOutstandingInterest(InterestAppData memory appData)
        private
        view
        returns (Funds memory)
    {
        uint256 numBlocks = block.number - appData.blocknumber;

        address[] memory assets = new address[](appData.principal.asset.length);
        uint256[] memory amounts = new uint256[](appData.principal.asset.length);

        Funds memory outstanding = Funds(assets, amounts);

        // copy all assets from the principal, and multiply by the interest rate
        for (uint256 i = 0; i < appData.principal.asset.length; i++) {
            outstanding.asset[i] = appData.principal.asset[i];
            outstanding.amount[i] = (appData.principal.amount[i] * numBlocks) / appData.bpy;
        }

        return outstanding;
    }

    // Ensures that the given outcome does not unfairly allocate to the intermediary.
    function requireFairOutcomeAdjustment(
        Outcome.SingleAssetExit[] memory initialOutcome,
        Outcome.SingleAssetExit[] memory finalOutcome,
        Funds memory outstandingInterest
    ) private pure {
        for (uint256 i = 0; i < outstandingInterest.asset.length; i++) {
            address asset = outstandingInterest.asset[i];
            uint256 earned = outstandingInterest.amount[i];

            for (uint256 j = 0; j < finalOutcome.length; j++) {
                if (finalOutcome[j].asset == asset) {
                    require(initialOutcome[j].asset == asset, 'Asset mismatch');

                    requireFairAssetAdjustment(initialOutcome[j], finalOutcome[j], earned);
                }
            }
        }
    }

    // Ensures that the given asset outcome does not unfairly allocate to the intermediary.
    function requireFairAssetAdjustment(
        Outcome.SingleAssetExit memory initial,
        Outcome.SingleAssetExit memory adjusted,
        uint256 earned
    ) private pure {
        require(
            initial.allocations[uint256(AllocationIndicies.serviceProvider)].destination ==
                adjusted.allocations[uint256(AllocationIndicies.serviceProvider)].destination,
            'payee mismatch'
        );
        uint256 initialProviderBalance = initial
            .allocations[uint256(AllocationIndicies.serviceProvider)]
            .amount;
        uint256 adjustedProviderBalance = adjusted
            .allocations[uint256(AllocationIndicies.serviceProvider)]
            .amount;
        uint256 claimed = initialProviderBalance - adjustedProviderBalance;

        require(claimed <= earned, 'earned<claimed');
    }

    function daysSince(uint256 blocknumber) private view returns (uint32) {
        return uint32((block.number - blocknumber) / 7200); // 7200 == 24*60*60/12
    }
}
