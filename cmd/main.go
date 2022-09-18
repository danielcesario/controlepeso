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
	entryRepository := entry.NewPGRepository(db)
	entryService := entry.NewService(entryRepository)
	entryHandler := handler.NewHandler(entryService)

	router := mux.NewRouter()
	router.Use(mux.CORSMethodMiddleware(router))
	router.Use(loggingMiddleware)
	router.HandleFunc("/entries", entryHandler.HandleCreateEntry).Methods("POST")
	router.HandleFunc("/entries", entryHandler.HandleListEntries).Methods("GET")
	router.HandleFunc("/entries/{id}", entryHandler.HandleGetEntry).Methods("GET")
	router.HandleFunc("/entries/{id}", entryHandler.HandleDeleteEntry).Methods("DELETE")
	router.HandleFunc("/entries/{id}", entryHandler.HandleUpdateEntry).Methods("PUT")
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

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
