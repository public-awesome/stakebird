package types

import (
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Default parameter namespace
const (
	DefaultParamspace = ModuleName
	// TODO: Define your default parameters
)

// Parameter store keys
var (
// TODO: Define your keys for the parameter store
// KeyParamName          = []byte("ParamName")
)

// ParamKeyTable for funding module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// Params - used for initializing default parameter for funding at genesis
type Params struct {
	// TODO: Add your parameters to the parameter struct
	// KeyParamName string `json:"key_param_name"`
}

// NewParams creates a new Params object
func NewParams( /* TODO: Pass in the parameters*/ ) Params {
	return Params{
		// TODO: Create your Params Type
	}
}

// String implements the stringer interface for Params
func (p Params) String() string {
	// TODO: Return all the params as a string
	return ""
}

// ParamSetPairs - Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		// TODO: Pair your key with the param
		// params.NewParamSetPair(KeyParamName, &p.ParamName),
	}
}

// DefaultParams defines the parameters for this module
func DefaultParams() Params {
	return NewParams( /* TODO: Pass in your default Params */ )
}
