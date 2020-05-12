package keeper

import (
	"errors"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/davecgh/go-spew/spew"
	"github.com/rocket-protocol/stakebird/x/stake/types"
)

// Perform a delegation
func (k Keeper) Delegate(ctx sdk.Context, vendorID, postID uint64, downvote bool,
	delAddr sdk.AccAddress, valAddress sdk.ValAddress, amount sdk.Coin) (err error) {

	// check if post exist, if not, create it and begin the voting period
	post, found := k.GetPost(ctx, vendorID, postID)
	if !found {
		post = k.CreatePost(ctx, postID, vendorID, "", k.VotingPeriod(ctx))
	}

	// find validator
	validator, found := k.sk.GetValidator(ctx, valAddress)
	if !found {
		return errors.New("validator not found")
	}

	// add delegation to voting delegation queue
	shares := amount.Amount.ToDec()
	delegation := stakingtypes.NewDelegation(delAddr, valAddress, shares)
	k.InsertVotingDelegationQueue(ctx, vendorID, postID, delegation, post.VoteEnd)

	// perform delegation on chain
	_, err = k.sk.Delegate(ctx, delAddr, amount.Amount, sdk.Unbonded, validator, false)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) Undelegate(ctx sdk.Context, endTime time.Time, vendorID, postID, stakeID uint64) error {
	d := k.getDelegation(ctx, endTime, vendorID, postID, stakeID)
	_, err := k.sk.Unbond(ctx, d.DelegatorAddress, d.ValidatorAddress, d.Shares)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) getDelegation(ctx sdk.Context, endTime time.Time, vendorID, postID, stakeID uint64) stakingtypes.Delegation {
	store := ctx.KVStore(k.storeKey)
	key := types.VotingDelegationQueueKey(endTime, vendorID, postID, stakeID)
	value := store.Get(key)
	var delegation stakingtypes.Delegation
	k.cdc.UnmarshalBinaryBare(value, &delegation)

	return delegation
}

func (k Keeper) InsertVotingDelegationQueue(ctx sdk.Context, vendorID, postID uint64,
	delegation stakingtypes.Delegation, completionTime time.Time) {
	// get current stake index
	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.KeyIndexStakeID)
	stakeID := types.StakeIndexFromKey(value)

	// create voting delegation queue key
	queueKey := types.VotingDelegationQueueKey(completionTime, vendorID, postID, stakeID)

	value = store.Get(queueKey)
	if len(value) == 0 {
		// add to queue at the right time
		k.setVotingDelegationQueue(ctx, queueKey, delegation)
	}

	// store incremented index
	store.Set(types.KeyIndexStakeID, types.StakeIndexToKey(stakeID+1))
}

func (k Keeper) RemoveFromVotingDelegationQueue(ctx sdk.Context, endTime time.Time, vendorID, postID, stakeID uint64) {
	store := ctx.KVStore(k.storeKey)
	key := types.VotingDelegationQueueKey(endTime, vendorID, postID, stakeID)
	store.Delete(key)
}

func (k Keeper) setVotingDelegationQueue(ctx sdk.Context, key []byte, delegation stakingtypes.Delegation) {
	store := ctx.KVStore(k.storeKey)
	value := k.cdc.MustMarshalBinaryBare(&delegation)
	store.Set(key, value)
}

func (k Keeper) IterateVotingDelegationQueue(ctx sdk.Context, endTime time.Time,
	cb func(endTime time.Time, vendorID, postID, stakeID uint64, delegation stakingtypes.Delegation) (stop bool)) {

	iterator := k.VotingDelegationQueueIterator(ctx, endTime)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		spew.Dump(iterator.Key())
		endTime, vendorID, postID, stakeID := types.SplitVotingDelegationQueueKey(iterator.Key())
		// spew.Dump("vendorID, postID, stakeID", vendorID, postID, stakeID)
		var delegation stakingtypes.Delegation
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &delegation)

		if cb(endTime, vendorID, postID, stakeID, delegation) {
			break
		}
	}
}

// Returns all delegation timeslices from time 0 until completion time
func (k Keeper) VotingDelegationQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(types.KeyPrefixVotingDelegationQueue,
		sdk.InclusiveEndBytes(types.VotingDelegationQueueTimeKeyPrefix(endTime)))
}
