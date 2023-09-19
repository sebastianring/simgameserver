package main

import (
	"fmt"
	// "net/http/httptest"

	// sg "github.com/sebastianring/simulationgame"
	"testing"
)

func TestNewSingleSimulationFromApi(t *testing.T) {
	fmt.Println("Testing a GET method to api url /api/new_sim which generates a new simulation and returns a Json file with the results per round.")

	initServer()

	service := NewAPIServer(":8080")
	service.Run()

	// request := httptest.NewRequest("GET", "127.0.0.1/api/single_sim", nil)
	// w := httptest.NewRecorder()

	// if err != nil {
	// 	t.Error("Failed test due to: ", err.Error())
	// 	t.Fatal()
	// }
}
