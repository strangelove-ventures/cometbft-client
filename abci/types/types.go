package types

import "github.com/strangelove-ventures/cometbft-client/proto/tendermint/crypto"

const (
	CodeTypeOK uint32 = 0
)

// ValidatorUpdates is a list of validators that implements the Sort interface
type ValidatorUpdates []ValidatorUpdate

type ValidatorUpdate struct {
	PubKey crypto.PublicKey `protobuf:"bytes,1,opt,name=pub_key,json=pubKey,proto3" json:"pub_key"`
	Power  int64            `protobuf:"varint,2,opt,name=power,proto3" json:"power,omitempty"`
}

type ExecTxResult struct {
	Code      uint32  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Data      []byte  `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	Log       string  `protobuf:"bytes,3,opt,name=log,proto3" json:"log,omitempty"`
	Info      string  `protobuf:"bytes,4,opt,name=info,proto3" json:"info,omitempty"`
	GasWanted int64   `protobuf:"varint,5,opt,name=gas_wanted,proto3" json:"gas_wanted,omitempty"`
	GasUsed   int64   `protobuf:"varint,6,opt,name=gas_used,proto3" json:"gas_used,omitempty"`
	Events    []Event `protobuf:"bytes,7,rep,name=events,proto3" json:"events,omitempty"`
	Codespace string  `protobuf:"bytes,8,opt,name=codespace,proto3" json:"codespace,omitempty"`
}

// IsOK returns true if Code is OK.
func (r ExecTxResult) IsOK() bool {
	return r.Code == CodeTypeOK
}

// IsErr returns true if Code is something other than OK.
func (r ExecTxResult) IsErr() bool {
	return r.Code != CodeTypeOK
}

// -----------------------------------------------
// construct Result data

// Event allows application developers to attach additional information to
// ResponseFinalizeBlock and ResponseCheckTx.
// Later, transactions may be queried using these events.
type Event struct {
	Type       string           `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Attributes []EventAttribute `protobuf:"bytes,2,rep,name=attributes,proto3" json:"attributes,omitempty"`
}

// EventAttribute is a single key-value pair, associated with an event.
type EventAttribute struct {
	Key   string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	Index bool   `protobuf:"varint,3,opt,name=index,proto3" json:"index,omitempty"`
}

type ResponseInfo struct {
	Data             string `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	Version          string `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	AppVersion       uint64 `protobuf:"varint,3,opt,name=app_version,json=appVersion,proto3" json:"app_version,omitempty"`
	LastBlockHeight  int64  `protobuf:"varint,4,opt,name=last_block_height,json=lastBlockHeight,proto3" json:"last_block_height,omitempty"`
	LastBlockAppHash []byte `protobuf:"bytes,5,opt,name=last_block_app_hash,json=lastBlockAppHash,proto3" json:"last_block_app_hash,omitempty"`
}

type ResponseQuery struct {
	Code uint32 `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	// bytes data = 2; // use "value" instead.
	Log       string    `protobuf:"bytes,3,opt,name=log,proto3" json:"log,omitempty"`
	Info      string    `protobuf:"bytes,4,opt,name=info,proto3" json:"info,omitempty"`
	Index     int64     `protobuf:"varint,5,opt,name=index,proto3" json:"index,omitempty"`
	Key       []byte    `protobuf:"bytes,6,opt,name=key,proto3" json:"key,omitempty"`
	Value     []byte    `protobuf:"bytes,7,opt,name=value,proto3" json:"value,omitempty"`
	ProofOps  *ProofOps `protobuf:"bytes,8,opt,name=proof_ops,json=proofOps,proto3" json:"proofOps,omitempty"`
	Height    int64     `protobuf:"varint,9,opt,name=height,proto3" json:"height,omitempty"`
	Codespace string    `protobuf:"bytes,10,opt,name=codespace,proto3" json:"codespace,omitempty"`
}

func (r *ResponseQuery) IsOK() bool {
	return r.Code == CodeTypeOK
}

// ProofOps is Merkle proof defined by the list of ProofOps
type ProofOps struct {
	Ops []ProofOp `protobuf:"bytes,1,rep,name=ops,proto3" json:"ops"`
}

// ProofOp defines an operation used for calculating Merkle root
// The data could be arbitrary format, providing nessecary data
// for example neighbouring node hash
type ProofOp struct {
	Type string `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Key  []byte `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	Data []byte `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
}

type ResponseCheckTx struct {
	Code      uint32  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Data      []byte  `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	Log       string  `protobuf:"bytes,3,opt,name=log,proto3" json:"log,omitempty"`
	Info      string  `protobuf:"bytes,4,opt,name=info,proto3" json:"info,omitempty"`
	GasWanted int64   `protobuf:"varint,5,opt,name=gas_wanted,proto3" json:"gas_wanted,omitempty"`
	GasUsed   int64   `protobuf:"varint,6,opt,name=gas_used,proto3" json:"gas_used,omitempty"`
	Events    []Event `protobuf:"bytes,7,rep,name=events,proto3" json:"events,omitempty"`
	Codespace string  `protobuf:"bytes,8,opt,name=codespace,proto3" json:"codespace,omitempty"`
}
