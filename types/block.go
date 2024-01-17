package types

import (
	"time"

	"github.com/strangelove-ventures/cometbft-client/libs/bits"
	cmtbytes "github.com/strangelove-ventures/cometbft-client/libs/bytes"
	cmtsync "github.com/strangelove-ventures/cometbft-client/libs/sync"
)

// Block defines the atomic unit of a CometBFT blockchain.
type Block struct {
	mtx cmtsync.Mutex

	Header     `json:"header"`
	Data       `json:"data"`
	Evidence   EvidenceData `json:"evidence"`
	LastCommit *Commit      `json:"last_commit"`
}

//-----------------------------------------------------------------------------

// Header defines the structure of a CometBFT block header.
// NOTE: changes to the Header should be duplicated in:
// - header.Hash()
// - abci.Header
// - https://github.com/cometbft/cometbft/blob/v0.38.x/spec/blockchain/blockchain.md
type Header struct {
	// basic block info
	//Version cmtversion.Consensus `json:"version"`
	ChainID string    `json:"chain_id"`
	Height  int64     `json:"height"`
	Time    time.Time `json:"time"`

	// prev block info
	LastBlockID BlockID `json:"last_block_id"`

	// hashes of block data
	LastCommitHash cmtbytes.HexBytes `json:"last_commit_hash"` // commit from validators from the last block
	DataHash       cmtbytes.HexBytes `json:"data_hash"`        // transactions

	// hashes from the app output from the prev block
	ValidatorsHash     cmtbytes.HexBytes `json:"validators_hash"`      // validators for the current block
	NextValidatorsHash cmtbytes.HexBytes `json:"next_validators_hash"` // validators for the next block
	ConsensusHash      cmtbytes.HexBytes `json:"consensus_hash"`       // consensus params for current block
	AppHash            cmtbytes.HexBytes `json:"app_hash"`             // state after txs from the previous block
	// root hash of all results from the txs from the previous block
	// see `deterministicExecTxResult` to understand which parts of a tx is hashed into here
	LastResultsHash cmtbytes.HexBytes `json:"last_results_hash"`

	// consensus info
	EvidenceHash    cmtbytes.HexBytes `json:"evidence_hash"`    // evidence included in the block
	ProposerAddress Address           `json:"proposer_address"` // original proposer of the block
}

//-------------------------------------

// BlockIDFlag indicates which BlockID the signature is for.
type BlockIDFlag byte

// CommitSig is a part of the Vote included in a Commit.
type CommitSig struct {
	BlockIDFlag      BlockIDFlag `json:"block_id_flag"`
	ValidatorAddress Address     `json:"validator_address"`
	Timestamp        time.Time   `json:"timestamp"`
	Signature        []byte      `json:"signature"`
}

//-------------------------------------

// ExtendedCommitSig contains a commit signature along with its corresponding
// vote extension and vote extension signature.
type ExtendedCommitSig struct {
	CommitSig                 // Commit signature
	Extension          []byte // Vote extension
	ExtensionSignature []byte // Vote extension signature
}

//-------------------------------------

// Commit contains the evidence that a block was committed by a set of validators.
// NOTE: Commit is empty for height 1, but never nil.
type Commit struct {
	// NOTE: The signatures are in order of address to preserve the bonded
	// ValidatorSet order.
	// Any peer with a block can gossip signatures by index with a peer without
	// recalculating the active ValidatorSet.
	Height     int64       `json:"height"`
	Round      int32       `json:"round"`
	BlockID    BlockID     `json:"block_id"`
	Signatures []CommitSig `json:"signatures"`

	// Memoized in first call to corresponding method.
	// NOTE: can't memoize in constructor because constructor isn't used for
	// unmarshaling.
	hash cmtbytes.HexBytes
}

//-------------------------------------

// ExtendedCommit is similar to Commit, except that its signatures also retain
// their corresponding vote extensions and vote extension signatures.
type ExtendedCommit struct {
	Height             int64
	Round              int32
	BlockID            BlockID
	ExtendedSignatures []ExtendedCommitSig

	bitArray *bits.BitArray
}

//-------------------------------------

// Data contains the set of transactions included in the block
type Data struct {
	// Txs that will be applied by state @ block.Height+1.
	// NOTE: not all txs here are valid.  We're just agreeing on the order first.
	// This means that block.AppHash does not include these txs.
	Txs Txs `json:"txs"`

	// Volatile
	hash cmtbytes.HexBytes
}

//-----------------------------------------------------------------------------

// EvidenceData contains any evidence of malicious wrong-doing by validators
type EvidenceData struct {
	Evidence EvidenceList `json:"evidence"`

	// Volatile. Used as cache
	hash     cmtbytes.HexBytes
	byteSize int64
}

//--------------------------------------------------------------------------------

// BlockID
type BlockID struct {
	Hash          cmtbytes.HexBytes `json:"hash"`
	PartSetHeader PartSetHeader     `json:"parts"`
}
