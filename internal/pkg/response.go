package pkg

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func WriteCreatedResponse(w http.ResponseWriter, id uuid.UUID, resourcePath string) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	location := resourcePath
	location = strings.Replace(location, "{id}", id.String(), 1)
	w.Header().Set("Location", location)
	response := map[string]interface{}{
		"message": "Resource created successfully",
		"id":      id,
	}
	json.NewEncoder(w).Encode(response)
}

func WriteUpdatedResponse(w http.ResponseWriter, id uuid.UUID, resourcePath string) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	location := resourcePath
	location = strings.Replace(location, "{id}", id.String(), 1)
	w.Header().Set("Location", location)
	response := map[string]interface{}{
		"message": "Resource updated successfully",
		"id":      id,
	}
	json.NewEncoder(w).Encode(response)
}
