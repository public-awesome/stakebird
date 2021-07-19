package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Implements the interface of `sdk.Msg`
// verify interface at compile time
var _ sdk.Msg = &MsgAfterMintSocialToken{}

// msg types
const (
	TypeMsgAfterMintSocialToken = "claim_after_mint_social_token"
)

// NewMsgAfterMintSocialToken creates a new MsgAfterMintSocialToken instance
func NewMsgAfterMintSocialToken(sender sdk.AccAddress,
) *MsgAfterMintSocialToken {
	return &MsgAfterMintSocialToken{
		Sender: sender.String(),
	}
}

// Route implements sdk.Msg
func (msg MsgAfterMintSocialToken) Route() string { return RouterKey }

// Type implements sdk.Msg
func (msg MsgAfterMintSocialToken) Type() string { return TypeMsgAfterMintSocialToken }

// GetSigners implements sdk.Msg
func (msg MsgAfterMintSocialToken) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgAfterMintSocialToken) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// ValidateBasic validity check for the AnteHandler
func (msg MsgAfterMintSocialToken) ValidateBasic() error {
	if strings.TrimSpace(msg.Sender) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "empty sender")
	}
	return nil
}
