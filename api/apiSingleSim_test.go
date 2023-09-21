package api

import (
	"fmt"
	"net/http/httptest"

	// sg "github.com/sebastianring/simulationgame"
	"testing"
)

func APIServer_GetHandleSingleSimulation(t *testing.T) {
	fmt.Println("Testing a GET method to api url /api/new_sim which generates a new simulation and returns a Json file with the results per round.")

	// initServer()

	req := httptest.NewRequest("GET", "/api/new_single_sim", nil)

	rr := httptest.NewRecorder()

	s := NewAPIServer(":8080")

	s.HandleSingleSimulation(rr, req)

	fmt.Println(rr.Body.String())
}
