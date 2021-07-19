package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/public-awesome/stargaze/simapp"
	"github.com/public-awesome/stargaze/x/alloc/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

type KeeperTestSuite struct {
	suite.Suite

	app *simapp.SimApp
	ctx sdk.Context
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.app = simapp.Setup(false)
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{Height: 1, Time: time.Now().UTC()})
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestDistributionToDAOAndDevs() {
	denom := suite.app.StakingKeeper.BondDenom(suite.ctx)
	allocKeeper := suite.app.AllocKeeper
	params := suite.app.AllocKeeper.GetParams(suite.ctx)
	devRewardsReceiver := sdk.AccAddress([]byte("addr1---------------"))
	params.DistributionProportions.DaoRewards = sdk.NewDecWithPrec(40, 2)
	params.DistributionProportions.DeveloperRewards = sdk.NewDecWithPrec(10, 2)
	params.WeightedDeveloperRewardsReceivers = []types.WeightedAddress{
		{
			Address: devRewardsReceiver.String(),
			Weight:  sdk.NewDec(1),
		},
	}
	suite.app.AllocKeeper.SetParams(suite.ctx, params)

	feePool := suite.app.DistrKeeper.GetFeePool(suite.ctx)
	feeCollector := suite.app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
	suite.Equal(
		"0",
		suite.app.BankKeeper.GetAllBalances(suite.ctx, feeCollector).AmountOf(denom).String())
	suite.Equal(
		sdk.NewDec(0),
		feePool.CommunityPool.AmountOf(denom))

	mintCoin := sdk.NewCoin(denom, sdk.NewInt(100000))
	mintCoins := sdk.Coins{mintCoin}
	err := suite.app.MintKeeper.MintCoins(suite.ctx, mintCoins)
	suite.NoError(err)

	err = suite.app.BankKeeper.SetBalances(suite.ctx, feeCollector, mintCoins)
	suite.NoError(err)

	feeCollector = suite.app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
	suite.Equal(
		mintCoin.Amount.String(),
		suite.app.BankKeeper.GetAllBalances(suite.ctx, feeCollector).AmountOf(denom).String())

	suite.Equal(
		sdk.NewDec(0),
		feePool.CommunityPool.AmountOf(denom))

	allocKeeper.DistributeInflation(suite.ctx)

	feeCollector = suite.app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
	feePool = suite.app.DistrKeeper.GetFeePool(suite.ctx)
	modulePortion := params.DistributionProportions.DaoRewards.
		Add(params.DistributionProportions.DeveloperRewards)
	suite.Equal(
		mintCoin.Amount.ToDec().Mul(modulePortion).RoundInt().String(),
		suite.app.BankKeeper.GetAllBalances(suite.ctx, feeCollector).AmountOf(denom).String())
	suite.Equal(
		mintCoin.Amount.ToDec().Mul(params.DistributionProportions.DeveloperRewards).TruncateInt(),
		suite.app.BankKeeper.GetBalance(suite.ctx, devRewardsReceiver, denom).Amount)
	// since the DAO is not setup yet, funds go into the communtiy pool
	suite.Equal(
		mintCoin.Amount.ToDec().Mul(params.DistributionProportions.DaoRewards),
		feePool.CommunityPool.AmountOf(denom))
}
