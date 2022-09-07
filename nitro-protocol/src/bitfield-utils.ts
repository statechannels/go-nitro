import {Uint256} from './contract/types';

export function getSignersNum(signedByStr: Uint256): number {
  let signedBy = parseFloat(signedByStr);
  let amount = 0;

  for (; signedBy > 0; amount++) {
    signedBy &= signedBy - 1;
  }

  return amount;
}

export function getSignersIndices(signedByStr: Uint256): number[] {
  let signedBy = parseFloat(signedByStr);

  const signerIndices: number[] = [];
  let signerNum = 0;
  let acceptedSigners = 0;

  for (; signedBy > 0; signerNum++) {
    if (signedBy % 2 == 1) {
      signerIndices[acceptedSigners] = signerNum;
      acceptedSigners++;
    }
    signedBy >>= 1;
  }

  return signerIndices;
}

export function getSignedBy(signerIndices: number[]): Uint256 {
  let signedBy = 0;
  for (const sIdx of signerIndices) {
    signedBy += 2 ** sIdx;
  }
  return signedBy.toString();
}
