set -e

GONITRO_DIR=$(pwd)
NITRO_PROTOCOL_DIR=$GONITRO_DIR/packages/nitro-protocol
ARTIFACTS_DIR=$NITRO_PROTOCOL_DIR/artifacts/contracts
GETH_DIR=$(go list -m -f '{{.Dir}}' github.com/ethereum/go-ethereum)

echo "Compiling contracts..."
cd $NITRO_PROTOCOL_DIR
npx hardhat compile

echo "Generating .abi and .bin files for each contract..."

parseJson() {
  cd $ARTIFACTS_DIR
  cat ${1}.sol/${1}.json | jq -cM '.abi' > ${1}.sol/${1}.abi
  cat ${1}.sol/${1}.json | jq -r '.bytecode' > ${1}.sol/${1}.bin
}

parseJson "NitroAdjudicator"
parseJson "ConsensusApp"
parseJson "Token"
parseJson "VirtualPaymentApp"

echo "Using abigen from $GETH_DIR..."

runAbigen() {
  cd $GETH_DIR
  go run ./cmd/abigen \
    --abi=$ARTIFACTS_DIR/${1}.sol/${1}.abi \
    --bin=$ARTIFACTS_DIR/${1}.sol/${1}.bin \
    --pkg=${1} \
    --out=$GONITRO_DIR/node/engine/chainservice/${2}/${1}.go 
}

runAbigen "NitroAdjudicator" "adjudicator"
runAbigen "ConsensusApp" "consensusapp"
runAbigen "Token" "erc20"
runAbigen "VirtualPaymentApp" "virtualpaymentapp"
