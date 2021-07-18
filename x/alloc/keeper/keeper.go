package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/public-awesome/stargaze/x/alloc/types"
	// this line is used by starport scaffolding # ibc/keeper/import
)

type (
	Keeper struct {
		cdc      codec.Marshaler
		storeKey sdk.StoreKey
		memKey   sdk.StoreKey
		// this line is used by starport scaffolding # ibc/keeper/attribute

		accountKeeper types.AccountKeeper
		bankKeeper    types.BankKeeper
		stakingKeeper types.StakingKeeper
		distrKeeper   types.DistrKeeper
	}
)

func NewKeeper(
	cdc codec.Marshaler,
	storeKey,
	memKey sdk.StoreKey,
	// this line is used by starport scaffolding # ibc/keeper/parameter
	accountKeeper types.AccountKeeper, bankKeeper types.BankKeeper,
	stakingKeeper types.StakingKeeper,
	distrKeeper types.DistrKeeper,
) *Keeper {
	return &Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		memKey:   memKey,
		// this line is used by starport scaffolding # ibc/keeper/return
		accountKeeper: accountKeeper, bankKeeper: bankKeeper,
		stakingKeeper: stakingKeeper,
		distrKeeper:   distrKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetModuleAccountBalance gets the airdrop coin balance of module account
func (k Keeper) GetModuleAccountAddress(ctx sdk.Context) sdk.AccAddress {
	return k.accountKeeper.GetModuleAddress(types.ModuleName)
}

// DistributeInflation distributes module-specific inflation
func (k Keeper) DistributeInflation(ctx sdk.Context) error {
	blockInflationAddr := k.accountKeeper.GetModuleAccount(ctx, authtypes.FeeCollectorName).GetAddress()
	blockInflation := k.bankKeeper.GetBalance(ctx, blockInflationAddr, k.stakingKeeper.BondDenom(ctx))
	blockInflationDec := sdk.NewDecFromInt(blockInflation.Amount)

	params, err := k.GetParams(ctx)
	if err != nil {
		return err
	}
	proportions := params.DistributionProportions

	daoRewardAmount := blockInflationDec.Mul(proportions.DaoRewards)
	daoRewardCoin := sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), daoRewardAmount.TruncateInt())
	// Distribute DAO incentives to the community pool until StargazeDAO is implemented
	err = k.distrKeeper.FundCommunityPool(ctx, sdk.NewCoins(daoRewardCoin), blockInflationAddr)
	if err != nil {
		return err
	}
	k.Logger(ctx).Info("funded community pool")

	devRewardAmount := blockInflationDec.Mul(proportions.DeveloperRewards)
	devRewardCoin := sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), devRewardAmount.TruncateInt())

	for _, w := range params.WeightedDeveloperRewardsReceivers {
		devRewardPortionCoins := sdk.NewCoins(k.GetProportions(ctx, devRewardCoin, w.Weight))
		if w.Address == "" {
			err := k.distrKeeper.FundCommunityPool(ctx, devRewardPortionCoins, blockInflationAddr)
			if err != nil {
				return err
			}
		} else {
			devRewardsAddr, err := sdk.AccAddressFromBech32(w.Address)
			if err != nil {
				return err
			}
			// // If recipient is vesting account, pay to account according to its vesting condition
			// err = k.bankKeeper.SendCoinsFromModuleToAccountOriginalVesting(
			// 	ctx, types.DeveloperVestingModuleAcctName, devRewardsAddr, devRewardPortionCoins)
			err = k.bankKeeper.SendCoins(ctx, blockInflationAddr, devRewardsAddr, devRewardPortionCoins)
			if err != nil {
				return err
			}
			k.Logger(ctx).Info("sent coins to developer")
		}
	}

	return nil
}

// GetProportions gets the balance of the `MintedDenom` from minted coins and returns coins according to the `AllocationRatio`
func (k Keeper) GetProportions(ctx sdk.Context, mintedCoin sdk.Coin, ratio sdk.Dec) sdk.Coin {
	return sdk.NewCoin(mintedCoin.Denom, mintedCoin.Amount.ToDec().Mul(ratio).TruncateInt())
}
