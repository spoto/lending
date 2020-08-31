package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/spoto/lending/x/lending/types"
)

// Keeper of the lending platform's store
type Keeper struct {

	// we include the keeper of Cosmos' Bank module
	bankKeeper bank.Keeper

	// we add an extra keeper to keep all debts created so far
	storeKey   sdk.StoreKey

	cdc        *codec.Codec
}

func NewKeeper(storeKey sdk.StoreKey, bankKeeper bank.Keeper, cdc *codec.Codec) Keeper {
	return Keeper{
		storeKey:   storeKey,
		bankKeeper: bankKeeper,
		cdc:        cdc,
	}
}

const debtStorePrefix = ":debt:"

func getDebtStoreKey(ID string) []byte {
	return []byte(debtStorePrefix + ID)
}

func (keeper Keeper) getDebtByID(ctx sdk.Context, id string) (types.Debt, error) {
	store := ctx.KVStore(keeper.storeKey)

	debtKey := getDebtStoreKey(id)
	if !store.Has(debtKey) {
		return types.Debt{}, fmt.Errorf("cannot find debt with ID %s", id)
	}

	var debt types.Debt
	keeper.cdc.MustUnmarshalBinaryBare(store.Get(debtKey), &debt)
	return debt, nil
}

func (keeper Keeper) CreateDebt(ctx sdk.Context, debt types.Debt) error {
	store := ctx.KVStore(keeper.storeKey)

	if !store.Has(getDebtStoreKey(debt.ID)) {
		store.Set(getDebtStoreKey(debt.ID), keeper.cdc.MustMarshalBinaryBare(&debt))
		return nil
	}

	return fmt.Errorf("cannot create a debt with an already used ID %s", debt.ID)
}

func (keeper Keeper) debtsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(keeper.storeKey)
	return sdk.KVStorePrefixIterator(store, []byte(debtStorePrefix))
}

// a functional type used to filter debts
type debtDiscrimination func(debt types.Debt) bool

func (keeper Keeper) getDebts(ctx sdk.Context, logic debtDiscrimination) []types.Debt {
	ri := keeper.debtsIterator(ctx)
	defer ri.Close()

	receivedResult := []types.Debt{}
	for ; ri.Valid(); ri.Next() {
		var debt types.Debt
		keeper.cdc.MustUnmarshalBinaryBare(ri.Value(), &debt)

		// we only accept debts that satisfy the filter
		if logic(debt) {
			receivedResult = append(receivedResult, debt)
		}
	}

	return receivedResult
}

func (keeper Keeper) GetAllDebts(ctx sdk.Context) []types.Debt {
	// we use a filter that accepts everything
	return keeper.getDebts(ctx, func(_ types.Debt) bool {
		return true
	})
}

func (keeper Keeper) GetDebtorDebts(ctx sdk.Context, address sdk.AccAddress) []types.Debt {
	// we use a filter that projects on the debtor
	return keeper.getDebts(ctx, func(debt types.Debt) bool {
		return debt.Debtor.Equals(address)
	})
}

func (keeper Keeper) GetCreditorDebts(ctx sdk.Context, address sdk.AccAddress) []types.Debt {
	// we use a filter that projects on the creditor
	return keeper.getDebts(ctx, func(debt types.Debt) bool {
		return debt.Creditor.Equals(address)
	})
}

func (keeper Keeper) PayDebt(ctx sdk.Context, msg types.MsgPayDebt) error {
	debt, err := keeper.getDebtByID(ctx, msg.ID)
	if err != nil {
		return err
	}

	if !msg.Debtor.Equals(debt.Debtor) {
		return fmt.Errorf("the debt with ID %s is not yours", msg.ID)
	}

	if err := keeper.bankKeeper.SendCoins(ctx, debt.Debtor, debt.Creditor, sdk.NewCoins(msg.Amount)); err != nil {
		return err
	}

	if msg.Amount.IsLT(debt.Amount) {
		debt.Amount = debt.Amount.Sub(msg.Amount)
	} else {
		debt.Amount = sdk.NewCoin(debt.Amount.Denom, sdk.NewInt(0))
	}

	return keeper.updateDebt(ctx, debt)
}

func (keeper Keeper) ChangeDebt(ctx sdk.Context, msg types.MsgChangeDebt) error {
	debt, err := keeper.getDebtByID(ctx, msg.ID)
	if err != nil {
		return err
	}

	if !msg.Creditor.Equals(debt.Creditor) {
		return fmt.Errorf("the debt with ID %s is not yours", msg.ID)
	}

	if msg.Amount.IsGTE(debt.Amount) {
		return fmt.Errorf("the new amount can only be smaller than the original %s", debt.Amount)
	}

	debt.Amount = debt.Amount.Sub(msg.Amount)
	if debt.Amount.IsNegative() {
		debt.Amount = sdk.NewCoin(debt.Amount.Denom, sdk.NewInt(0))
	}

	return keeper.updateDebt(ctx, debt)
}

func (keeper Keeper) updateDebt(ctx sdk.Context, debt types.Debt) error {
	store := ctx.KVStore(keeper.storeKey)

	if !store.Has(getDebtStoreKey(debt.ID)) {
		return fmt.Errorf("the debt with ID %s does not exist", debt.ID)
	}

	store.Set(getDebtStoreKey(debt.ID), keeper.cdc.MustMarshalBinaryBare(&debt))
	return nil
}
