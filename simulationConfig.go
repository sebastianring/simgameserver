package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sebastianring/simulationgame"
	"math/rand"
	"net/http"
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
	// getParameter() string
	getMin() int
	getMax() int
}

type intInterval struct {
	// parameter string
	min int
	max int
}

type uintInterval struct {
	// parameter string
	min uint
	max uint
}

func (ii *intInterval) getMin() int {
	return ii.min
}

func (ii *intInterval) getMax() int {
	return ii.max
}

// func (ii *intInterval) getParameter() string {
// 	return ii.parameter
// }

func (ui *uintInterval) getMin() int {
	return int(ui.min)
}

func (ui *uintInterval) getMax() int {
	return int(ui.max)
}

// func (ui *uintInterval) getParameter() string {
// 	return ui.parameter
// }

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

func getSimulationConfigFromUrl(r *http.Request) (*simulationgame.SimulationConfig, error) {
	sc := simulationgame.SimulationConfig{}
	finalValue := make(map[string]any)

	for key, rule := range parameterRules {
		parameter := r.URL.Query().Get(key)

		fmt.Println("Parameter fetched: " + key + " " + parameter + "its length: " + strconv.Itoa(len(parameter)))

		if len(parameter) == 0 {
			finalValue[parameter] = rule.StandardValue

		} else {
			_, ok := rule.StandardValue.(bool)

			if ok {
				if parameter == "true" {
					finalValue[parameter] = true

				} else if parameter == "false" {
					finalValue[parameter] = false

				} else {
					return nil, errors.New(rule.ErrorMsg)

				}
			} else {
				_, ifInt := rule.StandardValue.(int)
				_, ifUint := rule.StandardValue.(uint)

				if ifInt || ifUint {
					value, err := strconv.Atoi(parameter)

					if err != nil {
						return nil, errors.New("For value: " + key + " the value was not an int.")
					}
					//
					// if value < rule.MinVal || value > rule.MaxVal {
					// 	return nil, errors.New(rule.ErrorMsg)
					// }

					if ifUint {
						finalValue[parameter] = uint(value)

					} else {

						finalValue[parameter] = value
					}

				}
			}
		}

		cols, ok := finalValue["cols"].(int)

		if ok {
			sc.Cols = cols
		}

		rows, ok := finalValue["rows"].(int)

		if ok {
			sc.Rows = rows
		}

		draw, ok := finalValue["draw"].(bool)

		if ok {
			sc.Draw = draw
		}

		foods, ok := finalValue["foods"].(int)

		if ok {
			sc.Foods = foods
		}

		creature1, ok := finalValue["creature1"].(uint)

		if ok {
			sc.Creature1 = creature1
		}

		creature2, ok := finalValue["creature2"].(uint)

		if ok {
			sc.Creature2 = creature2
		}
	}

	return &sc, nil
}

func getRandomSimulationConfigFromUrl(r *http.Request) (*simulationgame.SimulationConfig, error) {
	parameters := mux.Vars(r)
	fmt.Println(parameters)

	intervalMap := getStandardIntervalMap()
	sc, err := getRandomSimulationConfigFromInterval(intervalMap)

	if err != nil {
		//http error
		//
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

func getRandomSimulationConfigFromInterval(intervalMap map[string]valueInterval) (*simulationgame.SimulationConfig, error) {
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

func (r *Rule) validateValue(value any) (any, bool) {
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

func getValidatedConfigFromMap(valueMap map[string]any) (*simulationgame.SimulationConfig, error) {
	sc := simulationgame.SimulationConfig{Draw: false}
	finalValue := make(map[string]any)

	for key, rule := range parameterRules {
		v, ok := rule.validateValue(valueMap[key])

		if !ok {
			finalValue[key] = rule.StandardValue
			// return nil, errors.New("Error validating value.")
		} else {
			finalValue[key] = v
		}
	}

	cols, ok := finalValue["cols"].(int)

	if ok {
		sc.Cols = cols
	}

	rows, ok := finalValue["rows"].(int)

	if ok {
		sc.Rows = rows
	}

	// draw, ok := finalValue["draw"].(bool)
	//
	// if ok {
	// 	sc.Draw = draw
	// }

	foods, ok := finalValue["foods"].(int)

	if ok {
		sc.Foods = foods
	}

	creature1, ok := finalValue["creature1"].(uint)

	if ok {
		sc.Creature1 = creature1
	}

	creature2, ok := finalValue["creature2"].(uint)

	if ok {
		sc.Creature2 = creature2
	}

	return &sc, nil

}
