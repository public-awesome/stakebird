package keeper

import (
	"github.com/public-awesome/stargaze/x/alloc/types"
)

var _ types.QueryServer = Keeper{}
