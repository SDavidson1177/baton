package baton_test

import (
	"testing"

	keepertest "baton/testutil/keeper"
	"baton/testutil/nullify"
	"baton/x/baton"
	"baton/x/baton/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.BatonKeeper(t)
	baton.InitGenesis(ctx, *k, genesisState)
	got := baton.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
