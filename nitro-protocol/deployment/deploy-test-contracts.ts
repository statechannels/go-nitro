// NOTE: this script manages deploying contracts for testing purposes ONLY
// DO NOT USE THIS SCRIPT TO DEPLOY CONTRACTS TO PRODUCTION NETWORKS
import {TEST_ACCOUNTS} from '@statechannels/devtools';
import {ContractFactory, providers, Wallet} from 'ethers';

import countingAppArtifact from '../artifacts/contracts/CountingApp.sol/CountingApp.json';
import nitroAdjudicatorArtifact from '../artifacts/contracts/NitroAdjudicator.sol/NitroAdjudicator.json';
import singleAssetPaymentsArtifact from '../artifacts/contracts/examples/SingleAssetPayments.sol/SingleAssetPayments.json';
import hashLockedSwapArtifact from '../artifacts/contracts/examples/HashLockedSwap.sol/HashLockedSwap.json';
import testForceMoveArtifact from '../artifacts/contracts/test/TESTForceMove.sol/TESTForceMove.json';
import testNitroUtilsArtifact from '../artifacts/contracts/test/TESTNitroUtils.sol/TESTNitroUtils.json';
import testStrictTurnTakingArtifact from '../artifacts/contracts/test/TESTStrictTurnTaking.sol/TESTStrictTurnTaking.json';
import testConsensusArtifact from '../artifacts/contracts/test/TESTConsensus.sol/TESTConsensus.json';
import testNitroAdjudicatorArtifact from '../artifacts/contracts/test/TESTNitroAdjudicator.sol/TESTNitroAdjudicator.json';
import tokenArtifact from '../artifacts/contracts/Token.sol/Token.json';
import badTokenArtifact from '../artifacts/contracts/test/BadToken.sol/BadToken.json';
import trivialAppArtifact from '../artifacts/contracts/TrivialApp.sol/TrivialApp.json';
import consensusAppArtifact from '../artifacts/contracts/ConsensusApp.sol/ConsensusApp.json';
import virtualPaymentAppArtifact from '../artifacts/contracts/VirtualPaymentApp.sol/VirtualPaymentApp.json';
import ledgerFinancingAppArtifact from '../artifacts/contracts/LedgerFinancingApp.sol/LedgerFinancingApp.json';
const rpcEndPoint = 'http://localhost:' + process.env.GANACHE_PORT;
const provider = new providers.JsonRpcProvider(rpcEndPoint);

// Factories

const [
  countingAppFactory,
  nitroAdjudicatorFactory,
  singleAssetPaymentsFactory,
  hashLockedSwapFactory,
  testForceMoveFactory,
  testNitroUtilsFactory,
  testStrictTurnTakingFactory,
  testConsensusFactory,
  testNitroAdjudicatorFactory,
  tokenFactory,
  badTokenFactory,
  trivialAppFactory,
  consensusAppFactory,
  virtualPaymentAppFactory,
  ledgerFinancingAppFactory,
] = [
  countingAppArtifact,
  nitroAdjudicatorArtifact,
  singleAssetPaymentsArtifact,
  hashLockedSwapArtifact,
  testForceMoveArtifact,
  testNitroUtilsArtifact,
  testStrictTurnTakingArtifact,
  testConsensusArtifact,
  testNitroAdjudicatorArtifact,
  tokenArtifact,
  badTokenArtifact,
  trivialAppArtifact,
  consensusAppArtifact,
  virtualPaymentAppArtifact,
  ledgerFinancingAppArtifact,
].map(artifact =>
  new ContractFactory(artifact.abi, artifact.bytecode).connect(provider.getSigner(0))
);

export async function deploy(): Promise<Record<string, string>> {
  const NITRO_ADJUDICATOR_ADDRESS = (await nitroAdjudicatorFactory.deploy()).address;
  const COUNTING_APP_ADDRESS = (await countingAppFactory.deploy()).address;

  const HASH_LOCK_ADDRESS = (await hashLockedSwapFactory.deploy()).address;
  const SINGLE_ASSET_PAYMENTS_ADDRESS = (await singleAssetPaymentsFactory.deploy()).address;
  const CONSENSUS_APP_ADDRESS = await (await consensusAppFactory.deploy()).address;
  const VIRTUAL_PAYMENT_APP_ADDRESS = await (await virtualPaymentAppFactory.deploy()).address;
  const LEDGER_FINANCING_APP_ADDRESS = await (await ledgerFinancingAppFactory.deploy()).address;

  const TEST_NITRO_ADJUDICATOR_ADDRESS = (await testNitroAdjudicatorFactory.deploy()).address;
  const TRIVIAL_APP_ADDRESS = (await trivialAppFactory.deploy()).address;
  const TEST_FORCE_MOVE_ADDRESS = (await testForceMoveFactory.deploy()).address;
  const TEST_NITRO_UTILS_ADDRESS = (await testNitroUtilsFactory.deploy()).address;
  const TEST_STRICT_TURN_TAKING_ADDRESS = (await testStrictTurnTakingFactory.deploy()).address;
  const TEST_CONSENSUS_ADDRESS = (await testConsensusFactory.deploy()).address;

  const TEST_TOKEN_ADDRESS = (
    await tokenFactory.deploy(new Wallet(TEST_ACCOUNTS[0].privateKey).address)
  ).address;

  const BAD_TOKEN_ADDRESS = (
    await badTokenFactory.deploy(new Wallet(TEST_ACCOUNTS[0].privateKey).address)
  ).address;

  return {
    NITRO_ADJUDICATOR_ADDRESS,
    COUNTING_APP_ADDRESS,
    HASH_LOCK_ADDRESS,
    SINGLE_ASSET_PAYMENTS_ADDRESS,
    TRIVIAL_APP_ADDRESS,
    CONSENSUS_APP_ADDRESS,
    VIRTUAL_PAYMENT_APP_ADDRESS,
    LEDGER_FINANCING_APP_ADDRESS,
    TEST_FORCE_MOVE_ADDRESS,
    TEST_NITRO_UTILS_ADDRESS,
    TEST_STRICT_TURN_TAKING_ADDRESS,
    TEST_CONSENSUS_ADDRESS,
    TEST_NITRO_ADJUDICATOR_ADDRESS,
    TEST_TOKEN_ADDRESS,
    BAD_TOKEN_ADDRESS,
  };
}
