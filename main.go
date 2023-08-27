package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/sebastianring/simulationgame"
	"net/http"
	"strconv"
)

var parameterRules map[string]*Rule

type Rule struct {
	StandardValue any
	MinVal        int
	MaxVal        int
	ErrorMsg      string
	FinalValue    any
}

type DBboard struct {
	Id   uuid.UUID `json:"id"`
	Rows int       `json:"rows"`
	Cols int       `json:"cols"`
}

func main() {
	initServer()
	runServer()
}

func runServer() {
	http.HandleFunc("/api/new_sim", newSimulation)
	http.HandleFunc("/api/sim", getBoardFromDb)

	fmt.Println("Server running at port 8080")
	http.ListenAndServe(":8080", nil)
}

func newSimulation(w http.ResponseWriter, r *http.Request) {
	sc, err := getSimulationConfig(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err.Error())
	}

	fmt.Println("Starting simulation with config: ")
	fmt.Println("sc.creature1: " + strconv.Itoa(int(sc.Creature1)))

	resultBoard, err := simulationgame.RunSimulation(sc)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err.Error())
	}

	jsonBytes, err := json.MarshalIndent(resultBoard, "", " ")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err.Error())
	}

	// fmt.Println("Respond with json data: \n" + string(jsonBytes))

	w.Header().Set("Content-type", "application/json")
	w.Write(jsonBytes)
}

func initServer() {
	rowsRule := Rule{
		StandardValue: int(40),
		MinVal:        5,
		MaxVal:        200,
		ErrorMsg:      "Invalid rows value, value should be between 5-200.",
	}

	colsRule := Rule{
		StandardValue: int(100),
		MinVal:        5,
		MaxVal:        200,
		ErrorMsg:      "Invalid cols value, value should be between 5-200.",
	}

	drawRule := Rule{
		StandardValue: false,
		ErrorMsg:      "Invalid value for draw parameters, must be either true or false.",
	}

	foodsRule := Rule{
		StandardValue: int(75),
		MinVal:        1,
		MaxVal:        150,
		ErrorMsg:      "Invalid value for foods parameter, value should be between 1-150.",
	}

	creature1Rule := Rule{
		StandardValue: uint(10),
		MinVal:        0,
		MaxVal:        50,
		ErrorMsg:      "Invalid value for creature1, should be between 0-50",
	}

	creature2Rule := Rule{
		StandardValue: uint(10),
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

	fmt.Println("Server initialized.")
}

func getSimulationConfig(r *http.Request) (*simulationgame.SimulationConfig, error) {
	sc := simulationgame.SimulationConfig{}

	for key, rule := range parameterRules {
		parameter := r.URL.Query().Get(key)

		fmt.Println("Parameter fetched: " + key + " " + parameter + "its length: " + strconv.Itoa(len(parameter)))

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
					return nil, errors.New(rule.ErrorMsg)

				}
			} else {
				_, ifInt := rule.StandardValue.(int)
				_, ifUint := rule.StandardValue.(uint)

				if ifInt || ifUint {
					value, err := strconv.Atoi(parameter)

					if err != nil {
						return nil, errors.New("For value: " + key + " the value was not an int, as expected.")
					}

					if value < rule.MinVal || value > rule.MaxVal {
						return nil, errors.New(rule.ErrorMsg)
					}

					if ifUint {
						rule.FinalValue = uint(value)

					} else {

						rule.FinalValue = value
					}

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

		creature2, ok := parameterRules["creature2"].FinalValue.(uint)

		if ok {
			sc.Creature2 = creature2
		}
	}

	return &sc, nil
}

func getBoardFromDb(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, "No id given, please check parameter id, currently given id: "+id, http.StatusInternalServerError)
		return

	} else {
		fmt.Println("This is the ID given, looking for this in the db: " + id)
	}

	db, err := openDbConnection()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer db.Close()

	query := "SELECT * FROM simulation_game.boards WHERE id = $1"
	rows, err := db.Query(query, id)

	if err != nil {
		fmt.Println(err.Error())
	}

	var results []DBboard

	for rows.Next() {
		dbboard := DBboard{}

		if err := rows.Scan(&dbboard.Id, &dbboard.Rows, &dbboard.Cols); err != nil {
			http.Error(w, "Database scan error", http.StatusInternalServerError)
			return
		}

		results = append(results, dbboard)
	}

	jsonResponse, err := json.Marshal(results)

	if err != nil {
		http.Error(w, "Error marshaling json", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
	fmt.Println("Responded with json file: " + string(jsonResponse))
}

func testSimulation() {
	testSC := simulationgame.SimulationConfig{
		Rows:      40,
		Cols:      100,
		Foods:     75,
		Draw:      true,
		Creature1: 20,
		Creature2: 10,
	}

	simulationgame.RunSimulation(&testSC)
}
