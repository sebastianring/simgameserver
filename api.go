package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type APIServer struct {
	listenAddr string
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

type ApiError struct {
	Error string
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func NewAPIServer(listenAddr string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/api/new_single_sim", makeHTTPHandleFunc(s.HandleSingleSimulation))
	router.HandleFunc("/api/new_multiple_sim_conc", makeHTTPHandleFunc(s.HandleMultipleRandomSimulationsConcurrent))
	router.HandleFunc("/api/new_random_sim", makeHTTPHandleFunc(s.HandleSingleRandomSimulation))
	router.HandleFunc("/new_sim_form", makeHTTPHandleFunc(s.HandleSimForm))
	router.HandleFunc("/api/sim/{id: [0-9a-fA-F-]+}", makeHTTPHandleFunc(s.HandleSims))

	log.Println("API server started running on port", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) HandleSingleSimulation(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.newSingleSimulation(w, r)
	}

	return fmt.Errorf("Method not allowed, %s", r.Method)
}

func (s *APIServer) HandleMultipleRandomSimulationsConcurrent(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.newMultipleRandomSimulationsConcurrent(w, r)
	}

	return fmt.Errorf("Method not allowed, %s", r.Method)
}

func (s *APIServer) HandleSingleRandomSimulation(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.newRandomSimulation(w, r)
	}

	return fmt.Errorf("Method not allowed, %s", r.Method)
}

func (s *APIServer) HandleSimForm(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.getSimulationForm(w, r)
	} else if r.Method == "POST" {
		return s.newSingleSimulation(w, r)
	}

	return fmt.Errorf("Method not allowed, %s", r.Method)
}

func (s *APIServer) HandleSims(w http.ResponseWriter, r *http.Request) error {
	log.Println("testing..")
	if r.Method == "GET" {
		return s.getBoardFromDb(w, r)
	} else if r.Method == "DELETE" {
		// BELOW FUNCTION IS NOT WORKING YET
		return s.delBoardFromDb(w, r)
	}

	return fmt.Errorf("Method not allowed, %s", r.Method)
}
