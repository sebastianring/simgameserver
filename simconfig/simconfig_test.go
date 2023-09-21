package simconfig

import (
	"fmt"
	"net/url"
	"reflect"
	// sg "github.com/sebastianring/simulationgame"
	"testing"
)

func TestGetRandomSimulationConfig(t *testing.T) {
	fmt.Println("Testing to get a random simulation config.")
	initRules()

	intervalMap := GetStandardIntervalMap()
	sc, err := GetRandomSimulationConfigFromInterval(intervalMap)

	if err != nil {
		t.Error("Error creating a new sc: ", err.Error())
	}

	fmt.Println(sc)
}

func TestGetMultipleRandomSimulationConfig(t *testing.T) {
	fmt.Println("Testing to get a random simulation config.")
	initRules()

	iterations := 25
	counter := 0

	for i := 0; i < iterations; i++ {
		intervalMap := GetStandardIntervalMap()
		_, err := GetRandomSimulationConfigFromInterval(intervalMap)

		if err != nil {
			t.Error("Error creating a new sc: ", err.Error())
		} else {
			counter++
		}
	}

	if counter != iterations {
		t.Error("Config failed this many times: ", iterations-1-counter)
	}
}

func TestParsingUrl(t *testing.T) {
	fmt.Println("Testing url parsing")
	u, err := url.Parse("http://127.0.0.1:8080/api/new_sim?cols=20&rows=30&draw=true&creature2=30")

	if err != nil {
		t.Error("Issue parsing url")
	}

	q := u.Query()

	fmt.Println(reflect.TypeOf(q))

	resultMap, err := CleanUrlParametersToMap(q)

	if err != nil {
		t.Error("Issue cleaning url")
	}

	finalMap, err := GetValidatedConfigFromMap(resultMap)

	if err != nil {
		t.Error("Issue validating config.")
	}

	fmt.Println(finalMap)
}
