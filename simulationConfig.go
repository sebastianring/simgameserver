package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	sg "github.com/sebastianring/simulationgame"
	"math/rand"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
)

var parameterRules map[string]*Rule

type Rule struct {
	StandardValue any
	MinVal        any
	MaxVal        any
	ErrorMsg      string
}

type valueInterval interface {
	getMin() int
	getMax() int
}

// Potential generics?
type intInterval struct {
	min int
	max int
}

// Potential generics?
type uintInterval struct {
	min uint
	max uint
}

func (ii *intInterval) getMin() int {
	return ii.min
}

func (ii *intInterval) getMax() int {
	return ii.max
}

func (ui *uintInterval) getMin() int {
	return int(ui.min)
}

func (ui *uintInterval) getMax() int {
	return int(ui.max)
}

func initRules() {
	rowsRule := Rule{
		StandardValue: int(40),
		MinVal:        int(5),
		MaxVal:        int(200),
		ErrorMsg:      "Invalid rows value, value should be between 5-200.",
	}

	colsRule := Rule{
		StandardValue: int(100),
		MinVal:        int(5),
		MaxVal:        int(200),
		ErrorMsg:      "Invalid cols value, value should be between 5-200.",
	}

	drawRule := Rule{
		StandardValue: false,
		ErrorMsg:      "Invalid value for draw parameters, must be either true or false.",
	}

	foodsRule := Rule{
		StandardValue: int(75),
		MinVal:        int(1),
		MaxVal:        int(150),
		ErrorMsg:      "Invalid value for foods parameter, value should be between 1-150.",
	}

	creature1Rule := Rule{
		StandardValue: uint(10),
		MinVal:        uint(0),
		MaxVal:        uint(50),
		ErrorMsg:      "Invalid value for creature1, should be between 0-50",
	}

	creature2Rule := Rule{
		StandardValue: uint(10),
		MinVal:        uint(0),
		MaxVal:        uint(50),
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

func getSimulationConfigFromUrlValues(urlvalues url.Values) (*sg.SimulationConfig, error) {
	finalValue, err := cleanUrlParametersToMap(urlvalues)

	if err != nil {
		fmt.Println("Issue cleaning parameters from url values")

		return nil, err
	}

	sc, err := getValidatedConfigFromMap(finalValue)

	if err != nil {
		fmt.Println("Issue validating configuration from map")

		return nil, err
	}

	return sc, nil
}

func getRandomSimulationConfigFromUrl(r *http.Request) (*sg.SimulationConfig, error) {
	// Does not consider any parameters yet! Please add
	parameters := mux.Vars(r)
	fmt.Println(parameters)

	sc, err := getRandomSimulationConfig()

	if err != nil {
		return nil, err
	}

	return sc, nil
}

func getRandomSimulationConfig() (*sg.SimulationConfig, error) {
	intervalMap := getStandardIntervalMap()
	sc, err := getRandomSimulationConfigFromInterval(intervalMap)

	if err != nil {
		return nil, err
	}

	return sc, nil
}

func getStandardIntervalMap() map[string]valueInterval {
	standardInterval := make(map[string]valueInterval)

	standardInterval["rows"] = &intInterval{min: 50, max: 150}
	standardInterval["cols"] = &intInterval{min: 50, max: 150}
	standardInterval["foods"] = &intInterval{min: 50, max: 200}
	standardInterval["creature1"] = &uintInterval{min: 5, max: 25}
	standardInterval["creature2"] = &uintInterval{min: 5, max: 25}

	return standardInterval
}

func getRandomSimulationConfigFromInterval(intervalMap map[string]valueInterval) (*sg.SimulationConfig, error) {
	valueMap := make(map[string]any)

	for key, rule := range parameterRules {
		if key == "draw" {
			continue
		}

		if _, ok := rule.StandardValue.(uint); ok {
			valueMap[key] = uint(randomValueInInterval(intervalMap[key]))
		} else {
			valueMap[key] = randomValueInInterval(intervalMap[key])
		}
	}

	sc, err := getValidatedConfigFromMap(valueMap)

	if err != nil {
		return nil, errors.New("Validation of configuration failed.")
	}

	return sc, nil
}

func randomValueInInterval(interval valueInterval) int {
	value := rand.Intn(interval.getMax()-interval.getMin()+1) + interval.getMin()

	return value
}

func (r *Rule) validateGenericValue(value any) (any, bool) {
	if value == nil {
		return nil, false
	}

	fmt.Println(reflect.TypeOf(r.StandardValue), reflect.TypeOf(value))

	if reflect.TypeOf(value) == reflect.TypeOf(r.StandardValue) {
		switch v := value.(type) {
		case int:
			min, ok := r.MinVal.(int)

			if !ok {
				return nil, false
			}

			max, ok := r.MaxVal.(int)

			if !ok {
				return nil, false
			}

			if v >= min && v <= max {
				return value, true
			}
		case uint:
			min, ok := r.MinVal.(uint)

			if !ok {
				return nil, false
			}

			max, ok := r.MaxVal.(uint)

			if !ok {
				return nil, false
			}

			if v >= min && v <= max {
				return value, true
			}

		default:
			return nil, false
		}
	}

	return nil, false
}

func (r *Rule) validateValue(value any) (any, bool) {
	if value == nil {
		return nil, false
	}

	if reflect.TypeOf(value) == reflect.TypeOf(r.StandardValue) {
		switch v := value.(type) {
		case bool:
			return value, true
		case int:
			min, ok := r.MinVal.(int)

			if !ok {
				return nil, false
			}

			max, ok := r.MaxVal.(int)

			if !ok {
				return nil, false
			}

			if v >= min && v <= max {
				return value, true
			}
		case uint:
			min, ok := r.MinVal.(uint)

			if !ok {
				return nil, false
			}

			max, ok := r.MaxVal.(uint)

			if !ok {
				return nil, false
			}

			if v >= min && v <= max {
				return value, true
			}

		default:
			return nil, false
		}
	}

	return nil, false
}

func cleanUrlParametersToMap(input url.Values) (map[string]any, error) {
	returnMap := make(map[string]any)

	for key, value := range input {
		switch key {
		case "draw":
			if value[0] == "true" {
				returnMap[key] = true
			} else if value[0] == "false" {
				returnMap[key] = false
			}

		case "rows", "cols", "foods":
			intV, err := strconv.Atoi(value[0])

			if err != nil {
				fmt.Println("Issue with converting url parameter - not a string")
				return nil, errors.New("Issue with converting url parameter - not a string")
			}

			returnMap[key] = intV
		case "creature1", "creature2":
			intV, err := strconv.Atoi(value[0])

			if err != nil {
				fmt.Println("Issue with converting url parameter - not a string")
				return nil, errors.New("Issue with converting url parameter - not a string")
			}

			returnMap[key] = uint(intV)
		}
	}

	return returnMap, nil
}

func getValidatedConfigFromMap(valueMap map[string]any) (*sg.SimulationConfig, error) {
	sc := sg.SimulationConfig{}
	finalValue := make(map[string]any)

	for key, rule := range parameterRules {
		v, ok := rule.validateValue(valueMap[key])

		if !ok {
			// Currently, if there is an issue with a validation,
			// standard value is set instead.
			// Maybe should return an error instead.
			finalValue[key] = rule.StandardValue
		} else {
			finalValue[key] = v
		}
	}

	cols, ok := finalValue["cols"].(int)

	if ok {
		sc.Cols = cols
	} else {
		return nil, errors.New("Error, the value of cols was not an int.")
	}

	rows, ok := finalValue["rows"].(int)

	if ok {
		sc.Rows = rows
	} else {
		return nil, errors.New("Error, the value of rows was not an int.")
	}

	draw, ok := finalValue["draw"].(bool)

	if ok {
		sc.Draw = draw
	} else {
		return nil, errors.New("Error, the value of draw was not a bool.")
	}

	foods, ok := finalValue["foods"].(int)

	if ok {
		sc.Foods = foods
	} else {
		return nil, errors.New("Error, the value of foods was not an int.")
	}

	creature1, ok := finalValue["creature1"].(uint)

	if ok {
		sc.Creature1 = creature1
	} else {
		return nil, errors.New("Error, the value of creature1 was not an uint.")
	}

	creature2, ok := finalValue["creature2"].(uint)

	if ok {
		sc.Creature2 = creature2
	} else {
		return nil, errors.New("Error, the value of creature2 was not an uint.")
	}

	return &sc, nil

}
