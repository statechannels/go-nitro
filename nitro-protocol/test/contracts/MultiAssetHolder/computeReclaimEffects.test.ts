import {BigNumber, BytesLike, constants} from 'ethers';
import {Allocation, AllocationType} from '@statechannels/exit-format';
import {it} from '@jest/globals';

import {getTestProvider, setupContract} from '../../test-helpers';
import {TESTNitroAdjudicator} from '../../../typechain-types/TESTNitroAdjudicator';
// eslint-disable-next-line import/order
import TESTNitroAdjudicatorArtifact from '../../../artifacts/contracts/test/TESTNitroAdjudicator.sol/TESTNitroAdjudicator.json';

const testNitroAdjudicator = setupContract(
  getTestProvider(),
  TESTNitroAdjudicatorArtifact,
  process.env.TEST_NITRO_ADJUDICATOR_ADDRESS
) as unknown as TESTNitroAdjudicator;

import {computeReclaimEffects} from '../../../src/contract/multi-asset-holder';
import {encodeGuaranteeData} from '../../../src/contract/outcome';

const Alice = '0x000000000000000000000000000000000000000000000000000000000000000a';
const Bob = '0x000000000000000000000000000000000000000000000000000000000000000b';

interface TestCaseInputs {
  sourceAllocations: Allocation[];
  targetAllocations: Allocation[];
  indexOfTargetInSource: number;
}

interface TestCaseOutputs {
  newSourceAllocations: Allocation[];
}
interface TestCase {
  inputs: TestCaseInputs;
  outputs: TestCaseOutputs;
}

interface AllocationT {
  destination: string;
  amount: BigNumber;
  allocationType: number;
  metadata: BytesLike;
}

const testcase1: TestCase = {
  inputs: {
    indexOfTargetInSource: 2,
    sourceAllocations: [
      {
        destination: Alice,
        amount: '0x02',
        allocationType: AllocationType.simple,
        metadata: '0x',
      },
      {
        destination: Bob,
        amount: '0x02',
        allocationType: AllocationType.simple,
        metadata: '0x',
      },
      {
        destination: constants.HashZero,
        amount: '0x06',
        allocationType: AllocationType.guarantee,
        metadata: encodeGuaranteeData({left: Alice, right: Bob}),
      },
    ],
    targetAllocations: [
      {
        destination: Alice,
        amount: '0x01',
        allocationType: AllocationType.simple,
        metadata: '0x',
      },
      {
        destination: Bob,
        amount: '0x05',
        allocationType: AllocationType.simple,
        metadata: '0x',
      },
    ],
  },
  outputs: {
    newSourceAllocations: [
      {
        destination: Alice,
        amount: '0x03',
        allocationType: AllocationType.simple,
        metadata: '0x',
      },
      {
        destination: Bob,
        amount: '0x07',
        allocationType: AllocationType.simple,
        metadata: '0x',
      },
    ],
  },
};

const testCases: TestCase[][] = [[testcase1]];

describe('computeReclaimEffects', () => {
  it.each(testCases)('off chain method matches expectation', (testCase: TestCase) => {
    const offChainNewSourceAllocations = computeReclaimEffects(
      testCase.inputs.sourceAllocations,
      testCase.inputs.targetAllocations,
      testCase.inputs.indexOfTargetInSource
    );

    expect(offChainNewSourceAllocations).toMatchObject(testCase.outputs.newSourceAllocations);
  });

  it.each(testCases)('on chain method matches expectation', async (testCase: TestCase) => {
    const onChainNewSourceAllocations = await testNitroAdjudicator.compute_reclaim_effects(
      testCase.inputs.sourceAllocations,
      testCase.inputs.targetAllocations,
      testCase.inputs.indexOfTargetInSource
    );

    expect(onChainNewSourceAllocations.map(convertAmountToHexString)).toMatchObject(
      testCase.outputs.newSourceAllocations
    );
  });

  const convertAmountToHexString = (a: AllocationT) => ({...a, amount: a.amount.toHexString()});
});
