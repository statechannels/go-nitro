import { Interface } from 'ethers/lib/utils';
import { ethers } from 'hardhat';

async function main() {
  let privateKey = '0x9a14bf0eb618a3407a12a83a74dfe7bbed098ccc6347985b92ab08e81996cfc9';
  let provider = new ethers.providers.JsonRpcProvider();
  let wallet = new ethers.Wallet(privateKey, provider);
  
  console.log("Pinging contracts with the account:", wallet.address);
  console.log("Account balance:", (await wallet.getBalance()).toString());

  let abi = new Interface([
    "function getNumber() public pure returns(uint256)"
  ]);
  let NitroAdjudicatorAddress = '0xCc388ae2496E15ff8C6df70566171c750B5118E2';
  let NitroAdjudicator = new ethers.Contract(NitroAdjudicatorAddress, abi, wallet);

  console.log("The number is:", (await NitroAdjudicator.getNumber()).toNumber());
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
