package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/shivanshkc/ledgerkeep/src/configs"
	"github.com/shivanshkc/ledgerkeep/src/database"
	"github.com/shivanshkc/ledgerkeep/src/handlers"
	"github.com/shivanshkc/ledgerkeep/src/logger"
	"github.com/shivanshkc/ledgerkeep/src/middlewares"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	conf := configs.Get()
	log := logger.Get()

	// Creating database indexes. This also initiates a connection with the database upon application startup.
	go func() {
		indexData := []mongo.IndexModel{
			{Keys: bson.D{{Key: "account_id", Value: 1}}}, // Ascending B-tree index on "account_id".
			{Keys: bson.D{{Key: "notes", Value: "text"}}}, // Text index on "notes".
		}

		err := database.CreateIndexOnTransactionField(context.Background(), indexData)
		if err != nil {
			panic(err)
		}
	}()

	log.Info(context.Background(),
		&logger.Entry{Payload: fmt.Sprintf("Server listening at: %s", conf.HTTPServer.Addr)})

	// Starting the HTTP server.
	if err := http.ListenAndServe(conf.HTTPServer.Addr, getHandler()); err != nil {
		log.Error(context.Background(),
			&logger.Entry{Payload: fmt.Errorf("failed to start http server: %w", err)})
	}
}

func getHandler() http.Handler {
	router := mux.NewRouter()

	// Attaching global middlewares.
	router.Use(middlewares.Recovery)
	router.Use(middlewares.RequestContext)
	router.Use(middlewares.AccessLogger)
	router.Use(middlewares.CORS)

	// Auth middleware.
	router.Use(middlewares.Auth)

	router.HandleFunc("/api", handlers.BasicHandler).
		Methods(http.MethodGet, http.MethodOptions)

	router.HandleFunc("/api/accounts", handlers.CreateAccountHandler).
		Methods(http.MethodPost, http.MethodOptions)

	router.HandleFunc("/api/accounts", handlers.ListAccountsHandler).
		Methods(http.MethodGet, http.MethodOptions)

	router.HandleFunc("/api/accounts/{account_id}", handlers.UpdateAccountHandler).
		Methods(http.MethodPatch, http.MethodOptions)

	router.HandleFunc("/api/accounts/{account_id}", handlers.DeleteAccountHandler).
		Methods(http.MethodDelete, http.MethodOptions)

	router.HandleFunc("/api/transactions", handlers.CreateTransactionHandler).
		Methods(http.MethodPost, http.MethodOptions)

	router.HandleFunc("/api/transactions/{transaction_id}", handlers.GetTransactionHandler).
		Methods(http.MethodGet, http.MethodOptions)

	router.HandleFunc("/api/transactions", handlers.ListTransactionsHandler).
		Methods(http.MethodGet, http.MethodOptions)

	router.HandleFunc("/api/transactions/{transaction_id}", handlers.UpdateTransactionHandler).
		Methods(http.MethodPatch, http.MethodOptions)

	router.HandleFunc("/api/transactions/{transaction_id}", handlers.DeleteTransactionHandler).
		Methods(http.MethodDelete, http.MethodOptions)

	router.HandleFunc("/api/tags", handlers.CreateTagHandler).
		Methods(http.MethodPost, http.MethodOptions)

	router.HandleFunc("/api/tags", handlers.ListTagsHandler).
		Methods(http.MethodGet, http.MethodOptions)

	router.HandleFunc("/api/tags/{tag_id}", handlers.UpdateTagHandler).
		Methods(http.MethodPatch, http.MethodOptions)

	router.HandleFunc("/api/tags/{tag_id}", handlers.DeleteTagHandler).
		Methods(http.MethodDelete, http.MethodOptions)

	router.HandleFunc("/api/balances/by_category", handlers.GetBalancesByCategoryHandler).
		Methods(http.MethodGet, http.MethodOptions)

	router.HandleFunc("/api/balances/by_time", handlers.GetBalancesByTimeHandler).
		Methods(http.MethodGet, http.MethodOptions)

	router.HandleFunc("/api/balances/by_tags", handlers.GetBalancesByTagsHandler).
		Methods(http.MethodGet, http.MethodOptions)

	return router
}
