// SPDX-License-Identifier: MIT
pragma solidity 0.7.6;
pragma experimental ABIEncoderV2;

import '../interfaces/INitroTypes.sol';
import {NitroUtils} from '../libraries/NitroUtils.sol';


/**
 * @dev This contract extends the NitroUtils contract to enable it to be more easily unit-tested. It exposes public or external functions call into internal functions. It should not be deployed in a production environment.
 */
contract TESTNitroUtils {


    /**
     * @dev Wrapper for otherwise internal function. Given a digest and digital signature, recover the signer
     * @param _d message digest
     * @param sig ethereum digital signature
     * @return signer
     */
    function recoverSigner(bytes32 _d, INitroTypes.Signature memory sig) public pure returns (address) {
        return NitroUtils.recoverSigner(_d, sig);
    }

    /**
     * @notice Check if supplied participantIndex bit is set to 1 in signedBy bit mask.
     * @dev Check if supplied partitipationIndex bit is set to 1 in signedBy bit mask.
     * @param signedBy Bit mask field to check.
     * @param participantIndex Bit to check.
     */
    function isSignedBy(uint256 signedBy, uint8 participantIndex) public pure returns (bool) {
        return NitroUtils.isSignedBy(signedBy,participantIndex);
    }


}
