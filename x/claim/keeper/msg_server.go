package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/public-awesome/stargaze/x/claim/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the curating MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (k msgServer) AfterMintSocialToken(
	goCtx context.Context, msg *types.MsgAfterMintSocialToken) (*types.MsgAfterMintSocialTokenResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	k.Keeper.AfterMintSocialToken(ctx, sender)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	return &types.MsgAfterMintSocialTokenResponse{}, nil
}
