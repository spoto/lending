package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spoto/lending/x/lending/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_queryGetAllDebts(t *testing.T) {
	debtor, _ := sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
	creditor, _ := sdk.AccAddressFromBech32("cosmos1tupew4x3rhh0lpqha9wvzmzxjr4e37mfy3qefm")

	ID := "A1"
	amount := sdk.NewCoin("foo", sdk.NewInt(20000))

	tests := []struct {
		name  string
		debts []types.Debt
	}{
		{
			"no debts in store",
			nil,
		},
		{
			"one debt in store",
			[]types.Debt{
				{
					ID:       ID,
					Debtor:   debtor,
					Amount:   amount,
					Creditor: creditor,
				},
			},
		},
		{
			"debts in store",
			[]types.Debt{
				{
					ID:       ID,
					Debtor:   debtor,
					Amount:   amount,
					Creditor: creditor,
				},
				{
					ID:       ID + "A",
					Debtor:   creditor,
					Amount:   amount,
					Creditor: debtor,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cdc, ctx, _, _, keeper := SetupTestInput()

			for _, debt := range tt.debts {
				require.NoError(t, keeper.CreateDebt(ctx, debt))
			}

			result, err := queryGetAllDebts(ctx, nil, keeper)

			require.NoError(t, err)

			var d []types.Debt

			require.NotPanics(t, func() {
				cdc.MustUnmarshalJSON(result, &d)
			})

			require.Equal(t, tt.debts, d)
		})
	}
}