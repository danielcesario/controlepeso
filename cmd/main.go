package main

import (
	"log"
	"net/http"
	"time"

	"github.com/danielcesario/entry/cmd/handler"
	"github.com/danielcesario/entry/internal/entry"
	"github.com/gorilla/mux"
)

func main() {
	// db := entry.InitializeDB("adm_controlepeso", "controlepeso", "localhost", "controlepeso")
	db := entry.InitializeDB("tgjetrdu", "uLQTuIwOwGwIQRKSK2l8H9qICgz7qjS4", "jelani.db.elephantsql.com", "tgjetrdu")
	repository := entry.NewPGRepository(db)
	service := entry.NewService(repository)
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
