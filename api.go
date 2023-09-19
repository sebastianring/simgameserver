package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	sg "github.com/sebastianring/simulationgame"
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
	router.HandleFunc("/api/single_sim", makeHTTPHandleFunc(s.HandleSingleSimulation))

	log.Println("API server started running on port: ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) HandleSingleSimulation(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.getNewSingleSimulation(w, r)
	}

	return fmt.Errorf("Method not allowed, %s", r.Method)
}

func (s *APIServer) getNewSingleSimulation(w http.ResponseWriter, r *http.Request) error {
	sc, err := getSimulationConfigFromUrlValues(r.URL.Query())

	if err != nil {
		fmt.Println("Error occured during getting simulation config from URL values: ", err)
		return err
	}

	fmt.Println("Starting simulation with config: ")
	fmt.Println("sc.creature1: " + strconv.Itoa(int(sc.Creature1)))

	resultBoard, err := sg.RunSimulation(sc)

	if err != nil {
		fmt.Println("Error occured during running the simulation: ", err)
		return err
	}

	roundData, err := getRoundData(resultBoard, AliveAtEnd)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, roundData)
}
