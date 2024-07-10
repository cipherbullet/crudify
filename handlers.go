package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
)

func (s *Server) indexAccountsHandler(w http.ResponseWriter, r *http.Request) error {
	var accounts []Account

	ctx := r.Context()
	if err := s.Storage.All(ctx, &accounts); err != nil {
		log.Fatal(err)
	}

	return WriteJson(w, http.StatusOK, accounts)
}

func (s *Server) showAccountHandler(w http.ResponseWriter, r *http.Request) error {
	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return err
	}

	var account Account
	ctx := r.Context()
	err = s.Storage.Find(ctx, &account, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return err
	}

	return WriteJson(w, http.StatusOK, account)
}

func (s *Server) createAccountHandler(w http.ResponseWriter, r *http.Request) error {

	var account Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		return err
	}

	ctx := r.Context()
	err := s.Storage.Create(ctx, &account)
	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusCreated, account)
}

func (s *Server) deleteAccountHandler(w http.ResponseWriter, r *http.Request) error {
    idStr := mux.Vars(r)["id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid account ID", http.StatusBadRequest)
        return err
    }
    ctx := r.Context()
    err = s.Storage.Delete(ctx, &Account{Id: id})
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return err
    }
    return WriteJson(w, http.StatusOK, nil)
}