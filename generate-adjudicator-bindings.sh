set -e
cd nitro-protocol

solc --base-path $(pwd) 
  @statechannels/exit-format/=node_modules/@statechannels/exit-format/ \
  @openzeppelin/contracts/=node_modules/@openzeppelin/contracts/ \
  contracts/NitroAdjudicator.sol contracts/ConsensusApp.sol contracts/Token.sol contracts/VirtualPaymentApp.sol contracts/deploy/Create2Deployer.sol \
  --optimize --bin --abi -o tmp-build --via-ir

function runAbigen {
  abigen --abi=$(pwd)/tmp-build/${1}.abi \
    --bin=$(pwd)/tmp-build/${1}.bin \
    --pkg=${1} \
    --out=$(pwd)/../client/engine/chainservice/${2}/${1}.go 

}

runAbigen "NitroAdjudicator" "adjudicator"
runAbigen "ConsensusApp" "consensusapp"
runAbigen "Token" "erc20"
runAbigen "VirtualPaymentApp" "virtualpaymentapp"
runAbigen "Create2Deployer" "create2deployer"

rm -rf $(pwd)/tmp-build
echo "Deleted tmp-build directory."