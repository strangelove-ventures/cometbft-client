package client

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Results is a generalized type used to handle the breaking changes in the CometBFT type,
// ResultBlockResults, that were introduced in v0.38.0.
type Results struct {
	TxsResults []*TxResult
	Events     sdk.StringEvents
}

// TxResult is a generalized type used to handle the breaking changes in the CometBFT type,
// ResultBlockResults, that were introduced in v0.38.0.
type TxResult struct {
	Code   uint32
	Events sdk.StringEvents
}
