package handler

import (
	"encoding/json"
	"net/http"

	"github.com/danielcesario/controlepeso/internal/controlepeso"
)

type Service interface {
	CreateEntry(entry controlepeso.Entry) (*controlepeso.Entry, error)
}

type Handler struct {
	Service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		Service: service,
	}
}

func (a *Handler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var newEntry controlepeso.Entry
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newEntry); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	createdEntry, err := a.Service.CreateEntry(newEntry)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	respondWithJSON(w, http.StatusCreated, createdEntry)
}
