package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

func (s *CategoryStore) GetByID(id int64) (Category, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	category, ok := s.categories[id]

	return category, ok
}

type Category struct {
	ID     int64  `json:"id"`
	UserID int64  `json:"user_id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
}
type CategoryStore struct {
	mu         sync.Mutex
	categories map[int64]Category
	nextID     int
}

func newCategoryStore() *CategoryStore {
	return &CategoryStore{
		categories: make(map[int64]Category),
		nextID:     1,
	}
}
func (s *CategoryStore) Create(userID int64, name string, typeCategory string) (Category, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	category := Category{
		ID:     int64(s.nextID),
		UserID: userID,
		Name:   name,
		Type:   typeCategory,
	}
	s.categories[category.ID] = category
	s.nextID++

	return category, nil
}

type createCategoryRequest struct {
	UserID int64  `json:"user_id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
}

func createCategoryHandler(logger *slog.Logger, categoryStore *CategoryStore, userStore *userStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createCategoryRequest
		err := json.NewDecoder(r.Body).Decode(&req)

		if err != nil {
			writeError(w, http.StatusBadRequest, "bad body")
			return
		}
		name := strings.TrimSpace(req.Name)
		categoryType := strings.ToLower(strings.TrimSpace(req.Type))

		if name == "" || categoryType != "income" && categoryType != "expense" || categoryType == "" || req.UserID <= 0 {
			writeError(w, http.StatusBadRequest, "bad body")
			return
		}
		if _, ok := userStore.GetByID(req.UserID); !ok {
			if err := writeError(w, http.StatusBadRequest, "user not found"); err != nil {
				logger.Error("failed to encode error response", "error", err)
			}
			return
		}
		category, err := categoryStore.Create(req.UserID, name, categoryType)

		if err != nil {
			writeError(w, http.StatusInternalServerError, "bad server")
			return
		}
		writeJSON(w, http.StatusCreated, category)

	}
}
func (s *CategoryStore) ListByUserID(userID int64) []Category {
	s.mu.Lock()
	defer s.mu.Unlock()

	categories := make([]Category, 0)

	for _, category := range s.categories {
		if category.UserID == userID {
			categories = append(categories, category)
		}
	}
	return categories
}
func getCategoriesHandler(logger *slog.Logger, categoryStore *CategoryStore, userStore *userStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userIDParam := r.URL.Query().Get("user_id")

		if userIDParam == "" {
			if err := writeError(w, http.StatusBadRequest, "user_id is required"); err != nil {
				logger.Error("failed to encode error response", "error", err)
			}
			return
		}
		userID, err := strconv.ParseInt(userIDParam, 10, 64)

		if err != nil || userID <= 0 {
			if err := writeError(w, http.StatusBadRequest, "user_id must be a positive integer"); err != nil {
				logger.Error("failed to encode error response", "error", err)
			}
			return
		}
		if _, ok := userStore.GetByID(userID); !ok {
			if err := writeError(w, http.StatusBadRequest, "user not found"); err != nil {
				logger.Error("failed to encode error response", "error", err)
			}
			return
		}
		categories := categoryStore.ListByUserID(userID)

		if err := writeJSON(w, http.StatusOK, categories); err != nil {
			logger.Error("failed to encode response", "error", err)
		}
	}
}
