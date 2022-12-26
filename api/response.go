package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-kit/log"
	domain "github.com/kl09/counter-api"
)

func encodeJSONResponse(w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// encodeErrorResp converts errors returned by endpoint.Endpoint,
// request decoder/response encoder (JSON serialization errors, e.g., EOF) into HTTP response.
// Business logic errors are not sent here.
func encodeErrorResp(err error, w http.ResponseWriter, logger log.Logger) {
	logger.Log("error", err)

	if errors.Is(err, domain.ErrInvalidKey) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	// For security reason we shouldn't return detailed error for 5xx errors
	json.NewEncoder(w).Encode(map[string]string{
		"error": "internal error",
	})
}
