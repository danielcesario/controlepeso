package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/danielcesario/controlepeso/internal/controlepeso"
)

type Service interface {
	CreateEntry(entry controlepeso.Entry) error
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
	fmt.Println("HandleCreate")

	entry := controlepeso.Entry{
		UserId: 1,
		Weight: 110.5,
		Date:   "2022-09-10 00:25:00",
	}
	a.Service.CreateEntry(entry)
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}
