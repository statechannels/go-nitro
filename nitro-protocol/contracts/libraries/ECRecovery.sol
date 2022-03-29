//SPDX-License-Identifier: MIT License
pragma solidity 0.7.4;

library ECRecovery {

  /**
   * @dev Recover signer address from a message by using his signature
   * @param hash bytes32 message, the hash is the signed message. What is recovered is the signer address.
   * @param v v part of a signature
   * @param r r part of a signature
   * @param s s part of a signature
   */
  function recover(bytes32 hash, uint8 v, bytes32 r, bytes32 s) internal pure returns (address) {

    // Version of signature should be 27 or 28, but 0 and 1 are also possible versions
    if (v < 27) {
      v += 27;
    }

    // If the version is correct return the signer address
    if (v != 27 && v != 28) {
      return (address(0));
    } else {
      return ecrecover(hash, v, r, s);
    }
  }

}
