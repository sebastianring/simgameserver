package main

import (
	"encoding/json"
	// "errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/sebastianring/simulationgame"
)

var parameterRules map[string]*Rule

type Rule struct {
	StandardValue any
	MinVal        int
	MaxVal        int
	ErrorMsg      string
	FinalValue    any
}

func main() {
	initServer()
	runServer()
}

func initServer() {
	rowsRule := Rule{
		StandardValue: 100,
		MinVal:        5,
		MaxVal:        200,
		ErrorMsg:      "Invalid rows value, value should be between 5-200.",
	}

	colsRule := Rule{
		StandardValue: 100,
		MinVal:        5,
		MaxVal:        200,
		ErrorMsg:      "Invalid cols value, value should be between 5-200.",
	}

	drawRule := Rule{
		StandardValue: true,
		ErrorMsg:      "Invalid value for draw parameters, must be either true or false.",
	}

	foodsRule := Rule{
		StandardValue: 75,
		MinVal:        1,
		MaxVal:        150,
		ErrorMsg:      "Invalid value for foods parameter, value should be between 1-150.",
	}

	creature1Rule := Rule{
		StandardValue: 10,
		MinVal:        0,
		MaxVal:        50,
		ErrorMsg:      "Invalid value for creature1, should be between 0-50",
	}

	creature2Rule := Rule{
		StandardValue: 10,
		MinVal:        0,
		MaxVal:        50,
		ErrorMsg:      "Invalid value for creature2, should be between 0-50",
	}

	parameterRules = make(map[string]*Rule, 1)

	parameterRules["rows"] = &rowsRule
	parameterRules["cols"] = &colsRule
	parameterRules["draw"] = &drawRule
	parameterRules["foods"] = &foodsRule
	parameterRules["creature1"] = &creature1Rule
	parameterRules["creature2"] = &creature2Rule
}

func runServer() {
	http.HandleFunc("/api/new_sim", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Someone triggered a new simulation")
		sc := simulationgame.SimulationConfig{}

		for key, rule := range parameterRules {
			parameter := r.URL.Query().Get(key)

			fmt.Println(parameter)

			if len(parameter) == 0 {
				rule.FinalValue = rule.StandardValue
			} else {
				_, ok := rule.StandardValue.(bool)

				if ok {
					if parameter == "true" {
						rule.FinalValue = true

					} else if parameter == "false" {
						rule.FinalValue = false

					} else {
						http.Error(w, rule.ErrorMsg, http.StatusBadRequest)

					}
				} else {
					_, ok := rule.StandardValue.(int)

					if ok {
						fmt.Println("Its an int!")
						value, err := strconv.Atoi(parameter)

						if err != nil {
							http.Error(w, "For value: "+key+" the value was not an int, as expected.", http.StatusBadRequest)
						}

						if value < rule.MinVal || value > rule.MaxVal {
							http.Error(w, rule.ErrorMsg, http.StatusBadRequest)
						}

						rule.FinalValue = value
					}
				}
			}

			cols, ok := parameterRules["cols"].FinalValue.(int)

			if ok {
				sc.Cols = cols
			}

			rows, ok := parameterRules["rows"].FinalValue.(int)

			if ok {
				sc.Rows = rows
			}

			draw, ok := parameterRules["draw"].FinalValue.(bool)

			if ok {
				sc.Draw = draw
			}

			foods, ok := parameterRules["foods"].FinalValue.(int)

			if ok {
				sc.Foods = foods
			}

			creature1, ok := parameterRules["creature1"].FinalValue.(uint)

			if ok {
				sc.Creature1 = creature1
			}

			creature2, ok := parameterRules["creature1"].FinalValue.(uint)

			if ok {
				sc.Creature2 = creature2
			}
		}

		fmt.Println("Cols: " + strconv.Itoa(sc.Cols))
		fmt.Println("Rows: " + strconv.Itoa(sc.Rows))

		resultBoard, err := simulationgame.RunSimulation(&sc)

		if err != nil {
			fmt.Println(err.Error())
			panic(err.Error())
		}

		jsonBytes, err := json.Marshal(resultBoard)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println(err.Error())
			panic(err.Error())
		}

		fmt.Println("Respond with json data: \n" + string(jsonBytes))

		w.Header().Set("Content-type", "application/json")
		w.Write(jsonBytes)
	})

	fmt.Println("Server running at port 8080")
	http.ListenAndServe(":8080", nil)
}
