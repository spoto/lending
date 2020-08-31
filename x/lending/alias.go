package lending

import (
	"github.com/spoto/lending/x/lending/keeper"
	"github.com/spoto/lending/x/lending/types"
)

const (
	ModuleName        = types.ModuleName
	RouterKey         = types.RouterKey
	StoreKey          = types.StoreKey
	DefaultParamspace = types.DefaultParamspace
	QuerierRoute      = types.QuerierRoute
)

var (
	// functions aliases
	NewKeeper           = keeper.NewKeeper
	NewQuerier          = keeper.NewQuerier
	RegisterCodec       = types.RegisterCodec
	NewGenesisState     = types.NewGenesisState
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis

	// variable aliases
	ModuleCdc = types.ModuleCdc

	NewMsgCreateDebt = types.NewMsgCreateDebt
	NewMsgPayDebt    = types.NewMsgPayDebt
	NewMsgChangeDebt = types.NewMsgChangeDebt
)

type (
	Keeper       = keeper.Keeper
	GenesisState = types.GenesisState
	Params       = types.Params

	MsgCreateDebt = types.MsgCreateDebt
	MsgPayDebt    = types.MsgPayDebt
	MsgChangeDebt = types.MsgChangeDebt
)