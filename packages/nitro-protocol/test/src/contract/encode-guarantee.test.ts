import {it} from '@jest/globals';

import {encodeGuaranteeData, decodeGuaranteeData, Guarantee} from '../../../src/contract/outcome';

const guarantee: Guarantee = {
  left: '0x14bcc435f49d130d189737f9762feb25c44ef5b886bef833e31a702af6be4748',
  right: '0x14bcc435f49d130d189736bef833e31a702af6be47487f9762feb25c44ef5b88',
};

const description0 = 'Encodes and decodes guarantee';

describe('outcome', () => {
  describe('encoding and decoding', () => {
    it.each`
      description     | encodeFunction         | decodeFunction         | data
      ${description0} | ${encodeGuaranteeData} | ${decodeGuaranteeData} | ${guarantee}
    `('$description', (tc) => {
        const {encodeFunction, decodeFunction, data} = tc as { encodeFunction:Function, decodeFunction:Function, data:Guarantee };
      const encodedData = encodeFunction(data);
      const decodedData = decodeFunction(encodedData);
      expect(decodedData).toEqual(data);
    });
  });
});
