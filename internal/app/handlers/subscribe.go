package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/utils"
	"github.com/gorilla/mux"
	"net/http"
)

// swagger:model SubscribeResp
type SubscribeResp struct {
	IsOK bool `json:"is_ok"`
}

// swagger:operation GET /subscribe/{address} subscribe
// ---
// summary: Subscribe address for a listening new transactions
// description: Set up listening for address. Transactions available for getting through /get_transactions/{address} method
// parameters:
// - name: address
//   in: path
//   description: Ethereum address
//   type: string
//   required: true
// responses:
//   200:
//     description: Is everything ok status
//     schema:
//       $ref: "#/definitions/SubscribeResp"

func (h *Handler) subscribe(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	vars := mux.Vars(r)
	address, ok := vars["address"]
	if !ok {
		h.sendErrResponse(w, errors.New("address is not provided"), http.StatusBadRequest)
	}

	address = utils.ClearString(address)

	err := h.parser.Subscribe(ctx, address)
	if err != nil {
		h.sendErrResponse(w, err, http.StatusBadRequest)
		return
	}

	resp := SubscribeResp{IsOK: true}

	respRaw, err := json.Marshal(resp)
	if err != nil {
		h.sendErrResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.sendOKResponse(w, respRaw)
}
