package cli

import (
	"bufio"
	"fmt"

	"github.com/spoto/lending/x/lending/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"
)

func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		createDebtCmd(cdc),
		payDebtCmd(cdc),
		changeDebtCmd(cdc),
	)

	return txCmd
}

func changeDebtCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "change [ID] [amount]",
		Short: "Change an amount (less than the original) for the debt",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return changeDebtCmdFunc(cmd, args, cdc)
		},
	}

	cmd = flags.PostCommands(cmd)[0]

	return cmd
}

func changeDebtCmdFunc(cmd *cobra.Command, args []string, cdc *codec.Codec) error {
	inBuf := bufio.NewReader(cmd.InOrStdin())
	cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
	txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

	ID := args[0]
	creditor := cliCtx.GetFromAddress()
	amount, err := sdk.ParseCoin(args[1])
	if err != nil {
		return err
	}

	msg := types.NewMsgChangeDebt(ID, amount, creditor)

	if err := msg.ValidateBasic(); err != nil {
		return err
	}

	return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
}

func payDebtCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pay [ID] [amount]",
		Short: "Pay an amount for the debt",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return payDebtCmdFunc(cmd, args, cdc)
		},
	}

	cmd = flags.PostCommands(cmd)[0]

	return cmd
}

func payDebtCmdFunc(cmd *cobra.Command, args []string, cdc *codec.Codec) error {
	inBuf := bufio.NewReader(cmd.InOrStdin())
	cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
	txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

	ID := args[0]
	debtor := cliCtx.GetFromAddress()
	amount, err := sdk.ParseCoin(args[1])
	if err != nil {
		return err
	}

	msg := types.NewMsgPayDebt(ID, amount, debtor)

	if err := msg.ValidateBasic(); err != nil {
		return err
	}

	return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
}

func createDebtCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [ID] [amount] [creditor]",
		Short: "Creates a debt that should be collected",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return createDebtCmdFunc(cmd, args, cdc)
		},
	}

	cmd = flags.PostCommands(cmd)[0]

	return cmd
}

func createDebtCmdFunc(cmd *cobra.Command, args []string, cdc *codec.Codec) error {
	inBuf := bufio.NewReader(cmd.InOrStdin())
	cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
	txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

	ID := args[0]
	creditor := cliCtx.GetFromAddress()
	amount, err := sdk.ParseCoin(args[1])
	if err != nil {
		return err
	}
	debtor, err := sdk.AccAddressFromBech32(args[2])
	if err != nil {
		return err
	}

	msg := types.MsgCreateDebt(types.Debt{
		ID:       ID,
		Debtor:   debtor,
		Amount:   amount,
		Creditor: creditor,
	})
	if err := msg.ValidateBasic(); err != nil {
		return err
	}

	return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
}