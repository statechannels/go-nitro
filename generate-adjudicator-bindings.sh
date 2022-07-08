set -e
cd nitro-protocol

solc --base-path $(pwd) @statechannels/exit-format/=node_modules/@statechannels/exit-format/ @openzeppelin/contracts/=node_modules/@openzeppelin/contracts/ contracts/NitroAdjudicator.sol contracts/ConsensusApp.sol contracts/Token.sol --optimize --bin --abi -o tmp-build

# NitroAdjudicator
abigen --abi=$(pwd)/tmp-build/NitroAdjudicator.abi --bin=$(pwd)/tmp-build/NitroAdjudicator.bin --pkg=NitroAdjudicator --out=$(pwd)/../client/engine/chainservice/adjudicator/NitroAdjudicator.go
# ConsensusApp
abigen --abi=$(pwd)/tmp-build/ConsensusApp.abi --bin=$(pwd)/tmp-build/ConsensusApp.bin --pkg=ConsensusApp --out=$(pwd)/../client/engine/chainservice/consensusapp/ConsensusApp.go 
# Token
abigen --abi=$(pwd)/tmp-build/Token.abi --bin=$(pwd)/tmp-build/Token.bin --pkg=Token --type=Token --out=$(pwd)/../client/engine/chainservice/erc20/Token.go

rm -rf $(pwd)/tmp-build
echo "Deleted tmp-build directory."
