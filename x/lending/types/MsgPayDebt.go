package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgPayDebt{}

type MsgPayDebt struct {
	ID     string         `json:"id"`
	Amount sdk.Coin       `json:"amount"`
	Debtor sdk.AccAddress `json:"debtor"`
}

func NewMsgPayDebt(id string, amount sdk.Coin, debtor sdk.AccAddress) MsgPayDebt {
	return MsgPayDebt{
		ID:     id,
		Amount: amount,
		Debtor: debtor,
	}
}

const PayDebtConst = "PayDebt"

func (msg MsgPayDebt) Route() string { return RouterKey }
func (msg MsgPayDebt) Type() string  { return PayDebtConst }
func (msg MsgPayDebt) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Debtor}
}
func (msg MsgPayDebt) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}
func (msg MsgPayDebt) ValidateBasic() error {
	if msg.ID == "" {
		return sdkErr.Wrap(sdkErr.ErrInvalidRequest, "ID can't be empty")
	}

	if msg.Debtor.Empty() {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, msg.Debtor.String())
	}

	if msg.Amount.IsNegative() {
		return sdkErr.Wrap(sdkErr.ErrInvalidRequest, "Amount should be positive")
	}

	return nil
}
