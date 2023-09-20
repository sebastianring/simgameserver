package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	sg "github.com/sebastianring/simulationgame"
)

func (s *APIServer) newMultipleRandomSimulationsConcurrent(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	var iterations uint

	if len(vars["iterations"]) == 0 {
		iterations = 10
	} else {
		temp, err := strconv.Atoi(vars["iterations"])

		if err != nil {
			msg := "Error converting parameter iterations to uint: " + err.Error()
			log.Println(msg)
			return errors.New(msg)
		}

		if temp < 1 || temp > 100 {
			msg := "Either too few or too many iterations, interval should be between 1-100." + err.Error()
			fmt.Println(msg)
			return errors.New(msg)
		}

		iterations = uint(temp)
	}

	boardMap := [][]*simpleRoundData{}
	wg := sync.WaitGroup{}

	for i := uint(0); i < iterations; i++ {
		wg.Add(1)
		go s.runRandomSimAsGoRoutine(&boardMap, &wg, i)
	}

	wg.Wait()

	return WriteJSON(w, http.StatusOK, boardMap)
}

func (s *APIServer) runRandomSimAsGoRoutine(target *[][]*simpleRoundData, wg *sync.WaitGroup, id uint) error {
	defer wg.Done()
	fmt.Printf("Routine nr: %d running.\n", id)
	sc, err := getRandomSimulationConfig()

	fmt.Println("Starting random simulation with this config: ", sc)

	if err != nil {
		log.Println(err)
		return err
	}

	resultBoard, err := sg.RunSimulation(sc)

	if err != nil {
		log.Println(err)
		return err
	}

	roundData, err := getRoundData(resultBoard, AliveAtEnd)

	if err != nil {
		log.Println(err)
		return err
	}

	*target = append(*target, roundData)

	return nil
}
