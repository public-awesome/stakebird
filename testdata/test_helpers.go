package testdata

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/public-awesome/stakebird/x/curating/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

var (
	PKs = simapp.CreateTestPubKeys(5)

	ValConsPk1 = PKs[0]
	ValConsPk2 = PKs[1]
	ValConsPk3 = PKs[2]

	ValConsAddr1 = sdk.ConsAddress(ValConsPk1.Address())
	ValConsAddr2 = sdk.ConsAddress(ValConsPk2.Address())

	distrAcc = auth.NewEmptyModuleAccount(types.ModuleName)
)

func CreateTestInput() (*codec.Codec, *SimApp, sdk.Context) {
	db := dbm.NewMemDB()
	logger := log.NewTMJSONLogger(log.NewSyncWriter(os.Stdout))

	opts := []func(*baseapp.BaseApp){baseapp.SetPruning(store.PruneNothing)}
	app := NewSimApp(logger, db, nil, true, 0, map[int64]bool{}, simapp.DefaultNodeHome, opts...)

	genesisState := ModuleBasics.DefaultGenesis(app.Codec())
	stateBytes, err := codec.MarshalJSONIndent(app.Codec(), genesisState)
	if err != nil {
		panic(err)
	}

	// Initialize the chain
	app.InitChain(
		abci.RequestInitChain{
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
		},
	)

	header := abci.Header{Height: app.LastBlockHeight() + 1, Time: time.Now()}
	ctx := app.NewContext(false, header)

	return codec.New(), app, ctx
}

type GenerateAccountStrategy func(int) []sdk.AccAddress

// AddTestAddrs constructs and returns accNum amount of accounts with an
// initial balance of accAmt in random order
func AddTestAddrsIncremental(app *SimApp, ctx sdk.Context, accNum int, accAmt sdk.Int) []sdk.AccAddress {
	return addTestAddrs(app, ctx, accNum, accAmt, createIncrementalAccounts)
}

// ConvertAddrsToValAddrs converts the provided addresses to ValAddress.
func ConvertAddrsToValAddrs(addrs []sdk.AccAddress) []sdk.ValAddress {
	valAddrs := make([]sdk.ValAddress, len(addrs))

	for i, addr := range addrs {
		valAddrs[i] = sdk.ValAddress(addr)
	}

	return valAddrs
}

func addTestAddrs(app *SimApp, ctx sdk.Context, accNum int, accAmt sdk.Int, strategy GenerateAccountStrategy) []sdk.AccAddress {
	testAddrs := strategy(accNum)

	stakeCoin := sdk.NewCoin("stake", accAmt)
	stbCoin := sdk.NewCoin("ustb", accAmt)
	atomCoin := sdk.NewCoin("uatom", accAmt)
	ibcCoin := sdk.NewCoin("transfer/ibczeroxfer/stake", accAmt)
	initCoins := sdk.NewCoins(stakeCoin, stbCoin, ibcCoin, atomCoin)
	setTotalSupply(app, ctx, accAmt, accNum)

	// fill all the addresses with some coins, set the loose pool tokens simultaneously
	for _, addr := range testAddrs {
		saveAccount(app, ctx, addr, initCoins)
	}

	return testAddrs
}

// setTotalSupply provides the total supply based on accAmt * totalAccounts.
func setTotalSupply(app *SimApp, ctx sdk.Context, accAmt sdk.Int, totalAccounts int) {
	totalSupply := sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), accAmt.MulRaw(int64(totalAccounts))))
	prevSupply := app.BankKeeper.GetSupply(ctx)
	app.BankKeeper.SetSupply(ctx, bank.NewSupply(prevSupply.GetTotal().Add(totalSupply...)))
}

// saveAccount saves the provided account into the simapp with balance based on initCoins.
func saveAccount(app *SimApp, ctx sdk.Context, addr sdk.AccAddress, initCoins sdk.Coins) {
	acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr)
	app.AccountKeeper.SetAccount(ctx, acc)
	_, err := app.BankKeeper.AddCoins(ctx, addr, initCoins)
	if err != nil {
		panic(err)
	}
}

// createIncrementalAccounts is a strategy used by addTestAddrs() in order to generated addresses in ascending order.
func createIncrementalAccounts(accNum int) []sdk.AccAddress {
	var addresses []sdk.AccAddress
	var buffer bytes.Buffer

	// start at 100 so we can make up to 999 test addresses with valid test addresses
	for i := 100; i < (accNum + 100); i++ {
		numString := strconv.Itoa(i)
		buffer.WriteString("A58856F0FD53BF058B4909A21AEC019107BA6") //base address string

		buffer.WriteString(numString) //adding on final two digits to make addresses unique
		res, _ := sdk.AccAddressFromHex(buffer.String())
		bech := res.String()
		addr, _ := FakeAddr(buffer.String(), bech)

		addresses = append(addresses, addr)
		buffer.Reset()
	}

	return addresses
}

func FakeAddr(addr string, bech string) (sdk.AccAddress, error) {
	res, err := sdk.AccAddressFromHex(addr)
	if err != nil {
		return nil, err
	}
	bechexpected := res.String()
	if bech != bechexpected {
		return nil, fmt.Errorf("bech encoding doesn't match reference")
	}

	bechres, err := sdk.AccAddressFromBech32(bech)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(bechres, res) {
		return nil, err
	}

	return res, nil
}
