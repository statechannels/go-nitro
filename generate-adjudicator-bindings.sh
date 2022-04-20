set -e
cd nitro-protocol
solc --base-path $(pwd) @statechannels/exit-format/=node_modules/@statechannels/exit-format/ @openzeppelin/contracts/=node_modules/@openzeppelin/contracts/ contracts/NitroAdjudicator.sol --optimize --bin --abi -o tmp-build
abigen --abi=$(pwd)/tmp-build/NitroAdjudicator.abi --bin=$(pwd)/tmp-build/NitroAdjudicator.bin --pkg=NitroAdjudicator --out=$(pwd)/../client/engine/chainservice/adjudicator/NitroAdjudicator.go
rm -rf $(pwd)/tmp-build
echo "Deleted tmp-build directory."
