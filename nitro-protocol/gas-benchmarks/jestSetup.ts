import {exec} from 'child_process';
import {promises, existsSync, truncateSync} from 'fs';

import waitOn from 'wait-on';
import kill from 'tree-kill';
import {BigNumber} from '@ethersproject/bignumber';
import {SnapshotRestorer, takeSnapshot} from '@nomicfoundation/hardhat-network-helpers';

import {deployContracts} from './localSetup';
import {TestChannel, challengeChannel} from './fixtures';

declare global {
  // eslint-disable-next-line @typescript-eslint/no-namespace
  namespace jest {
    interface Matchers<R> {
      toConsumeGas(benchmark: number): R;
    }
  }
}

const logFile = './hardhat-network-output.log';
const hardHatNetworkEndpoint = 'http://localhost:9546'; // the port should be unique

jest.setTimeout(30_000); // give hardhat network a chance to get going
if (existsSync(logFile)) truncateSync(logFile);
const hardhatProcess = exec('npx hardhat node --no-deploy --port 9546', (error, stdout) => {
  promises.appendFile(logFile, stdout);
});
const hardhatProcessExited = new Promise(resolve => hardhatProcess.on('exit', resolve));
const hardhatProcessClosed = new Promise(resolve => hardhatProcess.on('close', resolve));

let snapshot: SnapshotRestorer;

beforeAll(async () => {
  await waitOn({resources: [hardHatNetworkEndpoint]});

  await deployContracts();

  snapshot = await takeSnapshot();
});

beforeEach(async () => {
  await snapshot.restore();
  snapshot = await takeSnapshot();
});

afterAll(async () => {
  // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
  await kill(hardhatProcess.pid!);
  await hardhatProcessExited;
  await hardhatProcessClosed;
});

expect.extend({
  async toConsumeGas(
    received: any, // TransactionResponse
    benchmark: number
  ) {
    const {gasUsed: gasUsedBN} = await received.wait();
    const gasUsed = (gasUsedBN as BigNumber).toNumber();

    const pass = gasUsed === benchmark; // This could get replaced with a looser check with upper/lower bounds

    if (pass) {
      return {
        message: () => `expected to NOT consume ${benchmark} gas, but did`,
        pass: true,
      };
    } else {
      const format = (x: number) => {
        return x.toLocaleString().replace(/,/g, '_');
      };
      const green = (x: string) => `\x1b[32m${x}\x1b[0m`;
      const red = (x: string) => `\x1b[31m${x}\x1b[0m`;

      const diff = gasUsed - benchmark;
      const diffStr: string = diff > 0 ? red('+' + format(diff)) : green(format(diff));
      const diffPercent = `${Math.round((Math.abs(diff) / benchmark) * 100)}%`;

      return {
        message: () =>
          `expected to consume ${format(benchmark)} gas, but actually consumed ${format(
            gasUsed
          )} gas (${diffStr}, ${diffPercent}). Consider running npm run benchmark:update`,
        pass: false,
      };
    }
  },
});

/**
 * Constructs a support proof for the supplied channel, calls challenge,
 * and asserts the expected gas
 * @returns The proof and finalizesAt
 */
export async function challengeChannelAndExpectGas(
  channel: TestChannel,
  asset: string,
  expectedGas: number
): Promise<{proof: ReturnType<typeof channel.counterSignedSupportProof>; finalizesAt: number}> {
  const {challengeTx, proof, finalizesAt} = await challengeChannel(channel, asset);

  await expect(challengeTx).toConsumeGas(expectedGas);

  return {proof, finalizesAt};
}
