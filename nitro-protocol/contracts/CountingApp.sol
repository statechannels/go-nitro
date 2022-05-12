// SPDX-License-Identifier: MIT
pragma solidity 0.7.6;
pragma experimental ABIEncoderV2;

import './interfaces/IForceMoveApp.sol';
import {ExitFormat as Outcome} from '@statechannels/exit-format/contracts/ExitFormat.sol';

/**
 * @dev The CountingApp contracts complies with the ForceMoveApp interface and allows only for a simple counter to be incremented. Used for testing purposes.
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
        // TODO see https://github.com/statechannels/go-nitro/issues/558
        return signedVariableParts[signedVariableParts.length - 1].variablePart;
    }
}
