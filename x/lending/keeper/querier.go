package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/spoto/lending/x/lending/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case types.QueryAllDebts:
			return queryGetAllDebts(ctx, path[1:], keeper)
		case types.QueryDebtorDebts:
			return queryGetDebtorDebts(ctx, path[1:], keeper)
		case types.QueryCreditorDebts:
			return queryGetCreditorDebts(ctx, path[1:], keeper)
		default:
			return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("Unknown %s query endpoint", types.ModuleName))
		}
	}
}

func queryGetAllDebts(ctx sdk.Context, _ []string, keeper Keeper) ([]byte, error) {
	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, keeper.GetAllDebts(ctx))
	if err2 != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Could not marshal result to JSON")
	}

	return bz, nil
}

func queryGetDebtorDebts(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	addr := path[0]
	address, _ := sdk.AccAddressFromBech32(addr)

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, keeper.GetDebtorDebts(ctx, address))
	if err2 != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Could not marshal result to JSON")
	}

	return bz, nil
}

func queryGetCreditorDebts(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	addr := path[0]
	address, _ := sdk.AccAddressFromBech32(addr)

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, keeper.GetCreditorDebts(ctx, address))
	if err2 != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Could not marshal result to JSON")
	}

	return bz, nil
}