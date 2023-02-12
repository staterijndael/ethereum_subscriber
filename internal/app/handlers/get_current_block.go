package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

// swagger:model GetCurrentBlockResp
type GetCurrentBlockResp struct {
	// in: uint64
	CurrentBlock uint64 `json:"current_block"`
}

// swagger:route GET /get_current_block getCurrentBlock
// Returns last parsed block between all transactions.
// NOTE: Current block is not attached to last parsed transaction
// and indicates only block number that was handled by internal parser
//
// responses:
//
//	200: GetCurrentBlockResp
func (h *Handler) getCurrentBlock(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	currentBlock, err := h.parser.GetCurrentBlock(ctx)
	if err != nil {
		h.sendErrResponse(w, err, http.StatusBadRequest)
		return
	}

	if currentBlock == 0 {
		h.sendErrResponse(w, errors.New("current block is not parsed yet"), http.StatusBadRequest)
		return
	}

	resp := GetCurrentBlockResp{CurrentBlock: currentBlock}

	respRaw, err := json.Marshal(resp)
	if err != nil {
		h.sendErrResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.sendOKResponse(w, respRaw)
}
