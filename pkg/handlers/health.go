package handlers

import (
	"net/http"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	data := "Preparing my travel docs..."
	RespondWithData(w, r, data)
}
