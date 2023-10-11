package keeper

import (
	"baton/x/baton/types"
)

var _ types.QueryServer = Keeper{}
