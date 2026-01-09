package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func InitRouter() *mux.Router {
	r := mux.NewRouter()
	// CORS handling is only required during development
	headersOk := handlers.AllowedHeaders([]string{
		"Accept",
		"Accept-Language",
		"Content-Type",
		"Content-Length",
		"Authorization",
		"Origin",
	})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "DELETE", "PUT"})
	corsMiddleware := handlers.CORS(headersOk, originsOk, methodsOk)

	r.NotFoundHandler = corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&ApiResponse{Message: "Not found", Data: nil})
	}))

	r.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Use(corsMiddleware)
	return r
}

func RegisterRoutes(a *App) {
	r := a.Router
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "This is Markat - a mail client\n")
	}).Methods("GET")
	r.HandleFunc("/api/mb/list", a.MailboxListHandler).Methods("GET")
	r.HandleFunc("/api/mb/select", a.SelectMailBoxHandler).Methods("PUT")
	r.HandleFunc("/api/mb/envelopes", a.FetchEnvelopes).Methods("GET")
}
