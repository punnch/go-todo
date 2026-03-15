package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type ErrorDTO struct {
	Message string `json:"message"`
}

func NewErrorDTO(err error) ErrorDTO {
	return ErrorDTO{
		Message: err.Error(),
	}
}

func (e ErrorDTO) toString() string {
	b, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		panic(err)
	}

	return string(b)
}

func errorCompareJSON(w http.ResponseWriter, err, target error, code int) {
	errorDTO := NewErrorDTO(err)

	if errors.Is(err, target) {
		http.Error(w, errorDTO.toString(), code)
	} else {
		http.Error(w, errorDTO.toString(), http.StatusInternalServerError)
	}
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		fmt.Println("failed to write http response body:", err)
	}
}
