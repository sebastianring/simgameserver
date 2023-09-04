package main

import (
	"fmt"
	"net/url"
	"reflect"

	// sg "github.com/sebastianring/simulationgame"
	"testing"
)

func TestGettingRandomConfig(t *testing.T) {
	fmt.Println("Testing to get a random simulation config.")
	initRules()

	intervalMap := getStandardIntervalMap()
	sc, err := getRandomSimulationConfigFromInterval(intervalMap)

	if err != nil {
		t.Error("Error creating a new sc: ", err.Error())
	}

	fmt.Println(sc)
}

func TestParsingUrl(t *testing.T) {
	fmt.Println("Testing url parsing")
	u, err := url.Parse("http://127.0.0.1:8080/api/new_sim?cols=20&rows=30&draw=true&creature2=30")

	if err != nil {
		t.Error("Issue parsing url")
	}

	q := u.Query()

	fmt.Println(reflect.TypeOf(q))

	resultMap, err := cleanUrlParametersToMap(q)

	if err != nil {
		t.Error("Issue cleaning url")
	}

	finalMap, err := getValidatedConfigFromMap(resultMap)

	if err != nil {
		t.Error("Issue validating config.")
	}

	fmt.Println(finalMap)
}
