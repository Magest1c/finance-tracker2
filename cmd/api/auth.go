package main

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type authRegisterResponse struct {
	Message string `json:"message"`
	Email   string `json:"email"`
	Name    string `json:"name"`
}
type authRegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}
type authLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type authLoginResponse struct {
	Message string `json:"message"`
	Email   string `json:"email"`
	Name    string `json:"name"`
}

func authRegisterHandler(logger *slog.Logger, store *userStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req authRegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			if writeErr := writeError(w, http.StatusBadRequest, "invalid JSON body"); writeErr != nil {
				logger.Error("failed to encode error response", "error", writeErr)
			}
			logger.Error("failed to decode request body", "error", err)
			return

		}

		email := strings.ToLower(strings.TrimSpace(req.Email))
		name := strings.TrimSpace(req.Name)

		if email == "" || name == "" || req.Password == "" {
			if err := writeError(w, http.StatusBadRequest, "email, password and name are required"); err != nil {
				logger.Error("failed to encode error response", "error", err)
			}
			return

		}
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			if writeErr := writeError(w, http.StatusInternalServerError, "failed to hash password"); writeErr != nil {
				logger.Error("failed to encode error response", "error", writeErr)
			}
			logger.Error("failed to hash password", "error", err)
			return
		}
		user, err := store.Create(email, name, string(hash))
		if errors.Is(err, errUserAlreadyExists) {
			if writeErr := writeError(w, http.StatusConflict, "user with this email already exists"); writeErr != nil {
				logger.Error("failed to encode error response", "error", writeErr)
			}
			return
		}

		if err != nil {
			if writeErr := writeError(w, http.StatusInternalServerError, "internal server error"); writeErr != nil {
				logger.Error("failed to encode error response", "error", writeErr)
			}
			logger.Error("failed to create user", "error", err)
			return
		}
		resp := authRegisterResponse{
			Email:   user.Email,
			Name:    user.Name,
			Message: "user registered",
		}
		if err := writeJSON(w, http.StatusCreated, resp); err != nil {
			logger.Error("failed to encode response", "error", err)
		}
	}
}
func authLoginHandler(logger *slog.Logger, store *userStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req authLoginRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			writeError(w, http.StatusBadRequest, "bad request")
			return
		}
		email := strings.ToLower(strings.TrimSpace(req.Email))

		if email == "" || req.Password == "" {
			writeError(w, http.StatusBadRequest, "пустой запрос")
			return
		}
		user, ok := store.GetByEmail(email)
		if !ok {
			writeError(w, http.StatusUnauthorized, "invalid email or password")
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))

		if err != nil {
			writeError(w, http.StatusUnauthorized, "invalid email or password")
			return
		}
		resp := authLoginResponse{
			Message: "login successful",
			Email:   user.Email,
			Name:    user.Name,
		}
		writeJSON(w, http.StatusOK, resp)

	}
}

func (s *userStore) GetByID(id int64) (User, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, user := range s.userByEmail {
		if user.ID == id {
			return user, true
		}
	}
	return User{}, false
}
