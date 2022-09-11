package main

import (
	"log"
	"net/http"
	"time"

	"github.com/danielcesario/controlepeso/cmd/handler"
	"github.com/danielcesario/controlepeso/internal/controlepeso"
	"github.com/gorilla/mux"
)

func main() {
	db := controlepeso.InitializeDB("adm_controlepeso", "controlepeso", "controlepeso")
	repository := controlepeso.NewPGRepository(db)
	service := controlepeso.NewService(repository)
	apiHandler := handler.NewHandler(service)

	router := mux.NewRouter()
	router.HandleFunc("/entries", apiHandler.HandleCreateEntry).Methods("POST")
	router.HandleFunc("/entries", apiHandler.HandleListEntries).Methods("GET")
	runServer(router)
}

func runServer(router *mux.Router) {
	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
