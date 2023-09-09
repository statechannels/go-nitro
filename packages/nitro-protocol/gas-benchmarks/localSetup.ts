import {Contract, ContractFactory, utils} from 'ethers';
import hre from 'hardhat';

import '@nomiclabs/hardhat-ethers';
import {NitroAdjudicator} from '../typechain-types/NitroAdjudicator';
import {BatchOperator} from '../typechain-types/BatchOperator';
import {Token} from '../typechain-types/Token';
import {VirtualPaymentApp} from '../typechain-types/VirtualPaymentApp';
import {ConsensusApp} from '../typechain-types/ConsensusApp';
import nitroAdjudicatorArtifact from '../artifacts/contracts/NitroAdjudicator.sol/NitroAdjudicator.json';
import batchOperatorArtifact from '../artifacts/contracts/auxiliary/BatchOperator.sol/BatchOperator.json';
import tokenArtifact from '../artifacts/contracts/Token.sol/Token.json';
import consensusAppArtifact from '../artifacts/contracts/ConsensusApp.sol/ConsensusApp.json';
import virtualPaymentAppArtifact from '../artifacts/contracts/VirtualPaymentApp.sol/VirtualPaymentApp.json';

export let nitroAdjudicator: NitroAdjudicator & Contract;
export let batchOperator: BatchOperator & Contract;
export let token: Token & Contract;
export let consensusApp: ConsensusApp & Contract;
export let virtualPaymentApp: VirtualPaymentApp & Contract;

const provider = hre.ethers.provider;

const tokenFactory = new ContractFactory(tokenArtifact.abi, tokenArtifact.bytecode).connect(
  provider.getSigner(0)
);

const nitroAdjudicatorFactory = new ContractFactory(
  nitroAdjudicatorArtifact.abi,
  nitroAdjudicatorArtifact.bytecode
).connect(provider.getSigner(0));

const batchOperatorFactory = new ContractFactory(
  batchOperatorArtifact.abi,
  batchOperatorArtifact.bytecode
).connect(provider.getSigner(0));

const consensusAppFactory = new ContractFactory(
  consensusAppArtifact.abi,
  consensusAppArtifact.bytecode
).connect(provider.getSigner(0));

const virtualPaymentAppFactory = new ContractFactory(
  virtualPaymentAppArtifact.abi,
  virtualPaymentAppArtifact.bytecode
).connect(provider.getSigner(0));

export const consensusAppAddress = utils.getContractAddress({
  from: '0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266', // ASSUME: deployed by hardhat account 0
  nonce: 0, // ASSUME: this contract deployed in this account's first ever transaction
});

export const virtualPaymentAppAddress = utils.getContractAddress({
  from: '0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266', // ASSUME: deployed by hardhat account 0
  nonce: 1, // ASSUME: this contract deployed in this account's second ever transaction
});

export async function deployContracts() {
  consensusApp = (await consensusAppFactory.deploy(
    provider.getSigner(0).getAddress()
  )) as ConsensusApp & Contract; // THIS MUST BE DEPLOYED FIRST IN ORDER FOR THE ABOVE ADDRESS TO BE CORRECT
  virtualPaymentApp = (await virtualPaymentAppFactory.deploy(
    provider.getSigner(0).getAddress()
  )) as VirtualPaymentApp & Contract; // THIS MUST BE DEPLOYED SECOND IN ORDER FOR THE ABOVE ADDRESS TO BE CORRECT

  nitroAdjudicator = (await nitroAdjudicatorFactory.deploy()) as NitroAdjudicator & Contract;
  batchOperator = (await batchOperatorFactory.deploy(nitroAdjudicator.address)) as BatchOperator &
    Contract;
  token = (await tokenFactory.deploy(provider.getSigner(0).getAddress())) as Token & Contract;
}
