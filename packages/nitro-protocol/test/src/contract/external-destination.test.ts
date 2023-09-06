import {it} from '@jest/globals';

import {isExternalDestination} from '../../../src/contract/channel';

describe('isExternalDestination', () => {
  it.each`
    bytes32                                                                 | result
    ${'0x0'}                                                                | ${false}
    ${'0x0000000000000000000000002F0E2cB3c2c98E6AfB89A8c50cbEF0cB6B3DC35c'} | ${true}
    ${'0x0000000040000000000000002F0E2cB3c2c98E6AfB89A8c50cbEF0cB6B3DC35c'} | ${false}
  `('$bytes32 -- $result', (tc) => {
    const {bytes32, result} = tc as {bytes32: string; result: boolean};
    expect(isExternalDestination(bytes32)).toBe(result);
  });
});
