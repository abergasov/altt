package routes_test

import (
	"altt/internal/entities"
	testhelpers "altt/internal/test_helpers"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	chain   = entities.ChainEthereum
	address = "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"
	token   = entities.USDC
)

func TestServer_GetNativeBalance(t *testing.T) {
	// given
	tCtx := testhelpers.GetClean(t)
	srv := testhelpers.NewTestServer(t, tCtx)

	t.Run("known chain", func(t *testing.T) {
		// when
		resp := srv.Get(t, fmt.Sprintf("/%s/balance/%s", chain.String(), address))
		resp.RequireOk(t)

		// then
		var response entities.Balance
		resp.RequireUnmarshal(t, &response)
		require.Equal(t, entities.ChainEthereum, response.Chain)
		t.Log(fmt.Sprintf("balance is %s %s", response.TokenBalance, entities.MapChainToFuel(chain)))
	})

	t.Run("unknown chain", func(t *testing.T) {
		// when
		resp := srv.Get(t, "/abc/balance/"+address)
		resp.RequireNotFound(t)
	})
}

func TestServer_GetKnownTokenBalanceByAddress(t *testing.T) {
	// given
	tCtx := testhelpers.GetClean(t)
	srv := testhelpers.NewTestServer(t, tCtx)

	t.Run("known chain and token", func(t *testing.T) {
		// when
		resp := srv.Get(t, fmt.Sprintf("/%s/%s/balance/%s", chain.String(), token, address))
		resp.RequireOk(t)

		// then
		var response entities.Balance
		resp.RequireUnmarshal(t, &response)
		require.Equal(t, entities.ChainEthereum, response.Chain)
		t.Log(fmt.Sprintf("balance is %s %s", response.TokenBalance, token))
		t.Log("https://etherscan.io/tokenholdings?a=" + address)
	})
}
