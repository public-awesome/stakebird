package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// this line is used by starport scaffolding # genesis/types/import
// this line is used by starport scaffolding # ibc/genesistype/import

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: Params{
			DistributionProportions: DistributionProportions{
				DaoRewards:       sdk.NewDecWithPrec(35, 2), // 35%
				DeveloperRewards: sdk.NewDecWithPrec(10, 2), // 10%
			},
			WeightedDeveloperRewardsReceivers: []WeightedAddress{
				{"stars1xuuv5vucu9h74svhma4ykfvjzu4j0rxrsn0yfk", sdk.NewDecWithPrec(40, 2)},
				{"stars1s4ckh9405q0a3jhkwx9wkf9hsjh66nmuu53dwe", sdk.NewDecWithPrec(30, 2)},
				{"stars1kdfmfxg4tq68jxvl95h99wq9mvz9lxg6whrsjh", sdk.NewDecWithPrec(30, 2)},
			},
		},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// this line is used by starport scaffolding # ibc/genesistype/validate

	// this line is used by starport scaffolding # genesis/types/validate

	return nil
}
