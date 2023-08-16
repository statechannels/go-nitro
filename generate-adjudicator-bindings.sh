set -e

GONITRO_DIR=$(pwd)
NITRO_PROTOCOL_DIR=$GONITRO_DIR/nitro-protocol
TEMP_DIR=$NITRO_PROTOCOL_DIR/tmp-build
GETH_DIR=$(go list -m -f '{{.Dir}}' github.com/ethereum/go-ethereum)



trap '{
rm -rf $TEMP_DIR
echo "Deleted tmp-build directory."
 }' EXIT

echo "Compiling contracts..."

solc --base-path $NITRO_PROTOCOL_DIR \
  @statechannels/exit-format/=node_modules/@statechannels/exit-format/ \
  @openzeppelin/contracts/=node_modules/@openzeppelin/contracts/ \
  $NITRO_PROTOCOL_DIR/contracts/NitroAdjudicator.sol \
  $NITRO_PROTOCOL_DIR/contracts/ConsensusApp.sol \
  $NITRO_PROTOCOL_DIR/contracts/Token.sol \
  $NITRO_PROTOCOL_DIR/contracts/VirtualPaymentApp.sol \
  $NITRO_PROTOCOL_DIR/contracts/deploy/Create2Deployer.sol \
  --optimize --bin --abi -o $TEMP_DIR --via-ir

runAbigen() {
  cd $GETH_DIR
  go run ./cmd/abigen \
    --abi=$TEMP_DIR/${1}.abi \
    --bin=$TEMP_DIR/${1}.bin \
    --pkg=${1} \
    --out=$GONITRO_DIR/node/engine/chainservice/${2}/${1}.go 
}

echo "Using abigen from $GETH_DIR..."

runAbigen "NitroAdjudicator" "adjudicator"
runAbigen "ConsensusApp" "consensusapp"
runAbigen "Token" "erc20"
runAbigen "VirtualPaymentApp" "virtualpaymentapp"