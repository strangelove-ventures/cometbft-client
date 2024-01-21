package client

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/strangelove-ventures/cometbft-client/abci/types"
	"github.com/strangelove-ventures/cometbft-client/libs/bytes"
	"github.com/strangelove-ventures/cometbft-client/types"
)

// BlockResponse is used in place of the CometBFT type ResultBlockResults.
// This allows us to handle the decoding of events internally so that we can return events to consumers as raw strings.
type BlockResponse struct {
	Height           int64
	TxResponses      []*ExecTxResponse
	Events           sdk.StringEvents
	ValidatorUpdates []abci.ValidatorUpdate
	AppHash          []byte
}

// ExecTxResponse is used in place of the CometBFT type ExecTxResult.
// This allows us to handle the decoding of events internally so that we can return events to consumers as raw strings.
type ExecTxResponse struct {
	Code      uint32
	Data      []byte
	Log       string
	Info      string
	GasWanted int64
	GasUsed   int64
	Events    sdk.StringEvents
	Codespace string
}

func (e *ExecTxResponse) IsOK() bool {
	return e.Code == abci.CodeTypeOK
}

// TxResponse is used in place of the CometBFT type ResultTx.
// This allows us to handle the decoding of events internally so that we can return events to consumers as raw strings.
type TxResponse struct {
	Hash   bytes.HexBytes
	Height int64
	Index  uint32
	ExecTx ExecTxResponse
	Tx     types.Tx
	Proof  types.TxProof
}
