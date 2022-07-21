import {Contract, ContractFactory, utils} from 'ethers';
import hre from 'hardhat';

import {NitroAdjudicator} from '../typechain-types/NitroAdjudicator';
import {Token} from '../typechain-types/Token';
import {TrivialApp} from '../typechain-types/TrivialApp';
import nitroAdjudicatorArtifact from '../artifacts/contracts/NitroAdjudicator.sol/NitroAdjudicator.json';
import tokenArtifact from '../artifacts/contracts/Token.sol/Token.json';
import trivialAppArtifact from '../artifacts/contracts/TrivialApp.sol/TrivialApp.json';

export let nitroAdjudicator: NitroAdjudicator & Contract;
export let token: Token & Contract;
export let trivialApp: TrivialApp & Contract;

const provider = hre.ethers.provider;

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

export async function deployContracts() {
  trivialApp = (await trivialAppFactory.deploy(provider.getSigner(0).getAddress())) as TrivialApp &
    Contract; // THIS MUST BE DEPLOYED FIRST IN ORDER FOR THE ABOVE ADDRESS TO BE CORRECT

  nitroAdjudicator = (await nitroAdjudicatorFactory.deploy()) as NitroAdjudicator & Contract;
  token = (await tokenFactory.deploy(provider.getSigner(0).getAddress())) as Token & Contract;
}
