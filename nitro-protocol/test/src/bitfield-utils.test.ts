import {getSignersNum, getSignersIndices, getSignedBy} from '../../src/bitfield-utils';

describe('bitfield utils', () => {
  it('getSignersNum', () => {
    expect(getSignersNum('0')).toEqual(0);
    expect(getSignersNum('1')).toEqual(1);
    expect(getSignersNum('2')).toEqual(1);
    expect(getSignersNum('3')).toEqual(2);
    expect(getSignersNum('5')).toEqual(2);
    expect(getSignersNum('100')).toEqual(3);
    expect(getSignersNum('1000')).toEqual(6);
  });

  it('getSignersIndices', () => {
    expect(getSignersIndices('0')).toEqual([]);
    expect(getSignersIndices('1')).toEqual([0]);
    expect(getSignersIndices('2')).toEqual([1]);
    expect(getSignersIndices('3')).toEqual([0, 1]);
    expect(getSignersIndices('5')).toEqual([0, 2]);
    expect(getSignersIndices('100')).toEqual([2, 5, 6]);
    expect(getSignersIndices('1000')).toEqual([3, 5, 6, 7, 8, 9]);
  });

  it('getSignedBy', () => {
    expect(getSignedBy(0)).toEqual('1');
    expect(getSignedBy(1)).toEqual('2');
    expect(getSignedBy(5)).toEqual('32');

    expect(getSignedBy([])).toEqual('0');
    expect(getSignedBy([0])).toEqual('1');
    expect(getSignedBy([1])).toEqual('2');
    expect(getSignedBy([0, 1])).toEqual('3');
    expect(getSignedBy([0, 2])).toEqual('5');
    expect(getSignedBy([2, 5, 6])).toEqual('100');
    expect(getSignedBy([3, 5, 6, 7, 8, 9])).toEqual('1000');
  });
});
