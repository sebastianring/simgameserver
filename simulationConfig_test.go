package main

import (
	"fmt"
	// sg "github.com/sebastianring/simulationgame"
	"testing"
)

func TestGettingRandomConfig(t *testing.T) {
	fmt.Println("Testing to get a random simulation config.")
	initRules()

	intervalMap := getStandardIntervalMap()

	for k, v := range intervalMap {
		fmt.Println(k, v)
	}

	sc, err := getRandomSimulationConfigFromInterval(intervalMap)

	if err != nil {
		t.Error("Error creating a new sc: ", err.Error())
	}

	fmt.Println(sc)

}
