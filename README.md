# Finance Tracker API

A portfolio REST API for tracking personal income and expenses, built while learning Go.

## Current Features

- User registration and login with bcrypt password hashing
- Account and category management
- Transaction CRUD, filtering, sorting, and pagination
- Monthly summaries by account
- JSON error responses and structured logging

## Tech Stack

- Go 1.26
- chi router
- slog logging
- bcrypt password hashing

The project currently uses in-memory storage. PostgreSQL, JWT authentication, tests, and Docker support are planned next.

## Run Locally

```powershell
go run ./cmd/api
```

The API listens on `http://localhost:8080` by default. Set `HTTP_ADDR` to use a different address.

```powershell
$env:HTTP_ADDR = ":8090"
go run ./cmd/api
```

## Main Endpoints

- `GET /health`
- `POST /auth/register`
- `POST /auth/login`
- `GET /accounts`
- `POST /accounts`
- `GET /categories`
- `POST /categories`
- `GET /transactions`
- `POST /transactions`
- `PATCH /transactions/{id}`
- `DELETE /transactions/{id}`
- `GET /reports/monthly`

