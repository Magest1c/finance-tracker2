package main

import "log/slog"

type application struct {
	config       config
	logger       *slog.Logger
	users        *userStore
	accounts     *accountsStore
	categories   *CategoryStore
	transactions *transactionStore
}
