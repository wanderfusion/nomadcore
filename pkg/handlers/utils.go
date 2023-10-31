package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/wanderfusion/nomadcore/pkg/errors"

	"github.com/rs/zerolog/log"
)

// writing
type response struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  string      `json:"error,omitempty"`
}

func NewInternalServerError(err error) errors.APIError {
	return errors.APIError{
		Msg:    "Something went wrong...",
		Err:    err,
		Status: http.StatusInternalServerError,
	}
}

func RespondWithData(w http.ResponseWriter, r *http.Request, data interface{}) {
	res := &response{
		Status: http.StatusOK,
		Data:   data,
	}

	json, err := json.Marshal(res)
	if err != nil {
		log.Error().Err(err).Msg("Unable to marshall data json.")
		RespondWithError(w, r, NewInternalServerError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func RespondWithError(w http.ResponseWriter, r *http.Request, e errors.Error) {
	res := &response{
		Status: e.GetStatusCode(),
		Error:  e.GetMsg(),
	}

	json, err := json.Marshal(res)
	if err != nil {
		log.Error().Err(err).Msg("Unable to marshall error json.")
		json = []byte("{\"status\": 200, \"error\": \"Something went wrong\"}")
	}

	log.Error().
		Err(err).
		Interface("status", e.GetStatusCode()).
		Msg(e.GetMsg())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.GetStatusCode())
	w.Write(json)
}

func FromRequest[T any](req *http.Request, v *T) error {
	decoder := json.NewDecoder(req.Body)
	return decoder.Decode(v)
}
