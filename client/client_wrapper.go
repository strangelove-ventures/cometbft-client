package client

import (
	"context"
	"encoding/base64"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/strangelove-ventures/cometbft-client/abci/types"
	_ "github.com/strangelove-ventures/cometbft-client/crypto/encoding"
	rpcclient "github.com/strangelove-ventures/cometbft-client/rpc/client"
	rpchttp "github.com/strangelove-ventures/cometbft-client/rpc/client/http"
	jsonrpc "github.com/strangelove-ventures/cometbft-client/rpc/jsonrpc/client"
)

// Client is a wrapper around the CometBFT RPC client.
type Client struct {
	rpcclient.Client
}

// NewClient returns a pointer to a new instance of Client.
func NewClient(addr string, timeout time.Duration) (*Client, error) {
	rpcClient, err := newRPCClient(addr, timeout)
	if err != nil {
		return nil, err
	}

	return &Client{rpcClient}, nil
}

// GetBlockResults fetches the block results at a specific height,
// it then parses the tx results and block events into our generalized types.
// This allows us to maintain backwards compatability with older versions of CometBFT.
func (c *Client) GetBlockResults(ctx context.Context, height *int64) (*Results, error) {
	res, err := c.BlockResults(ctx, height)
	if err != nil {
		return nil, err
	}

	var txRes []*TxResult
	for _, tx := range res.TxsResults {
		txRes = append(txRes, &TxResult{
			Code:   tx.Code,
			Events: parseEvents(tx.Events),
		})
	}

	if (res.BeginBlockEvents != nil && len(res.BeginBlockEvents) > 0) &&
		(res.EndBlockEvents != nil && len(res.EndBlockEvents) > 0) {
		events := res.BeginBlockEvents
		events = append(events, res.EndBlockEvents...)

		return &Results{
			TxsResults: txRes,
			Events:     parseEvents(events),
		}, nil
	}

	return &Results{
		TxsResults: txRes,
		Events:     parseEvents(res.FinalizeBlockEvents),
	}, nil
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
	sdkEvents := make(sdk.StringEvents, 0, len(events))

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
