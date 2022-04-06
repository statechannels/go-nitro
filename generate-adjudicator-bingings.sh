cd nitro-protocol
solc @statechannels/exit-format/=$(pwd)/node_modules/@statechannels/exit-format/ @openzeppelin/contracts/=$(pwd)/node_modules/@openzeppelin/contracts/ $(pwd)/contracts/NitroAdjudicator.sol --optimize --bin --abi -o tmp-build
abigen --abi=$(pwd)/tmp-build/NitroAdjudicator.abi --bin=$(pwd)/tmp-build/NitroAdjudicator.bin --pkg=NitroAdjudicator --out=$(pwd)/../client/engine/chainservice/adjudicator/NitroAdjudicator.go
rm -rf $(pwd)/tmp-build