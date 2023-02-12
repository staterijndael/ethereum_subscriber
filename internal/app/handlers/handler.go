// Ethereum Subscriber Api:
// version: 0.0.1
// title: Ethereum Subscriber Api
// Schemes: http, https
// Host: localhost:8080
// BasePath: /
// Produces:
// - application/json
//
// swagger:meta

package handlers

import (
	"context"
	"github.com/bluntenpassant/ethereum_subscriber/config"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/models"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

type Parser interface {
	GetCurrentBlock(ctx context.Context) (uint64, error)
	GetTransactions(ctx context.Context, address string) ([]*models.Transaction, error)
	Subscribe(ctx context.Context, address string) error
}

type Handler struct {
	parser Parser
}

func NewHandler(parser Parser) *Handler {
	return &Handler{
		parser: parser,
	}
}

func (h *Handler) Start(httpConfig config.Http) {
	r := mux.NewRouter()
	r.HandleFunc("/subscribe/{address}", h.subscribe)
	r.HandleFunc("/get_current_block", h.getCurrentBlock)
	r.HandleFunc("/get_transactions/{address}", h.getTransactions)

	// This will serve files under http://localhost:8000/static/<filename>
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./internal/app/handlers"))))

	opts := middleware.SwaggerUIOpts{SpecURL: "/static/swagger.yaml"}
	sh := middleware.SwaggerUI(opts, nil)
	r.Handle("/docs", sh)

	srv := &http.Server{
		Handler: r,
		Addr:    httpConfig.Host + ":" + httpConfig.Port,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func (h *Handler) sendOKResponse(w http.ResponseWriter, resp []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (h *Handler) sendErrResponse(w http.ResponseWriter, err error, status int) {
	w.WriteHeader(status)
	w.Write([]byte(err.Error()))
}
