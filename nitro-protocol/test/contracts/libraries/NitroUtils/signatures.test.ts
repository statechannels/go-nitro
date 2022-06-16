import {Contract, Wallet, ethers} from 'ethers';
const {arrayify, id} = ethers.utils;

import NitroUtilsArtifact from '../../../../artifacts/contracts/test/TESTNitroUtils.sol/TESTNitroUtils.json';
import {getTestProvider, setupContract} from '../../../test-helpers';
import {sign} from '../../../../src/signatures';
const provider = getTestProvider();
let NitroUtils: Contract;

const participants = ['', '', ''];
const wallets = new Array(3);

// Populate wallets and participants array
for (let i = 0; i < 3; i++) {
  wallets[i] = Wallet.createRandom();
  participants[i] = wallets[i].address;
}

beforeAll(async () => {
  NitroUtils = setupContract(provider, NitroUtilsArtifact, process.env.TEST_NITRO_UTILS_ADDRESS);
});

describe('_recoverSigner', () => {
  it('recovers the signer correctly', async () => {
    // Following https://docs.ethers.io/ethers.js/html/cookbook-signing.html
    const privateKey = '0x0123456789012345678901234567890123456789012345678901234567890123';
    const wallet = new Wallet(privateKey);
    const msgHash = id('Hello World');
    const msgHashBytes = arrayify(msgHash);
    const sig = await sign(wallet, msgHashBytes);
    expect(await NitroUtils.recoverSigner(msgHash, sig)).toEqual(wallet.address);
  });
});

describe('isSignedBy', () => {
  it('returns true when a participant bit is set', async () => {
    expect(await NitroUtils.isSignedBy(0b101, 0)).toBe(true);
  });
  it('returns false when a participant bit is not set', async () => {
    expect(await NitroUtils.isSignedBy(0b101, 1)).toBe(false);
  });
});
