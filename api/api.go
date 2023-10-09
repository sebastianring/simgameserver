package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	sc "github.com/sebastianring/simgameserver/simconfig"
	"log"
	"net/http"
	"os"
	"time"
)

func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/api/new_single_sim", makeHTTPHandleFunc(s.HandleSingleSimulation))
	router.HandleFunc("/api/new_multiple_sim/{iterations:[1-9][0-9]*}", makeHTTPHandleFunc(s.HandleMultipleRandomSimulationsConcurrent))
	router.HandleFunc("/api/new_random_sim", makeHTTPHandleFunc(s.HandleSingleRandomSimulation))
	router.HandleFunc("/new_sim_form", makeHTTPHandleFunc(s.HandleSimForm))
	router.HandleFunc("/api/sim/{id:[0-9a-fA-F-]+}", makeHTTPHandleFunc(s.HandleSims))

	log.Println("API server started running on port", s.listenAddr)
	err := http.ListenAndServe(s.listenAddr, router)

	if err != nil {
		log.Println(err)
	}
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
	if r.Method == "GET" {
		return s.getBoardFromDb(w, r)
	} else if r.Method == "DELETE" {
		return s.delBoardFromDb(w, r)
	}

	return fmt.Errorf("Method not allowed, %s", r.Method)
}

type APIServer struct {
	listenAddr string
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

type ApiError struct {
	Error string `json:"error"`
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10000*time.Millisecond)
		defer cancel()

		done := make(chan bool)
		err := errors.New("")

		go func() {
			err = f(w, r)
			close(done)

			if err != nil {
				WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
			}
		}()

		select {
		case <-ctx.Done():
			close(done)

			WriteJSON(w, http.StatusGatewayTimeout, ApiError{Error: "Operation timed out."})

		case <-done:

		}
	}
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v ", token.Header["alg"])
		}

		return []byte(secret), nil
	})
}

func WithJWTAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Calling middleware JWT Auth.")

		tokenString := r.Header.Get("x-jwt-token")

		_, err := validateJWT(tokenString)

		if err != nil {
			WriteJSON(w, http.StatusForbidden, ApiError{Error: "Issue validating JWT: " + err.Error()})
			return
		}

		handlerFunc(w, r)
	}
}

func NewAPIServer(listenAddr string) *APIServer {
	sc.InitRules()
	return &APIServer{
		listenAddr: listenAddr,
	}
}
