package client

import (
	"context"
	"encoding/base64"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/strangelove-ventures/cometbft-client/abci/types"
	_ "github.com/strangelove-ventures/cometbft-client/crypto/encoding"
	"github.com/strangelove-ventures/cometbft-client/libs/bytes"
	rpcclient "github.com/strangelove-ventures/cometbft-client/rpc/client"
	rpchttp "github.com/strangelove-ventures/cometbft-client/rpc/client/http"
	coretypes "github.com/strangelove-ventures/cometbft-client/rpc/core/types"
	jsonrpc "github.com/strangelove-ventures/cometbft-client/rpc/jsonrpc/client"
	"github.com/strangelove-ventures/cometbft-client/types"
)

// Client is a wrapper around the CometBFT RPC client.
type Client struct {
	rpcClient rpcclient.Client
}

// NewClient returns a pointer to a new instance of Client.
func NewClient(addr string, timeout time.Duration) (*Client, error) {
	rpcClient, err := newRPCClient(addr, timeout)
	if err != nil {
		return nil, err
	}

	return &Client{rpcClient}, nil
}

// BlockResults fetches the block results at a specific height,
// it then parses the tx results and block events into our generalized types.
// This allows us to maintain backwards compatability with older versions of CometBFT.
func (c *Client) BlockResults(ctx context.Context, height *int64) (*BlockResponse, error) {
	res, err := c.rpcClient.BlockResults(ctx, height)
	if err != nil {
		return nil, err
	}

	var txRes []*ExecTxResponse
	for _, tx := range res.TxsResults {
		txRes = append(txRes, &ExecTxResponse{
			Code:      tx.Code,
			Data:      tx.Data,
			Log:       tx.Log,
			Info:      tx.Info,
			GasWanted: tx.GasWanted,
			GasUsed:   tx.GasUsed,
			Events:    parseEvents(tx.Events),
			Codespace: tx.Codespace,
		})
	}

	if res.FinalizeBlockEvents != nil && len(res.FinalizeBlockEvents) > 0 {
		return &BlockResponse{
			Height:           res.Height,
			TxResponses:      txRes,
			Events:           parseEvents(res.FinalizeBlockEvents),
			ValidatorUpdates: res.ValidatorUpdates,
			AppHash:          res.AppHash,
		}, nil
	}

	events := res.BeginBlockEvents
	events = append(events, res.EndBlockEvents...)

	return &BlockResponse{
		Height:           res.Height,
		TxResponses:      txRes,
		Events:           parseEvents(events),
		ValidatorUpdates: res.ValidatorUpdates,
		AppHash:          res.AppHash,
	}, nil
}

func (c *Client) Tx(ctx context.Context, hash []byte, prove bool) (*TxResponse, error) {
	res, err := c.rpcClient.Tx(ctx, hash, prove)
	if err != nil {
		return nil, err
	}

	execTx := ExecTxResponse{
		Code:      res.TxResult.Code,
		Data:      res.TxResult.Data,
		Log:       res.TxResult.Log,
		Info:      res.TxResult.Info,
		GasWanted: res.TxResult.GasWanted,
		GasUsed:   res.TxResult.GasUsed,
		Events:    parseEvents(res.TxResult.Events),
		Codespace: res.TxResult.Codespace,
	}

	return &TxResponse{
		Hash:   res.Hash,
		Height: res.Height,
		Index:  res.Index,
		ExecTx: execTx,
		Tx:     res.Tx,
		Proof:  res.Proof,
	}, nil
}

func (c *Client) TxSearch(
	ctx context.Context,
	query string,
	prove bool,
	page *int,
	perPage *int,
	orderBy string,
) ([]*TxResponse, error) {
	res, err := c.rpcClient.TxSearch(ctx, query, prove, page, perPage, orderBy)
	if err != nil {
		return nil, err
	}

	result := make([]*TxResponse, len(res.Txs))

	for i, tx := range res.Txs {
		execTx := ExecTxResponse{
			Code:      tx.TxResult.Code,
			Data:      tx.TxResult.Data,
			Log:       tx.TxResult.Log,
			Info:      tx.TxResult.Info,
			GasWanted: tx.TxResult.GasWanted,
			GasUsed:   tx.TxResult.GasUsed,
			Events:    parseEvents(tx.TxResult.Events),
			Codespace: tx.TxResult.Codespace,
		}

		result[i] = &TxResponse{
			Hash:   tx.Hash,
			Height: tx.Height,
			Index:  tx.Index,
			ExecTx: execTx,
			Tx:     tx.Tx,
			Proof:  tx.Proof,
		}
	}

	return result, nil
}

