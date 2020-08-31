package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgCreateDebt{}

type MsgCreateDebt Debt

func NewMsgCreateDebt(debt Debt) MsgCreateDebt {
	return MsgCreateDebt(debt) // cast
}

const CreateDebtConst = "CreateDebt"

func (msg MsgCreateDebt) Route() string { return RouterKey }
func (msg MsgCreateDebt) Type() string  { return CreateDebtConst }
func (msg MsgCreateDebt) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creditor}
}
func (msg MsgCreateDebt) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}
func (msg MsgCreateDebt) ValidateBasic() error {
	return Debt(msg).Validate() // delegation to Debt
}
