package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
)

type Debt struct {
	ID     string         `json:"ID"`
	Debtor sdk.AccAddress `json:"debtor"`
	Amount   sdk.Coin       `json:"amount"`
	Creditor sdk.AccAddress `json:"creditor"`
}

func (d Debt) Validate() error {
	if d.ID == "" {
		return sdkErr.Wrap(sdkErr.ErrInvalidRequest, "ID can't be empty")
	}

	if d.Debtor.Empty() {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, d.Debtor.String())
	}

	if d.Amount.IsNegative() {
		return sdkErr.Wrap(sdkErr.ErrInvalidRequest, "Amount should be positive")
	}

	if d.Creditor.Empty() {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, (d.Creditor.String()))
	}

	return nil
}

func (d Debt) String() string {
	return strings.TrimSpace(fmt.Sprintf(`ID: %s
                Debtor: %s
                Amount: %s
                Creditor: %s`,
		d.ID,
		d.Debtor,
		d.Amount,
		d.Creditor))
}