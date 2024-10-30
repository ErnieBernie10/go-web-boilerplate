package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"framer/internal/pkg"
	"net/http"
	"strings"
)

func WriteJSONError(w http.ResponseWriter, statusCode int, err error) {
	errorMessages := strings.Split(err.Error(), "\n")

	errorResponse := ErrorResponseDto{
		Errors: []string{},
	}
	if len(errorMessages)%2 == 1 {
		errorResponse.Errors = append(errorResponse.Errors, errorMessages[0])
	} else {
		for i := 0; i < len(errorMessages); i += 2 {
			errorResponse.Errors = append(errorResponse.Errors, fmt.Sprintf("%s: %s", errorMessages[i], errorMessages[i+1]))
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse)
}

func HandleError(r *http.Request, w http.ResponseWriter, err error) {
	if errors.Is(err, pkg.ErrNotFound) {
		WriteJSONError(w, http.StatusNotFound, err)
	} else if errors.Is(err, sql.ErrNoRows) {
		WriteJSONError(w, http.StatusNotFound, err)
	} else if errors.Is(err, pkg.ErrValidation) {
		WriteJSONError(w, http.StatusBadRequest, err)
	} else if errors.Is(err, pkg.ErrUnauthorized) {
		WriteJSONError(w, http.StatusUnauthorized, err)
	} else {
		WriteJSONError(w, http.StatusInternalServerError, err)
	}
}

func WriteCreatedResponse(w http.ResponseWriter, resourcePath string, response CreatedResponseDto) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	location := resourcePath
	w.Header().Set("Location", location)
	json.NewEncoder(w).Encode(response)
}

func WriteUpdatedResponse(w http.ResponseWriter, resourcePath string, response CreatedResponseDto) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	location := resourcePath
	w.Header().Set("Location", location)
	json.NewEncoder(w).Encode(response)
}

func WriteOkResponse(w http.ResponseWriter, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func CreatedResponse(id string) CreatedResponseDto {
	return CreatedResponseDto{Id: id}
}

type CreatedResponseDto struct {
	Id string `json:"id"`
}

type ErrorResponseDto struct {
	Errors []string `json:"errors"`
}
