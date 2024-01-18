package client

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/strangelove-ventures/cometbft-client/libs/bytes"
	"github.com/stretchr/testify/require"
)

const url = "https://rpc.osmosis.strange.love:443"

// TODO: this hardcoded value makes the test brittle since the underlying node may not have this state persisted
var blockHeight = int64(13311684)

func testClient(t *testing.T) *Client {
	client, err := NewClient(url, 5*time.Second)
	require.NoError(t, err, "failed to initialize client")

	return client
}

func TestClientStatus(t *testing.T) {
	client := testClient(t)

	res, err := client.Status(context.Background())
	require.NoError(t, err, "failed to get client status")

	resJson, err := json.Marshal(res)
	require.NoError(t, err)

	t.Logf("Status Resp: %s \n", resJson)
}

func TestBlockResults(t *testing.T) {
	client := testClient(t)

	ctx := context.Background()
	res, err := client.BlockResults(ctx, nil)
	require.NoError(t, err, "failed to get block results")

	resJson, err := json.Marshal(res)
	require.NoError(t, err)

	t.Logf("Block Results: %s \n", resJson)

	res2, err := client.GetBlockResults(ctx, nil)
	require.NoError(t, err)

	res2Json, err := json.Marshal(res2)
	require.NoError(t, err)

	t.Logf("Block Results: %s \n", res2Json)
}

func TestABCIInfo(t *testing.T) {
	client := testClient(t)

	res, err := client.ABCIInfo(context.Background())
	require.NoError(t, err, "failed to get ABCI info")

	resJson, err := json.Marshal(res)
	require.NoError(t, err)

	t.Logf("ABCI Info: %s \n", resJson)
}

func TestABCIQuery(t *testing.T) {
	client := testClient(t)

	// TODO: pass in valid values for path and data
	path := ""
	data := bytes.HexBytes{}

	res, err := client.ABCIQuery(context.Background(), path, data)
	require.NoError(t, err, "failed to query ABCI")

	require.Equal(t, uint32(6), res.Response.Code)
	require.Equal(t, "no query path provided: unknown request", res.Response.Log)
	require.Equal(t, "sdk", res.Response.Codespace)

	resJson, err := json.Marshal(res)
	require.NoError(t, err)

	t.Logf("ABCI Query: %s \n", resJson)
}

func TestBlockByHeight(t *testing.T) {
	client := testClient(t)

	res, err := client.BlockResults(context.Background(), &blockHeight)
	require.NoError(t, err, "failed to get block results")

	resJson, err := json.Marshal(res)
	require.NoError(t, err)

	t.Logf("Block Results: %s \n", resJson)
}

func TestConsensusParams(t *testing.T) {
	client := testClient(t)

	res, err := client.ConsensusParams(context.Background(), &blockHeight)
	if err != nil {
		t.Fatalf("Failed to get consensus params: %v", err)
	}

	t.Logf("Consensus Params: %v \n", res)
}

func TestConsensusState(t *testing.T) {
	client := testClient(t)

	res, err := client.ConsensusState(context.Background())
	if err != nil {
		t.Fatalf("Failed to get consensus state: %v", err)
	}

	t.Logf("Consensus State: %v \n", res)
}

func TestDumpConsensusState(t *testing.T) {
	client := testClient(t)

	res, err := client.DumpConsensusState(context.Background())
	if err != nil {
		t.Fatalf("Failed to dump consensus state: %v", err)
	}

	t.Logf("Dump Consensus State: %v \n", res)
}

func TestGenesis(t *testing.T) {
	client := testClient(t)

	res, err := client.Genesis(context.Background())
	if err != nil && !strings.Contains(err.Error(), "genesis response is large, please use the genesis_chunked API instead") {
		t.Fatalf("Failed to get genesis: %v", err)
	}

	t.Logf("Genesis: %v \n", res)
}

func TestGenesisChunked(t *testing.T) {
	client := testClient(t)

	chunk := uint(1)
	res, err := client.GenesisChunked(context.Background(), chunk)
	if err != nil {
		t.Fatalf("Failed to get genesis chunk: %v", err)
	}

	t.Logf("Genesis Chunk: %v \n", res)
}

func TestHealth(t *testing.T) {
	client := testClient(t)

	res, err := client.Health(context.Background())
	if err != nil {
		t.Fatalf("Failed to get health status: %v", err)
	}

	t.Logf("Health Status: %v \n", res)
}

func TestNetInfo(t *testing.T) {
	client := testClient(t)

	res, err := client.NetInfo(context.Background())
	if err != nil {
		t.Fatalf("Failed to get network info: %v", err)
	}

	t.Logf("Network Info: %v \n", res)
}

func TestNumUnconfirmedTxs(t *testing.T) {
	client := testClient(t)

	res, err := client.NumUnconfirmedTxs(context.Background())
	if err != nil {
		t.Fatalf("Failed to get number of unconfirmed txs: %v \n", err)
	}

	t.Logf("Num Of Unconfirmed Txs: %v \n", res)
}

func TestUnconfirmedTxs(t *testing.T) {
	client := testClient(t)

	limit := 5
	res, err := client.UnconfirmedTxs(context.Background(), &limit)
	if err != nil {
		t.Fatalf("Failed to get unconfirmed txs with limit %d: %v \n", limit, err)
	}

	t.Logf("Unconfirmed Txs: %v \n", res)
	require.Equal(t, limit+1, res.Count) // TODO: upstream off by one error?
}

func TestValidators(t *testing.T) {
	client := testClient(t)

	page := 1
	perPage := 5

	res, err := client.Validators(context.Background(), &blockHeight, &page, &perPage)
	if err != nil {
		t.Fatalf("Failed to get validators: %v", err)
	}

	t.Logf("Validators: %v \n", res)
	require.Equal(t, perPage, res.Count)
}

func TestBlockByHash(t *testing.T) {

}

func TestBlockSearch(t *testing.T) {

}

func TestBlockchainMinMaxHeight(t *testing.T) {

}

// TODO: is this necessary?
func TestBroadcastEvidence(t *testing.T) {

}

func TestBroadcastTxAsync(t *testing.T) {

}

func TestBroadcastTxCommit(t *testing.T) {

}

func TestBroadcastTxSync(t *testing.T) {

}

func TestCheckTx(t *testing.T) {

}

func TestCommit(t *testing.T) {

}

func TestUnsubscribeByQuery(t *testing.T) {

}

func TestUnsubscribeAll(t *testing.T) {

}

func TestSubscribe(t *testing.T) {

}

func TestTxByHash(t *testing.T) {

}

func TestTxSearch(t *testing.T) {

}
