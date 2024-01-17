package cometbft_client

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/strangelove-ventures/cometbft-client/libs/bytes"
	"github.com/stretchr/testify/require"
)

const url = "https://rpc.osmosis.strange.love:80"

// TODO: this hardcoded value makes the test brittle since the underlying node may not have this state persisted
var blockHeight = int64(13311684)

func TestClientStatus(t *testing.T) {
	client := NewClient()
	err := client.Init(url, "0.37.2", 5*time.Second)
	if err != nil {
		t.Fatalf("Failed to initialize client: %v", err)
	}

	res, err := client.RPCClient.Status(context.Background())
	if err != nil {
		t.Fatalf("Failed to get client status: %v", err)
	}

	t.Logf("Status Resp: %v \n", res)
}

func TestBlockResults(t *testing.T) {
	client := NewClient()
	err := client.Init(url, "0.37.2", 5*time.Second)
	if err != nil {
		t.Fatalf("Failed to initialize client: %v", err)
	}

	res, err := client.RPCClient.BlockResults(context.Background(), nil)
	if err != nil {
		t.Fatalf("Failed to get block results: %v", err)
	}

	t.Logf("Block Results: %v \n", res)
}

func TestABCIInfo(t *testing.T) {
	client := NewClient()
	err := client.Init(url, "0.37.2", 5*time.Second)
	if err != nil {
		t.Fatalf("Failed to initialize client: %v", err)
	}

	res, err := client.RPCClient.ABCIInfo(context.Background())
	if err != nil {
		t.Fatalf("Failed to get ABCI info: %v", err)
	}

	t.Logf("ABCI Info: %v \n", res)
}

func TestABCIQuery(t *testing.T) {
	client := NewClient()
	err := client.Init(url, "0.37.2", 5*time.Second)
	if err != nil {
		t.Fatalf("Failed to initialize client: %v", err)
	}

	// TODO: pass in valid values for path and data
	path := ""
	data := bytes.HexBytes{}

	res, err := client.RPCClient.ABCIQuery(context.Background(), path, data)
	if err != nil {
		t.Fatalf("Failed to query ABCI: %v", err)
	}

	t.Logf("ABCI Query: %v \n", res)
}

func TestBlockByHeight(t *testing.T) {
	client := NewClient()
	err := client.Init(url, "0.37.2", 5*time.Second)
	if err != nil {
		t.Fatalf("Failed to initialize client: %v", err)
	}

	res, err := client.BlockResults(context.Background(), &blockHeight)
	if err != nil {
		t.Fatalf("Failed to get block results: %v", err)
	}

	t.Logf("Block Results: %v \n", res)
}

func TestConsensusParams(t *testing.T) {
	client := NewClient()
	err := client.Init(url, "0.37.2", 5*time.Second)
	if err != nil {
		t.Fatalf("Failed to initialize client: %v", err)
	}

	res, err := client.RPCClient.ConsensusParams(context.Background(), &blockHeight)
	if err != nil {
		t.Fatalf("Failed to get consensus params: %v", err)
	}

	t.Logf("Consensus Params: %v \n", res)
}

func TestConsensusState(t *testing.T) {
	client := NewClient()
	err := client.Init(url, "0.37.2", 5*time.Second)
	if err != nil {
		t.Fatalf("Failed to initialize client: %v", err)
	}

	res, err := client.RPCClient.ConsensusState(context.Background())
	if err != nil {
		t.Fatalf("Failed to get consensus state: %v", err)
	}

	t.Logf("Consensus State: %v \n", res)
}

func TestDumpConsensusState(t *testing.T) {
	client := NewClient()
	err := client.Init(url, "0.37.2", 5*time.Second)
	if err != nil {
		t.Fatalf("Failed to initialize client: %v", err)
	}

	res, err := client.RPCClient.DumpConsensusState(context.Background())
	if err != nil {
		t.Fatalf("Failed to dump consensus state: %v", err)
	}

	t.Logf("Dump Consensus State: %v \n", res)
}

func TestGenesis(t *testing.T) {
	client := NewClient()
	err := client.Init(url, "0.37.2", 5*time.Second)
	if err != nil {
		t.Fatalf("Failed to initialize client: %v", err)
	}

	res, err := client.RPCClient.Genesis(context.Background())
	if err != nil && !strings.Contains(err.Error(), "genesis response is large, please use the genesis_chunked API instead") {
		t.Fatalf("Failed to get genesis: %v", err)
	}

	t.Logf("Genesis: %v \n", res)
}

func TestGenesisChunked(t *testing.T) {
	client := NewClient()
	err := client.Init(url, "0.37.2", 30*time.Second)
	if err != nil {
		t.Fatalf("Failed to initialize client: %v", err)
	}

	chunk := uint(1)
	res, err := client.RPCClient.GenesisChunked(context.Background(), chunk)
	if err != nil {
		t.Fatalf("Failed to get genesis chunk: %v", err)
	}

	t.Logf("Genesis Chunk: %v \n", res)
}

func TestHealth(t *testing.T) {
	client := NewClient()
	err := client.Init(url, "0.37.2", 5*time.Second)
	if err != nil {
		t.Fatalf("Failed to initialize client: %v", err)
	}

	res, err := client.RPCClient.Health(context.Background())
	if err != nil {
		t.Fatalf("Failed to get health status: %v", err)
	}

	t.Logf("Health Status: %v \n", res)
}

func TestNetInfo(t *testing.T) {
	client := NewClient()
	err := client.Init(url, "0.37.2", 5*time.Second)
	if err != nil {
		t.Fatalf("Failed to initialize client: %v", err)
	}

	res, err := client.RPCClient.NetInfo(context.Background())
	if err != nil {
		t.Fatalf("Failed to get network info: %v", err)
	}

	t.Logf("Network Info: %v \n", res)
}

func TestNumUnconfirmedTxs(t *testing.T) {
	client := NewClient()
	err := client.Init(url, "0.37.2", 5*time.Second)
	if err != nil {
		t.Fatalf("Failed to initialize client: %v", err)
	}

	res, err := client.RPCClient.NumUnconfirmedTxs(context.Background())
	if err != nil {
		t.Fatalf("Failed to get number of unconfirmed txs: %v \n", err)
	}

	t.Logf("Num Of Unconfirmed Txs: %v \n", res)
}

func TestUnconfirmedTxs(t *testing.T) {
	client := NewClient()
	err := client.Init(url, "0.37.2", 5*time.Second)
	if err != nil {
		t.Fatalf("Failed to initialize client: %v", err)
	}

	limit := 5
	res, err := client.RPCClient.UnconfirmedTxs(context.Background(), &limit)
	if err != nil {
		t.Fatalf("Failed to get unconfirmed txs with limit %d: %v \n", limit, err)
	}

	t.Logf("Unconfirmed Txs: %v \n", res)
	require.Equal(t, limit+1, res.Count) // TODO: upstream off by one error?
}

func TestValidators(t *testing.T) {
	client := NewClient()
	err := client.Init(url, "0.37.2", 5*time.Second)
	if err != nil {
		t.Fatalf("Failed to initialize client: %v", err)
	}

	page := 1
	perPage := 5

	res, err := client.RPCClient.Validators(context.Background(), &blockHeight, &page, &perPage)
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
