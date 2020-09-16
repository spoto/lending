package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spoto/lending/x/lending/types"
)

func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "debts_are_positive", DebtsArePositive(k))
	ir.RegisterRoute(types.ModuleName, "debts_are_not_reflexive", DebtsAreNotReflexive(k))
}

func DebtsArePositive(keeper Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		debts := keeper.GetAllDebts(ctx)
		negative := false
		for _, debt := range debts {
			negative = negative || debt.Amount.IsNegative()
		}

		return sdk.FormatInvariant(types.ModuleName,
			"negative debt",
			"A debt has a negative amount"),
			negative
	}
}

func DebtsAreNotReflexive(keeper Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		debts := keeper.GetAllDebts(ctx)
		reflexive := false
		for _, debt := range debts {
			reflexive = reflexive || debt.Creditor.Equals(debt.Debtor)
		}

		return sdk.FormatInvariant(types.ModuleName,
			"reflexive debt",
			"The creditor of a debt coincides with its debtor"),
			reflexive
	}
}