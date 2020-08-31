package lending

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/spoto/lending/x/lending/types"
)

func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case types.MsgCreateDebt:
			return handleMsgCreateDebt(ctx, keeper, msg)
		case types.MsgPayDebt:
			return handleMsgPayDebt(ctx, keeper, msg)
		case types.MsgChangeDebt:
			return handleMsgChangeDebt(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized %s message type: %v", types.ModuleName, msg.Type())
			return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgChangeDebt(ctx sdk.Context, keeper Keeper, msg types.MsgChangeDebt) (*sdk.Result, error) {
	err := keeper.ChangeDebt(ctx, msg)
	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrInvalidRequest, err.Error())
	}

	return &sdk.Result{Log: "Debt changed successfully"}, nil
}

func handleMsgCreateDebt(ctx sdk.Context, keeper Keeper, msg types.MsgCreateDebt) (*sdk.Result, error) {
	err := keeper.CreateDebt(ctx, types.Debt(msg))
	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrInvalidRequest, err.Error())
	}

	return &sdk.Result{Log: "Debt created successfully"}, nil
}

func handleMsgPayDebt(ctx sdk.Context, keeper Keeper, msg types.MsgPayDebt) (*sdk.Result, error) {
	err := keeper.PayDebt(ctx, msg)
	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrInvalidRequest, err.Error())
	}

	return &sdk.Result{Log: "Debt financed successfully"}, nil
}