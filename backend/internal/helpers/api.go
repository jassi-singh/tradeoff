package helpers

import (
	"encoding/json"
	"net/http"
)

type CustomError struct {
	Message string `json:"message"`
	status  int    `json:"-"`
}

func (e *CustomError) Error() string {
	return e.Message
}

func NewCustomError(message string, status int) *CustomError {
	return &CustomError{
		Message: message,
		status:  status,
	}
}

func RespondWithError(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case *CustomError:
		RespondWithJSON(w, e.status, map[string]string{"error": e.Message})
	default:
		RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	}
}

func RespondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	response, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

