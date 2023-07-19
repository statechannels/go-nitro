set -e
cd nitro-protocol

solc --base-path $(pwd) \
  @statechannels/exit-format/=node_modules/@statechannels/exit-format/ \
  @openzeppelin/contracts/=node_modules/@openzeppelin/contracts/ \
  contracts/NitroAdjudicator.sol contracts/ConsensusApp.sol contracts/Token.sol contracts/VirtualPaymentApp.sol contracts/deploy/Create2Deployer.sol \
  --optimize --bin --abi -o tmp-build --via-ir

runAbigen() {
  abigen --v2 --abi=$(pwd)/tmp-build/${1}.abi \
    --bin=$(pwd)/tmp-build/${1}.bin \
    --pkg=${1} \
    --out=$(pwd)/../node/engine/chainservice/${2}/${1}.go 
}

runAbigen "NitroAdjudicator" "adjudicator"
runAbigen "ConsensusApp" "consensusapp"
# runAbigen "Token" "erc20" # TODO: getting an error generating these bindings for this one, not sure why
runAbigen "VirtualPaymentApp" "virtualpaymentapp"

rm -rf $(pwd)/tmp-build
echo "Deleted tmp-build directory."
