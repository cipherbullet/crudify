package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	Bind   		string
	Port   		string
	Storage 	StorageInterface
}

// func NewServer(bind string, port string, storage StorageInterface) *Server {
func NewServer(bind string, port string, storage StorageInterface) *Server {
	return &Server{
		Bind:    	bind,
		Port:    	port,
		Storage: 	storage,
		// Settings: 	&ServerSettings{},
	}
}

func (s *Server) Start() {
	router := s.SetupRouter()

	log.Println("Starting server on :", s.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", s.Bind, s.Port), router))

}


type apiFunc func(http.ResponseWriter, *http.Request) error

func MakeHandleFunc(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// WriteJson writes a JSON response
func WriteJson(w http.ResponseWriter, status int, payload interface{}) error {
    response, err := json.Marshal(payload)
    if err != nil {
        return err
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    w.Write(response)
    return nil
}