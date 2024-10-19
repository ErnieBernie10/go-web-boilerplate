package pkg

import (
	"encoding/json"
	"net/http"
)

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
