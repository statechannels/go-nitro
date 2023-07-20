pragma solidity 0.8.17;
import '../MultiAssetHolder.sol';

import '@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol';
import '@openzeppelin/contracts/token/ERC20/IERC20.sol';

/**
@dev This contract is used to batch deposit ERC20 tokens into .
 */
contract BatchOperator {
    using SafeERC20 for IERC20;
    MultiAssetHolder public adjudicator;

    constructor(address adjudicatorAddr) {
        adjudicator = MultiAssetHolder(adjudicatorAddr);
    }

    /**
     * @dev Deposits ETH (or native token for other runtime) into the adjudicator for multiple channels.
     */
    function deposit_batch_eth(
        bytes32[] calldata channelIds,
        uint256[] calldata expectedHelds,
        uint256[] calldata amounts
    ) external payable virtual {
        require(
            channelIds.length == expectedHelds.length && expectedHelds.length == amounts.length,
            'Array lengths must match'
        );
        for (uint256 i = 0; i < channelIds.length; i++) {
            adjudicator.deposit{value: amounts[i]}(
                address(0),
                channelIds[i],
                expectedHelds[i],
                amounts[i]
            );
        }
    }

    /**
     * @dev Deposits ERC20 tokens into the adjudicator for multiple channels.
     */
    function deposit_batch_erc(
        address asset,
        bytes32[] calldata channelIds,
        uint256[] calldata expectedHelds,
        uint256[] calldata amounts,
        uint256 totalAmount
    ) external payable virtual {
        require(
            channelIds.length == expectedHelds.length && expectedHelds.length == amounts.length,
            'Array lengths must match'
        );
        IERC20(asset).safeIncreaseAllowance(msg.sender, totalAmount);

        for (uint256 i = 0; i < channelIds.length; i++) {
            adjudicator.deposit(asset, channelIds[i], expectedHelds[i], amounts[i]);
        }
    }
}
