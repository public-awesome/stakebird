package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/public-awesome/stargaze/x/claim/types"
)

func (k Keeper) AfterMintSocialToken(ctx sdk.Context, sender sdk.AccAddress) {
	_, err := k.ClaimCoinsForAction(ctx, sender, types.ActionMintSocialToken)
	if err != nil {
		panic(err.Error())
	}
}

func (k Keeper) AfterBuySocialToken(ctx sdk.Context, sender sdk.AccAddress) {
	_, err := k.ClaimCoinsForAction(ctx, sender, types.ActionBuySocialToken)
	if err != nil {
		panic(err.Error())
	}
}

func (k Keeper) AfterProposalVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress) {
	_, err := k.ClaimCoinsForAction(ctx, voterAddr, types.ActionVote)
	if err != nil {
		panic(err.Error())
	}
}

func (k Keeper) AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	_, err := k.ClaimCoinsForAction(ctx, delAddr, types.ActionDelegateStake)
	if err != nil {
		panic(err.Error())
	}
}

// ________________________________________________________________________________________

// Hooks wrapper struct for slashing keeper
type Hooks struct {
	k Keeper
}

var _ stakingtypes.StakingHooks = Hooks{}

// Return the wrapper struct
func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

// gamm hooks
func (h Hooks) AfterPoolCreated(
	ctx sdk.Context, sender sdk.AccAddress, poolID uint64) {
	h.k.AfterMintSocialToken(ctx, sender)
}
func (h Hooks) AfterJoinPool(
	ctx sdk.Context, sender sdk.AccAddress, poolID uint64, enterCoins sdk.Coins, shareOutAmount sdk.Int) {
	h.k.AfterMintSocialToken(ctx, sender)
}
func (h Hooks) AfterExitPool(
	ctx sdk.Context, sender sdk.AccAddress, poolID uint64, shareInAmount sdk.Int, exitCoins sdk.Coins) {
}
func (h Hooks) AfterBuySocialToken(
	ctx sdk.Context, sender sdk.AccAddress, poolID uint64, input sdk.Coins, output sdk.Coins) {
	h.k.AfterBuySocialToken(ctx, sender)
}

// governance hooks
func (h Hooks) AfterProposalSubmission(ctx sdk.Context, proposalID uint64) {}
func (h Hooks) AfterProposalDeposit(ctx sdk.Context, proposalID uint64, depositorAddr sdk.AccAddress) {
}

func (h Hooks) AfterProposalVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress) {
	h.k.AfterProposalVote(ctx, proposalID, voterAddr)
}

func (h Hooks) AfterProposalInactive(ctx sdk.Context, proposalID uint64) {}
func (h Hooks) AfterProposalActive(ctx sdk.Context, proposalID uint64)   {}

// staking hooks
func (h Hooks) AfterValidatorCreated(ctx sdk.Context, valAddr sdk.ValAddress)   {}
func (h Hooks) BeforeValidatorModified(ctx sdk.Context, valAddr sdk.ValAddress) {}
func (h Hooks) AfterValidatorRemoved(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
}
func (h Hooks) AfterValidatorBonded(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
}
func (h Hooks) AfterValidatorBeginUnbonding(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
}
func (h Hooks) BeforeDelegationCreated(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
}
func (h Hooks) BeforeDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
}
func (h Hooks) BeforeDelegationRemoved(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
}
func (h Hooks) AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	h.k.AfterDelegationModified(ctx, delAddr, valAddr)
}
func (h Hooks) BeforeValidatorSlashed(ctx sdk.Context, valAddr sdk.ValAddress, fraction sdk.Dec) {}
