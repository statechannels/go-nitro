// SPDX-License-Identifier: MIT
pragma solidity 0.7.6;
pragma experimental ABIEncoderV2;

import '../interfaces/INitroTypes.sol';
import {ShortcuttingTurnTaking} from '../libraries/signature-logic/ShortcuttingTurnTaking.sol';

/**
 * @dev This contract extends the ShortcuttingTurnTaking contract to enable it to be more easily unit-tested. It exposes public or external functions call into internal functions. It should not be deployed in a production environment.
 */
contract TESTShortcuttingTurnTaking {
    function requireValidTurnTaking(
        INitroTypes.FixedPart memory fixedPart,
        INitroTypes.SignedVariablePart[] memory signedVariableParts
    ) public pure returns (bool) {
        ShortcuttingTurnTaking.requireValidTurnTaking(fixedPart, signedVariableParts);
        return true;
    }
}
