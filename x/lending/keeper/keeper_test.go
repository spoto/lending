package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spoto/lending/x/lending/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestKeeper_CreateDebtWithAvailableID(t *testing.T) {
	debtor, _ := sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
	creditor, _ := sdk.AccAddressFromBech32("cosmos1tupew4x3rhh0lpqha9wvzmzxjr4e37mfy3qefm")

	debt1 := types.Debt{
		ID:       "A1",
		Debtor:   debtor,
		Amount:   sdk.NewCoin("foo", sdk.NewInt(20000)),
		Creditor: creditor,
	}

	_, ctx, _, _, keeper := SetupTestInput()
	err := keeper.CreateDebt(ctx, debt1)

	require.NoError(t, err)

	debt2 := types.Debt{
		ID:       "A2",
		Debtor:   debtor,
		Amount:   sdk.NewCoin("foo", sdk.NewInt(20000)),
		Creditor: creditor,
	}

	err = keeper.CreateDebt(ctx, debt2)

	require.NoError(t, err)
}

func TestKeeper_CreateDebtWithUnavailableID(t *testing.T) {
	debtor, _ := sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
	creditor, _ := sdk.AccAddressFromBech32("cosmos1tupew4x3rhh0lpqha9wvzmzxjr4e37mfy3qefm")

	debt1 := types.Debt{
		ID:       "A1",
		Debtor:   debtor,
		Amount:   sdk.NewCoin("foo", sdk.NewInt(20000)),
		Creditor: creditor,
	}

	_, ctx, _, _, keeper := SetupTestInput()
	err := keeper.CreateDebt(ctx, debt1)
	require.NoError(t, err)

	debt2 := types.Debt{
		ID:       "A1",  // same ID as before!
		Debtor:   creditor,
		Amount:   sdk.NewCoin("foo", sdk.NewInt(10000)),
		Creditor: debtor,
	}

	err = keeper.CreateDebt(ctx, debt2)
	require.Error(t, err)
}

