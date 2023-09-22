package api_test

import (
	"fmt"
	"github.com/sebastianring/simgameserver/api"
	sc "github.com/sebastianring/simgameserver/simconfig"
	"net/http/httptest"

	// sg "github.com/sebastianring/simulationgame"
	"testing"
)

func TestAPIServer_GetHandleSingleSimulation(t *testing.T) {
	t.Setenv("sim_game", "valmet865")
	fmt.Println("Testing a GET method to api url /api/new_sim which generates a new simulation and returns a Json file with the results per round.")

	req := httptest.NewRequest("GET", "/api/new_single_sim", nil)

	rr := httptest.NewRecorder()

	s := api.NewAPIServer(":8080")

	sc.InitRules()

	s.HandleSingleSimulation(rr, req)

	fmt.Println(rr.Body.String())
}
