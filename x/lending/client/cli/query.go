package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	"github.com/spoto/lending/x/lending/types"
)

func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		getAllDebts(cdc),
		getDebtorDebts(cdc),
		getCreditorDebts(cdc),
	)

	return cmd
}

func getDebtorDebts(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "get-debtor-debts [user-address]",
		Short: "Get all the debts where address is debtor",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getDebtorDebtsFunc(cmd, args, cdc)
		},
	}
}

func getDebtorDebtsFunc(cmd *cobra.Command, args []string, cdc *codec.Codec) error {
	cliCtx := context.NewCLIContext().WithCodec(cdc)

	route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryDebtorDebts, args[0])
	res, _, err := cliCtx.QueryWithData(route, nil)

	if err != nil {
		return err
	}

	fmt.Println(string(res))

	return nil
}

func getCreditorDebts(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "get-creditor-debts [user-address]",
		Short: "Get all the debts where address is creditor",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getCreditorDebtsFunc(cmd, args, cdc)
		},
	}
}

func getCreditorDebtsFunc(cmd *cobra.Command, args []string, cdc *codec.Codec) error {
	cliCtx := context.NewCLIContext().WithCodec(cdc)

	route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryCreditorDebts, args[0])
	res, _, err := cliCtx.QueryWithData(route, nil)

	if err != nil {
		return err
	}

	fmt.Println(string(res))

	return nil
}

func getAllDebts(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "get-debts",
		Short: "Get all the debts",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getAllDebtsFunc(cmd, args, cdc)
		},
	}
}

func getAllDebtsFunc(cmd *cobra.Command, args []string, cdc *codec.Codec) error {
	cliCtx := context.NewCLIContext().WithCodec(cdc)

	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryAllDebts)
	res, _, err := cliCtx.QueryWithData(route, nil)

	if err != nil {
		return err
	}

	fmt.Println(string(res))

	return nil
}