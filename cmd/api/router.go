package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) router() http.Handler {
	router := chi.NewRouter()

	router.Get("/health", healthHandler(app.logger))

	router.Post("/echo", echoHandler(app.logger))

	router.Post("/auth/register", authRegisterHandler(app.logger, app.users))

	router.Post("/auth/login", authLoginHandler(app.logger, app.users))

	router.Post("/accounts", createAccountHandler(app.logger, app.accounts, app.users))

	router.Post("/categories", createCategoryHandler(app.logger, app.categories, app.users))

	router.Get("/categories", getCategoriesHandler(app.logger, app.categories, app.users))

	router.Get("/accounts", getAccountsHandler(app.logger, app.accounts, app.users))

	router.Post("/transactions", createTransactionHandler(app.logger, app.transactions, app.users, app.accounts, app.categories))

	router.Get("/transactions", getTransactionsHandler(app.logger, app.transactions, app.users))

	router.Patch("/transactions/{id}", updateTransactionHandler(app.logger, app.transactions, app.users))

	router.Delete("/transactions/{id}", deleteTransactionHandler(app.logger, app.transactions, app.users))

	router.Get("/reports/monthly", monthlySummaryHandler(app.logger, app.transactions, app.users, app.accounts))

	return router
}
