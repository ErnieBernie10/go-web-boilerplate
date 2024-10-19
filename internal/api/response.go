package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"framer/internal/core"
	"net/http"
	"os"
)

func WriteJSONError(w http.ResponseWriter, statusCode int, err error) {
	var errorResponse interface{}
	if errs, ok := err.(interface{ Unwrap() []error }); ok {
		var messages []string
		for _, e := range errs.Unwrap() {
			messages = append(messages, e.Error())
		}
		errorResponse = map[string][]string{"errors": messages}
	} else {
		errorResponse = map[string]string{"error": err.Error()}
	}
	errorJSON, _ := json.Marshal(errorResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(errorJSON)
}

func HandleError(r *http.Request, w http.ResponseWriter, err error, statusCode int) {

	switch statusCode {
	case http.StatusInternalServerError:
		GetLogger(r).Error(err.Error())
		if os.Getenv("APP_ENV") == string(core.Development) {
			WriteJSONError(w, statusCode, err)
		} else {
			WriteJSONError(w, statusCode, errors.New("internal server error"))
		}
	default:
		WriteJSONError(w, statusCode, err)
	}
}

func HandleDbError(r *http.Request, w http.ResponseWriter, err error) {
	if errors.Is(err, sql.ErrNoRows) {
		HandleError(r, w, err, http.StatusNotFound)
	} else {
		HandleError(r, w, err, http.StatusInternalServerError)
	}
}

func WriteCreatedResponse(w http.ResponseWriter, resourcePath string, response CreatedResponseDto) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	location := resourcePath
	w.Header().Set("Location", location)
	json.NewEncoder(w).Encode(response)
}

func WriteUpdatedResponse(w http.ResponseWriter, resourcePath string, response CreatedResponseDto) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	location := resourcePath
	w.Header().Set("Location", location)
	json.NewEncoder(w).Encode(response)
}

func CreatedResponse(id string) CreatedResponseDto {
	return CreatedResponseDto{Id: id}
}

type CreatedResponseDto struct {
	Id string `json:"id"`
}
