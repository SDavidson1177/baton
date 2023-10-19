package keeper

import (
	"baton/x/ledger/types"
)

var _ types.QueryServer = Keeper{}