func (c *Client) Commit(ctx context.Context, height *int64) (*coretypes.ResultCommit, error) {
	res, err := c.rpcClient.Commit(ctx, height)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) Validators(
	ctx context.Context,
	height *int64,
	page *int,
	perPage *int,
) (*coretypes.ResultValidators, error) {
	res, err := c.rpcClient.Validators(ctx, height, page, perPage)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) Status(ctx context.Context) (*coretypes.ResultStatus, error) {
	res, err := c.rpcClient.Status(ctx)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) Block(ctx context.Context, height *int64) (*coretypes.ResultBlock, error) {
	res, err := c.rpcClient.Block(ctx, height)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) BlockSearch(
	ctx context.Context,
	query string,
	page *int,
	perPage *int,
	orderBy string,
) (*coretypes.ResultBlockSearch, error) {
	res, err := c.rpcClient.BlockSearch(ctx, query, page, perPage, orderBy)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) BlockByHash(ctx context.Context, hash []byte) (*coretypes.ResultBlock, error) {
	res, err := c.rpcClient.BlockByHash(ctx, hash)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) BlockchainInfo(ctx context.Context, minHeight int64, maxHeight int64) (*coretypes.ResultBlockchainInfo, error) {
	res, err := c.rpcClient.BlockchainInfo(ctx, minHeight, maxHeight)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) BroadcastTxAsync(ctx context.Context, tx types.Tx) (*coretypes.ResultBroadcastTx, error) {
	res, err := c.rpcClient.BroadcastTxAsync(ctx, tx)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) BroadcastTxSync(ctx context.Context, tx types.Tx) (*coretypes.ResultBroadcastTx, error) {
	res, err := c.rpcClient.BroadcastTxSync(ctx, tx)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) BroadcastTxCommit(ctx context.Context, tx types.Tx) (*coretypes.ResultBroadcastTxCommit, error) {
	res, err := c.rpcClient.BroadcastTxCommit(ctx, tx)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) ABCIInfo(ctx context.Context) (*coretypes.ResultABCIInfo, error) {
	res, err := c.rpcClient.ABCIInfo(ctx)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) ABCIQuery(ctx context.Context, path string, data bytes.HexBytes) (*coretypes.ResultABCIQuery, error) {
	res, err := c.rpcClient.ABCIQuery(ctx, path, data)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) ABCIQueryWithOptions(
	ctx context.Context,
	path string,
	data bytes.HexBytes,
	opts rpcclient.ABCIQueryOptions,
) (*coretypes.ResultABCIQuery, error) {
	res, err := c.rpcClient.ABCIQueryWithOptions(ctx, path, data, opts)
	if err != nil {
		return nil, err
	}

	return res, nil
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

// parseEvents returns a slice of sdk.StringEvent objects that are composed from a slice of abci.Event objects.
// parseEvents will first attempt to base64 decode the abci.Event objects and if an error is encountered it will
// fall back to the stringifyEvents function.
func parseEvents(events []abci.Event) sdk.StringEvents {
	decodedEvents, err := base64DecodeEvents(events)
	if err == nil {
		return decodedEvents
	}

	return stringifyEvents(events)
}

// base64DecodeEvents attempts to base64 decode a slice of Event objects.
// An error is returned if base64 decoding any event in the slice fails.
func base64DecodeEvents(events []abci.Event) (sdk.StringEvents, error) {
	sdkEvents := make(sdk.StringEvents, len(events))

	for i, event := range events {
		evt := sdk.StringEvent{Type: event.Type}

		for _, attr := range event.Attributes {
			key, err := base64.StdEncoding.DecodeString(attr.Key)
			if err != nil {
				return nil, err
			}

			value, err := base64.StdEncoding.DecodeString(attr.Value)
			if err != nil {
				return nil, err
			}

			evt.Attributes = append(evt.Attributes, sdk.Attribute{
				Key:   string(key),
				Value: string(value),
			})
		}

		sdkEvents[i] = evt
	}

	return sdkEvents, nil
}

// stringifyEvents converts a slice of Event objects into a slice of StringEvent objects.
// This function is copied straight from the Cosmos SDK, so we can alter it to handle our abci.Event type.
func stringifyEvents(events []abci.Event) sdk.StringEvents {
	res := make(sdk.StringEvents, 0, len(events))

	for _, e := range events {
		res = append(res, stringifyEvent(e))
	}

	return res
}

// stringifyEvent converts an Event object to a StringEvent object.
// This function is copied straight from the Cosmos SDK, so we can alter it to handle our abci.Event type.
func stringifyEvent(e abci.Event) sdk.StringEvent {
	res := sdk.StringEvent{Type: e.Type}

	for _, attr := range e.Attributes {
		res.Attributes = append(
			res.Attributes,
			sdk.Attribute{Key: attr.Key, Value: attr.Value},
		)
	}

	return res
}
