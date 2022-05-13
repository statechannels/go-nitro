import {exec} from 'child_process';
import {promises, existsSync, truncateSync} from 'fs';

import {ContractFactory, Contract} from '@ethersproject/contracts';
import {providers, utils} from 'ethers';
import waitOn from 'wait-on';
import kill from 'tree-kill';
import {BigNumber} from '@ethersproject/bignumber';

import nitroAdjudicatorArtifact from '../artifacts/contracts/NitroAdjudicator.sol/NitroAdjudicator.json';
import tokenArtifact from '../artifacts/contracts/Token.sol/Token.json';
import trivialAppArtifact from '../artifacts/contracts/TrivialApp.sol/TrivialApp.json';
import {NitroAdjudicator} from '../typechain-types/NitroAdjudicator';
import {Token} from '../typechain-types/Token';
import {TrivialApp} from '../typechain-types/TrivialApp';
declare global {
  // eslint-disable-next-line @typescript-eslint/no-namespace
  namespace jest {
    interface Matchers<R> {
      toConsumeGas(benchmark: number): R;
    }
  }
}

export let nitroAdjudicator: NitroAdjudicator & Contract;
export let token: Token & Contract;
export let trivialApp: TrivialApp & Contract;

const logFile = './hardhat-network-output.log';
const hardHatNetworkEndpoint = 'http://localhost:9546'; // the port should be unique

jest.setTimeout(15_000); // give hardhat network a chance to get going
if (existsSync(logFile)) truncateSync(logFile);
const hardhatProcess = exec('npx hardhat node --no-deploy --port 9546', (error, stdout) => {
  promises.appendFile(logFile, stdout);
});
const hardhatProcessExited = new Promise(resolve => hardhatProcess.on('exit', resolve));
const hardhatProcessClosed = new Promise(resolve => hardhatProcess.on('close', resolve));

export const provider = new providers.JsonRpcProvider(hardHatNetworkEndpoint);

let snapshotId = 0;

const tokenFactory = new ContractFactory(tokenArtifact.abi, tokenArtifact.bytecode).connect(
  provider.getSigner(0)
);

const nitroAdjudicatorFactory = new ContractFactory(
  nitroAdjudicatorArtifact.abi,
  nitroAdjudicatorArtifact.bytecode
).connect(provider.getSigner(0));

const trivialAppFactory = new ContractFactory(
  trivialAppArtifact.abi,
  trivialAppArtifact.bytecode
).connect(provider.getSigner(0));

export const trivialAppAddress = utils.getContractAddress({
  from: '0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266', // ASSUME: deployed by hardhat account 0
  nonce: 0, // ASSUME: this contract deployed in this account's first ever transaction
});

beforeAll(async () => {
  await waitOn({resources: [hardHatNetworkEndpoint]});
  trivialApp = (await trivialAppFactory.deploy(provider.getSigner(0).getAddress())) as TrivialApp &
    Contract; // THIS MUST BE DEPLOYED FIRST IN ORDER FOR THE ABOVE ADDRESS TO BE CORRECT
  nitroAdjudicator = (await nitroAdjudicatorFactory.deploy()) as NitroAdjudicator & Contract;
  token = (await tokenFactory.deploy(provider.getSigner(0).getAddress())) as Token & Contract;

  snapshotId = await provider.send('evm_snapshot', []);
});

beforeEach(async () => {
  await provider.send('evm_revert', [snapshotId]);
  snapshotId = await provider.send('evm_snapshot', []);
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
          )} gas (${diffStr}, ${diffPercent}). Consider updating the appropriate number in gas.ts!`,
        pass: false,
      };
    }
  },
});
