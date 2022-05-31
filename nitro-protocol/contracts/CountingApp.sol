// SPDX-License-Identifier: MIT
pragma solidity 0.7.6;
pragma experimental ABIEncoderV2;

import {ExitFormat as Outcome} from '@statechannels/exit-format/contracts/ExitFormat.sol';
import {TurnTaking} from './libraries/signature-logic/TurnTaking.sol';
import './interfaces/IForceMoveApp.sol';

/**
 * @dev The CountingApp contract complies with the ForceMoveApp interface and TurnTaking logic and allows only for a simple counter to be incremented. Used for testing purposes.
 */
contract CountingApp is IForceMoveApp {
    struct CountingAppData {
        uint256 counter;
    }

    /**
     * @notice Decodes the appData.
     * @dev Decodes the appData.
     * @param appDataBytes The abi.encode of a CountingAppData struct describing the application-specific data.
     * @return A CountingAppData struct containing the application-specific data.
     */
    function appData(bytes memory appDataBytes) internal pure returns (CountingAppData memory) {
        bytes memory decodedAppData = abi.decode(appDataBytes, (bytes));
        return abi.decode(decodedAppData, (CountingAppData));
    }

    /**
     * @notice Encodes application-specific rules for a particular ForceMove-compliant state channel.
     * @dev Encodes application-specific rules for a particular ForceMove-compliant state channel.
     * @param fixedPart Fixed part of the state channel.
     * @param signedVariableParts Array of variable parts to find the latest of.
     * @return VariablePart Latest supported by application variable part from supplied array.
     */    
    function latestSupportedState(
        FixedPart calldata fixedPart,
        SignedVariablePart[] calldata signedVariableParts
    ) external pure override returns (VariablePart memory) {
        TurnTaking.requireValidTurnTaking(fixedPart, signedVariableParts);

        for (uint i = 1; i < signedVariableParts.length; i++) {
            _requireIncrementedCounter(signedVariableParts[i], signedVariableParts[i-1]);
            _requireEqualOutcomes(signedVariableParts[i], signedVariableParts[i-1]);
        }

        return signedVariableParts[signedVariableParts.length - 1].variablePart;
    }

    /**
     * @notice Checks that counter encoded in first variable part equals an incremented counter in second variable part.
     * @dev Checks that counter encoded in first variable part equals an incremented counter in second variable part.
     * @param b SignedVariablePart with incremented counter.
     * @param a SignedVariablePart with counter before incrementing.
     */
    function _requireIncrementedCounter(
        SignedVariablePart memory b,
        SignedVariablePart memory a
    ) internal pure {
        require(
            appData(b.variablePart.appData).counter == appData(a.variablePart.appData).counter + 1,
            'Counter must be incremented'
        );
    }

    /**
     * @notice Checks that supplied signed variable parts contain the same outcome.
     * @dev Checks that supplied signed variable parts contain the same outcome.
     * @param a First SignedVariablePart.
     * @param b Second SignedVariablePart.
     */
    function _requireEqualOutcomes(
        SignedVariablePart memory a,
        SignedVariablePart memory b
    ) internal pure {
        require(
            Outcome.exitsEqual(a.variablePart.outcome, b.variablePart.outcome),
            'Outcome must not change'
        );
    }
}
