import { ethers } from "hardhat";

async function main() {
  let privateKey = '0x9a14bf0eb618a3407a12a83a74dfe7bbed098ccc6347985b92ab08e81996cfc9';
  let provider = new ethers.providers.JsonRpcProvider();
  let deployer = new ethers.Wallet(privateKey, provider);
  
  console.log("Deploying contracts with the account:", deployer.address);
  console.log("Account balance:", (await deployer.getBalance()).toString());
  
  const NAFactory = await ethers.getContractFactory("NitroAdjudicator");
  const NitroAdjudicator = await NAFactory.deploy();

  await NitroAdjudicator.deployed();

  console.log("NitroAdjudicator address:", NitroAdjudicator.address);
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
