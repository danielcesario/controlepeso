package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/danielcesario/entry/internal/entry"
)

type Service interface {
	CreateEntry(entry entry.Entry) (*entry.Entry, error)
	ListEntries(start, count int) ([]entry.Entry, error)
}

type Handler struct {
	Service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		Service: service,
	}
}

func (handler *Handler) HandleCreateEntry(w http.ResponseWriter, r *http.Request) {
	var newEntry entry.Entry
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newEntry); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	createdEntry, err := handler.Service.CreateEntry(newEntry)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	respondWithJSON(w, http.StatusCreated, createdEntry)
}

func (handler *Handler) HandleListEntries(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count < 1 {
		count = 10
	}

	if start < 0 {
		start = 0
	}

	entries, err := handler.Service.ListEntries(start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	respondWithJSON(w, http.StatusOK, entries)
}
