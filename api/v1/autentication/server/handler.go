package server

import (
	"net/http"
)

func login(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
