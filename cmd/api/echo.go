package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type echoRequest struct {
	Message string `json:"message"`
}
type echoResponse struct {
	Message string `json:"message"`
}

func echoHandler(logger *slog.Logger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		var req echoRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			if writeErr := writeError(w, http.StatusBadRequest, "invalid JSON body"); writeErr != nil {
				logger.Error("failed to encode error response", "error", writeErr)
			}
			logger.Error("failed to decode request body", "error", err)
			return

		}

		if req.Message == "" {
			if err := writeError(w, http.StatusBadRequest, "message is required"); err != nil {
				logger.Error("failed to encode error response", "error", err)
			}
			return
		}

		resp := echoResponse{
			Message: req.Message,
		}

		if err := writeJSON(w, http.StatusOK, resp); err != nil {
			logger.Error("failed to encode response", "error", err)
		}
	}

}
