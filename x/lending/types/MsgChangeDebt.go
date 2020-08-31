package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgChangeDebt{}

type MsgChangeDebt struct {
	ID       string         `json:"id"`
	Amount   sdk.Coin       `json:"amount"`
	Creditor sdk.AccAddress `json:"creditor"`
}

func NewMsgChangeDebt(id string, amount sdk.Coin, creditor sdk.AccAddress) MsgChangeDebt {
	return MsgChangeDebt{
		ID:       id,
		Amount:   amount,
		Creditor: creditor,
	}
}

const PayChangeConst = "ChangeDebt"

func (msg MsgChangeDebt) Route() string { return RouterKey }
func (msg MsgChangeDebt) Type() string  { return PayChangeConst }
func (msg MsgChangeDebt) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creditor}
}
func (msg MsgChangeDebt) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}
func (msg MsgChangeDebt) ValidateBasic() error {
	if msg.ID == "" {
		return sdkErr.Wrap(sdkErr.ErrInvalidRequest, "ID can't be empty")
	}

	if msg.Amount.IsNegative() {
		return sdkErr.Wrap(sdkErr.ErrInvalidRequest, "Amount should be positive")
	}

	if msg.Creditor.Empty() {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, msg.Creditor.String())
	}

	return nil
}