func TestKeeper_PayDebt(t *testing.T) {
	debtor, _ := sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
	creditor, _ := sdk.AccAddressFromBech32("cosmos1tupew4x3rhh0lpqha9wvzmzxjr4e37mfy3qefm")

	ID := "A1"
	amount := sdk.NewCoin("foo", sdk.NewInt(20000))

	tests := []struct {
		name                  string
		preExistingDebt       *types.Debt
		msgPayDebt            types.MsgPayDebt
		startingDebtorBalance sdk.Coins
		wantErr               bool
	}{
		{
			"pay not existing debt",
			nil,
			types.MsgPayDebt{
				ID:     ID,
				Amount: amount,
				Debtor: debtor,
			},
			nil,
			true,
		},
		{
			"debtor has not enough funds to pay debt",
			&types.Debt{
				ID:       ID,
				Debtor:   debtor,
				Amount:   amount,
				Creditor: creditor,
			},
			types.MsgPayDebt{
				ID:     ID,
				Amount: amount,
				Debtor: debtor,
			},
			nil,
			true,
		},
		{
			"payer is not the debtor",
			&types.Debt{
				ID:       ID,
				Debtor:   debtor,
				Amount:   amount,
				Creditor: creditor,
			},
			types.MsgPayDebt{
				ID:     ID,
				Amount: amount,
				Debtor: creditor,
			},
			nil,
			true,
		},
		{
			"totally pay debt with exact amount",
			&types.Debt{
				ID:       ID,
				Debtor:   debtor,
				Amount:   amount,
				Creditor: creditor,
			},
			types.MsgPayDebt{
				ID:     ID,
				Amount: amount,
				Debtor: debtor,
			},
			sdk.NewCoins(amount),
			false,
		},

		{
			"totally pay debt with more amount than necessary",
			&types.Debt{
				ID:       ID,
				Debtor:   debtor,
				Amount:   amount,
				Creditor: creditor,
			},
			types.MsgPayDebt{
				ID:     ID,
				Amount: sdk.NewCoin("foo", sdk.NewInt(40000)),
				Debtor: debtor,
			},
			sdk.NewCoins(sdk.NewCoin("foo", sdk.NewInt(50000))),
			false,
		},
		{
			"partially pay debt",
			&types.Debt{
				ID:       ID,
				Debtor:   debtor,
				Amount:   amount,
				Creditor: creditor,
			},
			types.MsgPayDebt{
				ID:     ID,
				Amount: sdk.NewCoin("foo", sdk.NewInt(10000)),
				Debtor: debtor,
			},
			sdk.NewCoins(amount),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			require.NoError(t, tt.msgPayDebt.ValidateBasic())

			_, ctx, authKeeper, bankKeeper, keeper := SetupTestInput()

			if tt.preExistingDebt != nil {
				require.NoError(t, tt.preExistingDebt.Validate())
				require.NoError(t, keeper.CreateDebt(ctx, *tt.preExistingDebt))
				require.NoError(t, bankKeeper.SetCoins(ctx, tt.preExistingDebt.Creditor, sdk.NewCoins()))
			}

			if tt.startingDebtorBalance == nil {
				tt.startingDebtorBalance = sdk.NewCoins()
			}

			require.NoError(t, bankKeeper.SetCoins(ctx, tt.msgPayDebt.Debtor, tt.startingDebtorBalance))

			err := keeper.PayDebt(ctx, tt.msgPayDebt)

			if tt.wantErr {
				require.Error(t, err)
				if tt.preExistingDebt != nil {
					creditorAccount := authKeeper.GetAccount(ctx, creditor)
					require.True(t, creditorAccount.GetCoins().Empty())
				}

				debtorAccount := authKeeper.GetAccount(ctx, tt.msgPayDebt.Debtor)

				require.True(t, debtorAccount.GetCoins().IsEqual(tt.startingDebtorBalance))

				return
			}

			creditorAccount := authKeeper.GetAccount(ctx, tt.preExistingDebt.Creditor)
			require.True(t, creditorAccount.GetCoins().IsEqual(sdk.NewCoins(tt.msgPayDebt.Amount)))

			debtorAccount := authKeeper.GetAccount(ctx, tt.preExistingDebt.Debtor)
			require.True(t, debtorAccount.GetCoins().IsEqual(tt.startingDebtorBalance.Sub(sdk.NewCoins(tt.msgPayDebt.Amount))))

			newDebt, err := keeper.getDebtByID(ctx, tt.msgPayDebt.ID)
			require.NoError(t, err)

			expectedAmount := sdk.NewCoin(tt.preExistingDebt.Amount.Denom, sdk.NewInt(0))
			if tt.msgPayDebt.Amount.IsLT(tt.preExistingDebt.Amount) {
				expectedAmount = tt.preExistingDebt.Amount.Sub(tt.msgPayDebt.Amount)
			}

			require.True(t, newDebt.Amount.IsEqual(expectedAmount))
		})
	}
}

func TestKeeper_getDebtByID(t *testing.T) {
	debtor, _ := sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
	creditor, _ := sdk.AccAddressFromBech32("cosmos1tupew4x3rhh0lpqha9wvzmzxjr4e37mfy3qefm")

	debt := types.Debt{
		ID:       "A1",
		Debtor:   debtor,
		Amount:   sdk.NewCoin("foo", sdk.NewInt(20000)),
		Creditor: creditor,
	}

	tests := []struct {
		name            string
		ID              string
		preExistingDebt *types.Debt
		want            types.Debt
		wantErr         bool
	}{
		{
			"obtain existing debt",
			debt.ID,
			&debt,
			debt,
			false,
		},
		{
			"get debt from an empty store",
			debt.ID,
			nil,
			types.Debt{},
			true,
		},
		{
			"get not existing debt",
			debt.ID + "notExisting",
			&debt,
			types.Debt{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ctx, _, _, keeper := SetupTestInput()

			if tt.preExistingDebt != nil {
				require.NoError(t, keeper.CreateDebt(ctx, *tt.preExistingDebt))
			}

			result, err := keeper.getDebtByID(ctx, tt.ID)

			if tt.wantErr {
				require.Error(t, err)
			}

			require.Equal(t, tt.want, result)

		})
	}
}