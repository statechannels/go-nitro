package wallet

// Store is responsible for storing states, signatures, channel private keys and other necessary data
type Store interface {
	PrivateKey()
	SetPrivateKey()

	LockChannel()
	ReleaseChannel()

	Objective()
	SetObjective()

	// Bytecode returns the cached bytecode deployed at a given address on a given EVM chain
	Bytecode()
	SetBytecode()

	LedgerRequests()
	SetLedgerRequests()

	Funding()
	SetFunding()

	// Preimage returns some bytes that maps to a given hash under keccak256
	Preimage()
	SetPreiamge()

	NextNonce()
	// 	- concurrency control
	// - private key management
	// *- signing states
	// *- adding my signed states
	// - getting channel states
	// *- pushing messages
	//     - adding signed states
	//     - adding new objectives
	//     - collecting "dirty" channel ids
	// - getting objectives
	// - approving objectives
	// - marking objective statuses
	// - managing application bytecode
	// - managing ledger requests
	// *- adding signed states from peers
	// - creating channels, given channel constants
	// - storing data from processed blockchain events
	// - managing nonces
}
