package rest

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/gorilla/mux"
	"github.com/spoto/lending/x/lending/types"
	"net/http"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(fmt.Sprintf("/%s/createdebt", types.ModuleName), createDebtFn(cliCtx)).Methods("POST")
}

type createDebtRequest struct {
	BaseReq rest.BaseReq `json:"base_req"`
	ID     string         `json:"ID"`
	Debtor sdk.AccAddress `json:"debtor"`
	Amount   sdk.Coin       `json:"amount"`
	Creditor sdk.AccAddress `json:"creditor"`
}

func createDebtFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createDebtRequest

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()

		if !baseReq.ValidateBasic(w) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid request")
			return
		}

		// create the message
		msg := types.MsgCreateDebt {
			ID:       req.ID,
			Amount:   req.Amount,
			Creditor: req.Creditor,
			Debtor:   req.Debtor,
		}

		err := msg.ValidateBasic()

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}