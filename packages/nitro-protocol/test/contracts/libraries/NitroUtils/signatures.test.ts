import {Contract, Wallet, ethers} from 'ethers';
const {arrayify, id} = ethers.utils;

import NitroUtilsArtifact from '../../../../artifacts/contracts/test/TESTNitroUtils.sol/TESTNitroUtils.json';
import {getTestProvider, setupContract} from '../../../test-helpers';
import {sign} from '../../../../src/signatures';
import {TESTNitroUtils} from '../../../../typechain-types';
const provider = getTestProvider();
let NitroUtils: Contract & TESTNitroUtils;

beforeAll(async () => {
  NitroUtils = setupContract(
    provider,
    NitroUtilsArtifact,
    process.env.TEST_NITRO_UTILS_ADDRESS
  ) as Contract & TESTNitroUtils;
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

describe('isClaimedSignedBy', () => {
  // prettier-ignore
  it('returns true when a participant bit is set', async () => {
    expect(await NitroUtils.isClaimedSignedBy(0b101     ,0)).toBe(true);
    expect(await NitroUtils.isClaimedSignedBy(0b101     ,2)).toBe(true);
    expect(await NitroUtils.isClaimedSignedBy(0b001     ,0)).toBe(true);
    expect(await NitroUtils.isClaimedSignedBy(0b10000000,7)).toBe(true);
    expect(await NitroUtils.isClaimedSignedBy(8         ,3)).toBe(true);
  });
  // prettier-ignore
  it('returns false when a participant bit is not set', async () => {
    expect(await NitroUtils.isClaimedSignedBy(0b101     ,1)).toBe(false);
    expect(await NitroUtils.isClaimedSignedBy(0b001     ,3)).toBe(false);
    expect(await NitroUtils.isClaimedSignedBy(0b001     ,2)).toBe(false);
    expect(await NitroUtils.isClaimedSignedBy(0b001     ,1)).toBe(false);
  });
});

describe('isClaimedSignedOnlyBy', () => {
  // prettier-ignore
  it('returns true when only that participant bit is set', async () => {
    expect(await NitroUtils.isClaimedSignedOnlyBy(0b001     ,0)).toBe(true);
    expect(await NitroUtils.isClaimedSignedOnlyBy(0b10000000,7)).toBe(true);
    expect(await NitroUtils.isClaimedSignedOnlyBy(8         ,3)).toBe(true);
  });
  // prettier-ignore
  it('returns false when that participant bit is not set', async () => {
    expect(await NitroUtils.isClaimedSignedOnlyBy(0b011     ,0)).toBe(false);
    expect(await NitroUtils.isClaimedSignedOnlyBy(0b10010000,7)).toBe(false);
    expect(await NitroUtils.isClaimedSignedOnlyBy(9         ,3)).toBe(false);
    expect(await NitroUtils.isClaimedSignedOnlyBy(0b101     ,0)).toBe(false);
    expect(await NitroUtils.isClaimedSignedOnlyBy(0b101     ,2)).toBe(false);
    expect(await NitroUtils.isClaimedSignedOnlyBy(0b101     ,1)).toBe(false);
    expect(await NitroUtils.isClaimedSignedOnlyBy(0b001     ,3)).toBe(false);
    expect(await NitroUtils.isClaimedSignedOnlyBy(0b001     ,2)).toBe(false);
    expect(await NitroUtils.isClaimedSignedOnlyBy(0b001     ,1)).toBe(false);
  });
});

describe('getClaimedSignersNum', () => {
  // prettier-ignore
  it('counts the number of signers correctly', async () => {
    expect(await NitroUtils.getClaimedSignersNum(0b001)).toEqual(1)
    expect(await NitroUtils.getClaimedSignersNum(0b011)).toEqual(2)
    expect(await NitroUtils.getClaimedSignersNum(0b101)).toEqual(2)
    expect(await NitroUtils.getClaimedSignersNum(0b111)).toEqual(3)
    expect(await NitroUtils.getClaimedSignersNum(0b000)).toEqual(0)
  });
});

describe('getClaimedSignersIndices', () => {
  // prettier-ignore
  it('returns the correct indices', async () => {
    expect(await NitroUtils.getClaimedSignersIndices(0b001)).toEqual([0])
    expect(await NitroUtils.getClaimedSignersIndices(0b011)).toEqual([0,1])
    expect(await NitroUtils.getClaimedSignersIndices(0b101)).toEqual([0,2])
    expect(await NitroUtils.getClaimedSignersIndices(0b111)).toEqual([0,1,2])
    expect(await NitroUtils.getClaimedSignersIndices(0b000)).toEqual([])
  });
});
