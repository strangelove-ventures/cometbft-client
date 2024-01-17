package cometbft_client

import (
	"context"
	"time"

	rpcclient "github.com/strangelove-ventures/cometbft-client/rpc/client"
	rpchttp "github.com/strangelove-ventures/cometbft-client/rpc/client/http"
	jsonrpc "github.com/strangelove-ventures/cometbft-client/rpc/jsonrpc/client"
	"golang.org/x/mod/semver"
)

const cometBlockResultsThreshold = "v0.38.0-alpha"

// Client is a wrapper around the CometBFT RPC client.
type Client struct {
	RPCClient rpcclient.Client
}

// NewClient returns a pointer to a new instance of Client.
func NewClient(addr string, timeout time.Duration) (*Client, error) {
	rpcClient, err := newRPCClient(addr, timeout)
	if err != nil {
		return nil, err
	}

	return &Client{RPCClient: rpcClient}, nil
}

// Init attempts to initialize the CometBFT RPC client.
// Init MUST be called after instantiating a new instance of Client before it can be used.
func (c *Client) Init(addr, version string, timeout time.Duration) error {

	return nil
}

// BlockResults uses the appropriate CometBFT RPC client to fetch the block results at a specific height,
// it then parses the tx results and block events into our generalized types.
func (c *Client) BlockResults(ctx context.Context, height *int64) (*Results, error) {
	var results *Results

	// switch {
	// case c.cometLegacyBlockResults:
	// TODO: finish implementation for legacy RPC Client call to BlockResults
	//legacyRes, err := c.LegacyRPCClient.BlockResults(ctx, height)
	//if err != nil {
	//	return nil, err
	//}
	//
	//var events []abci.Event
	//events = append(events, ConvertEvents(legacyRes.BeginBlockEvents)...)
	//events = append(events, ConvertEvents(legacyRes.EndBlockEvents)...)
	//
	//results = &Results{
	//	TxsResults: ConvertTxResults(legacyRes.TxsResults),
	//	Events:     events,
	//}
	//default:
	res, err := c.RPCClient.BlockResults(ctx, height)
	if err != nil {
		return nil, err
	}

	var txRes []*TxResult
	for _, tx := range res.TxsResults {
		txRes = append(txRes, &TxResult{
			Code:   tx.Code,
			Events: tx.Events,
		})
	}

	results = &Results{
		TxsResults: txRes,
		Events:     res.FinalizeBlockEvents,
	}
	//}

	return results, nil
}

func newRPCClient(addr string, timeout time.Duration) (*rpchttp.HTTP, error) {
	httpClient, err := jsonrpc.DefaultHTTPClient(addr)
	if err != nil {
		return nil, err
	}

	httpClient.Timeout = timeout

	rpcClient, err := rpchttp.NewWithClient(addr, "/websocket", httpClient)
	if err != nil {
		return nil, err
	}

	return rpcClient, nil
}

func (c *Client) legacyBlockResults(version string) bool {
	return semver.Compare("v"+version, cometBlockResultsThreshold) < 0
}
