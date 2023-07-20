set -e

trap '{
rm -rf $(pwd)/tmp-build
echo "Deleted tmp-build directory."
 }' EXIT


cd nitro-protocol

solc --base-path $(pwd) \
  @statechannels/exit-format/=node_modules/@statechannels/exit-format/ \
  @openzeppelin/contracts/=node_modules/@openzeppelin/contracts/ \
  contracts/NitroAdjudicator.sol contracts/ConsensusApp.sol contracts/Token.sol contracts/VirtualPaymentApp.sol contracts/deploy/Create2Deployer.sol \
  --optimize --bin --abi -o tmp-build --via-ir

runAbigen() {
  abigen --abi=$(pwd)/tmp-build/${1}.abi \
    --bin=$(pwd)/tmp-build/${1}.bin \
    --pkg=${1} \
    --out=$(pwd)/../node/engine/chainservice/${2}/${1}.go 
}

runAbigenV2() {
  abigen --v2 --abi=$(pwd)/tmp-build/${1}.abi \
    --bin=$(pwd)/tmp-build/${1}.bin \
    --pkg=${1} \
    --out=$(pwd)/../node/engine/chainservice/${2}/${1}-v2.go 
}

runAbigen "NitroAdjudicator" "adjudicator"
runAbigenV2 "NitroAdjudicator" "adjudicatorv2"
runAbigen "ConsensusApp" "consensusapp"
runAbigen "Token" "erc20" # TODO: getting an error generating the v2 bindings for this one, not sure why
runAbigen "VirtualPaymentApp" "virtualpaymentapp"


