package types

import (
	abci "github.com/strangelove-ventures/cometbft-client/abci/types"
)

// ABCIResults wraps the deliver tx results to return a proof.
type ABCIResults []*abci.ExecTxResult
