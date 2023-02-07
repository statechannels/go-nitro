pragma solidity 0.8.17;
pragma experimental ABIEncoderV2;

import './interfaces/IForceMoveApp.sol';
import './libraries/NitroUtils.sol';
import {ExitFormat as Outcome} from '@statechannels/exit-format/contracts/ExitFormat.sol';

// LedgerFinancingApp is a ForceMoveApp that allows a intermediary to earn interest
// on a deposit. It fuctions as a ConsensusApp with the following additional rule:
//   - the intermediary can unilaterally transition from state n to state n+1,
//     forcing calculated interest into the intermediary's Outcome allocation
contract LedgerFinancingApp is IForceMoveApp {
    struct Funds {
        address[] asset;
        uint256[] amount;
    }

    struct InterestAppData {
        // a per-day simple interest rate (daily percentage yield), as a fraction.
        // ie, 1% per day is represented with dpyNum = 1 and dpyDen = 100.
        //
        // Lower numbers (eg, fraction in simplest terms) produce least risk of overflow.
        uint128 dpyNum;
        uint128 dpyDen;
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
        intermediary, // financier: makes initial deposit and earns interest
        serviceProvider // borrower: recovers service fees from intermediary's deposit
    }

    // Ensures that the given outcome does not unfairly allocate to the intermediary.
    function requireOutcomeIsEarned(
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
                    require(
                        initialOutcome[j]
                            .allocations[uint256(AllocationIndicies.intermediary)]
                            .destination ==
                            finalOutcome[j]
                                .allocations[uint256(AllocationIndicies.intermediary)]
                                .destination,
                        'payee mismatch'
                    );
                    // combine prior balance with outstandingInterest
                    earned += initialOutcome[j]
                        .allocations[uint256(AllocationIndicies.intermediary)]
                        .amount;
                    // and compare against the claimed outcome
                    uint256 claimed = finalOutcome[j]
                        .allocations[uint256(AllocationIndicies.intermediary)]
                        .amount;

                    require(claimed <= earned, 'earned<claimed');
                }
            }
        }
    }

    function daysSince(uint256 blocknumber) private view returns (uint32) {
        return uint32((block.number - blocknumber) / 7200); // 7200 == 24*60*60/12
    }

    // The outstanding interest is calculated based on:
    //  - the latest consensus principal
    //  - the channel's interest rate
    //  - the time elapsed since the last principal adjustment
    function getOutstandingInterest(InterestAppData memory appData)
        private
        view
        returns (Funds memory)
    {
        uint32 numDays = daysSince(appData.blocknumber);
        uint256 simpleInterestNum = numDays * appData.dpyNum;

        address[] memory assets = new address[](appData.principal.asset.length);
        uint256[] memory amounts = new uint256[](appData.principal.asset.length);

        Funds memory outstanding = Funds(assets, amounts);

        // copy all assets from the principal, and multiply by the interest rate
        for (uint256 i = 0; i < appData.principal.asset.length; i++) {
            outstanding.asset[i] = appData.principal.asset[i];
            outstanding.amount[i] =
                (appData.principal.amount[i] * simpleInterestNum) /
                appData.dpyDen;
        }

        return outstanding;
    }

    function requireStateSupported(
        FixedPart calldata fixedPart,
        RecoveredVariablePart[] calldata proof,
        RecoveredVariablePart calldata candidate
    ) external view override {
        if (proof.length == 0) {
            // unanimous consensus check
            require(
                NitroUtils.getClaimedSignersNum(candidate.signedBy) ==
                    fixedPart.participants.length,
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

            Funds memory outstandingInterest = getOutstandingInterest(
                abi.decode(proof[0].variablePart.appData, (InterestAppData))
            );
            requireOutcomeIsEarned(
                proof[0].variablePart.outcome,
                candidate.variablePart.outcome,
                outstandingInterest
            );
        } else {
            revert('|proof| > 1'); // does it pay to be this terse with revert messages?
        }
    }
}
