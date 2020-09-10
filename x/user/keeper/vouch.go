package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/public-awesome/stakebird/x/user/types"
)

// GetVouchesByVoucher returns the vouch if one exists
func (k Keeper) GetVouchesByVoucher(
	ctx sdk.Context, voucher sdk.AccAddress) (vouches []types.Vouch) {

	store := ctx.KVStore(k.storeKey)

	// iterator over vouches by a voucher
	it := sdk.KVStorePrefixIterator(store, types.VoucherPrefixKey(voucher))
	defer it.Close()

	for ; it.Valid(); it.Next() {
		var vouch types.Vouch
		k.cdc.MustUnmarshalBinaryBare(it.Value(), &vouch)
		vouches = append(vouches, vouch)
	}

	return vouches
}

// GetVouchByVouched returns the vouch if one exists
func (k Keeper) GetVouchByVouched(
	ctx sdk.Context, vouched sdk.AccAddress) (vouch types.Vouch, found bool, err error) {

	store := ctx.KVStore(k.storeKey)

	key := types.VouchedKey(vouched)
	value := store.Get(key)
	if value == nil {
		return vouch, false, nil
	}
	k.cdc.MustUnmarshalBinaryBare(value, &vouch)

	return vouch, true, nil
}

// IsVouched returns whether the given address has been previously vouched for
func (k Keeper) IsVouched(
	ctx sdk.Context, address sdk.AccAddress) (vouched bool) {

	store := ctx.KVStore(k.storeKey)

	key := types.VouchedKey(address)
	value := store.Get(key)

	return value != nil
}

// CanVouch returns whether the given address can vouch for someone
func (k Keeper) CanVouch(
	ctx sdk.Context, address sdk.AccAddress) (can bool) {

	vouches := k.GetVouchesByVoucher(ctx, address)

	// TODO: check the condition for the threshold amount

	// if already vouched enough time
	return uint32(len(vouches)) < k.GetParams(ctx).VouchCount
}

// CreateVouch registers a vouch on-chain.
func (k Keeper) CreateVouch(
	ctx sdk.Context, voucher, vouched sdk.AccAddress, comment string) error {

	if !k.CanVouch(ctx, voucher) {
		return sdkerrors.Wrap(
			sdkerrors.ErrInvalidRequest, fmt.Sprintf("given voucher cannot vouch %x", voucher))
	}

	if k.IsVouched(ctx, vouched) {
		return sdkerrors.Wrap(
			sdkerrors.ErrInvalidRequest, fmt.Sprintf("given account is already vouched for %x", vouched))
	}

	vouch := types.NewVouch(
		voucher, vouched, comment,
	)

	store := ctx.KVStore(k.storeKey)
	vouchedKey := types.VouchedKey(vouched)
	voucherKey := types.VoucherKey(voucher, vouched)
	value := k.cdc.MustMarshalBinaryBare(&vouch)
	store.Set(vouchedKey, value)
	store.Set(voucherKey, value)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypePost,
			// sdk.NewAttribute(types.AttributeKeyVendorID, fmt.Sprintf("%d", vendorID)),
			// sdk.NewAttribute(types.AttributeKeyPostID, postID),
			// sdk.NewAttribute(types.AttributeKeyCreator, creator.String()),
			// sdk.NewAttribute(types.AttributeKeyRewardAccount, rewardAccount.String()),
			// sdk.NewAttribute(types.AttributeKeyBody, body),
			// sdk.NewAttribute(types.AttributeKeyDeposit, d.String()),
			// sdk.NewAttribute(types.AttributeCurationEndTime, curationEndTime.Format(time.RFC3339)),
			// sdk.NewAttribute(types.AttributeKeyVoteDenom, types.DefaultVoteDenom),
		),
	})

	return nil
}
