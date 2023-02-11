package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/models"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/utils"
	"github.com/gorilla/mux"
	"net/http"
)

// swagger:parameters getTransactions
type _ struct {
	// Address
	// in:path
	Address string `json:"address"`
}

// swagger:model GetTransactionsResp
type GetTransactionsResp struct {
	Transactions []*models.Transaction `json:"transactions"`
}

// swagger:operation GET /get_transactions/{address} getTransactions
// ---
// summary: Get list of transaction by address that already listening
// description: Returns all history of transactions for a given address since subscribe until memory storage is cleaned.
// parameters:
// - name: address
//   in: path
//   description: Ethereum address
//   type: string
//   required: true
// responses:
//   200:
//     description: A list of transactions for the specified address
//     schema:
//       $ref: "#/definitions/GetTransactionsResp"

func (h *Handler) getTransactions(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	vars := mux.Vars(r)
	address, ok := vars["address"]
	if !ok {
		h.sendErrResponse(w, errors.New("address is not provided"), http.StatusBadRequest)
	}

	address = utils.ClearString(address)

	address = utils.ClearString(address)

	transactions, err := h.parser.GetTransactions(ctx, address)
	if err != nil {
		h.sendErrResponse(w, err, http.StatusBadRequest)
		return
	}

	resp := GetTransactionsResp{Transactions: transactions}

	respRaw, err := json.Marshal(resp)
	if err != nil {
		h.sendErrResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.sendOKResponse(w, respRaw)
}
