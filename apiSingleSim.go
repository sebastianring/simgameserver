package main

import (
	sg "github.com/sebastianring/simulationgame"
	"log"
	"net/http"
	"strconv"
)

func (s *APIServer) newSingleSimulation(w http.ResponseWriter, r *http.Request) error {
	sc, err := getSimulationConfigFromUrlValues(r.URL.Query())

	if err != nil {
		log.Println("Error occured during getting simulation config from URL values: ", err)
		return err
	}

	log.Println("Starting simulation with config: ")
	log.Println("sc.creature1: " + strconv.Itoa(int(sc.Creature1)))

	resultBoard, err := sg.RunSimulation(sc)

	if err != nil {
		log.Println("Error occured during running the simulation: ", err)
		return err
	}

	roundData, err := getRoundData(resultBoard, AliveAtEnd)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, roundData)
}

func (s *APIServer) newRandomSimulation(w http.ResponseWriter, r *http.Request) error {
	sc, err := getRandomSimulationConfigFromUrl(r)

	if err != nil {
	}

	log.Println("Starting simulation with config: ", sc)
	resultBoard, err := sg.RunSimulation(sc)

	if err != nil {
	}

	roundData, err := getRoundData(resultBoard, AliveAtEnd)

	if err != nil {
	}

	return WriteJSON(w, http.StatusOK, roundData)
}
