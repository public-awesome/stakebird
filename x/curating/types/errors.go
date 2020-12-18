package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/curating module sentinel errors
var (
	ErrPostNotFound  = sdkerrors.Register(ModuleName, 1, "Post not found")
	ErrDuplicatePost = sdkerrors.Register(ModuleName, 2, "Post already exists")
	ErrPostExpired   = sdkerrors.Register(ModuleName, 3, "Post already expired")
)
