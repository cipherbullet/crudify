package main

import (
	"log"
	"github.com/gorilla/mux"
)



func (s *Server) SetupRouter() *mux.Router {
	router := mux.NewRouter()

	if s == nil {
		log.Fatal("Server instance is nil")
	}

	routerApiV1 := router.PathPrefix("/v1/").Subrouter()

	accounts := routerApiV1.PathPrefix("/accounts").Subrouter()
	accounts.HandleFunc("/", MakeHandleFunc(s.indexAccountsHandler)).Methods("GET")
	accounts.HandleFunc("/{id}", MakeHandleFunc(s.showAccountHandler)).Methods("GET")
	accounts.HandleFunc("/", MakeHandleFunc(s.createAccountHandler)).Methods("POST")
	// accounts.HandleFunc("/{id}", MakeHandleFunc(s.UpdateAccountHandler)).Methods("PUT")
	accounts.HandleFunc("/{id}", MakeHandleFunc(s.deleteAccountHandler)).Methods("DELETE")

	return router
}