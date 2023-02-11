package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/utils"
	"net/http"
	"strings"
)

type SubscribeResp struct {
	IsOK bool `json:"is_ok"`
}

func (h *Handler) subscribe(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) == 0 {
		h.sendErrResponse(w, errors.New("incorrect address url"), http.StatusBadRequest)
		return
	}
	address := parts[len(parts)-1]

	address = utils.ClearString(address)

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
