package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// RespondWithError helper to reduce code duplication
func RespondWithError(w http.ResponseWriter, status int, err string) {
	w.WriteHeader(status)
	fmt.Fprint(w, err)
}

// RespondWithSuccess helper to reduce code duplication
func RespondWithSuccess(w http.ResponseWriter, status int, body interface{}) {
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(body)
}
