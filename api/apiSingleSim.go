package api

import (
	"errors"
	"github.com/gorilla/mux"
	ldb "github.com/sebastianring/simgameserver/db"
	sc "github.com/sebastianring/simgameserver/simconfig"
	sg "github.com/sebastianring/simulationgame"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func (s *APIServer) newSingleSimulation(w http.ResponseWriter, r *http.Request) error {
	sc, err := sc.GetSimulationConfigFromUrlValues(r.URL.Query())

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
	sc, err := sc.GetRandomSimulationConfigFromUrl(r)

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

func (s *APIServer) getSimulationForm(w http.ResponseWriter, r *http.Request) error {
	t, err := template.ParseFiles("html/new_sim_form.html")

	if err != nil {
		return errors.New(err.Error())
	}

	p := struct {
		Title string
		Text  string
	}{
		Title: "Simulation configuration",
		Text:  "Please add your simulation configuration below",
	}

	t.Execute(w, p)

	return nil
}

func (s *APIServer) getBoardFromDb(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]

	log.Println("Trying to get a board from db")

	if id == "" {
		return errors.New("No id given, please check parameter id, currently given id: " + id)
	} else {
		log.Println("Looking for this board in the db: " + id)
	}

	db, err := ldb.OpenDbConnection()

	if err != nil {
		return errors.New("Error connecting to DB: " + err.Error())
	}

	defer db.Close()

	query := "SELECT * FROM simulation_game.boards WHERE id = $1"
	rows, err := db.Query(query, id)

	if err != nil {
		return err
	}

	var results []ldb.DBboard

	for rows.Next() {
		dbboard := ldb.DBboard{}

		if err := rows.Scan(&dbboard.Id, &dbboard.Rows, &dbboard.Cols); err != nil {
			return errors.New("Database scan error: " + err.Error())
		}

		results = append(results, dbboard)
	}

	return WriteJSON(w, http.StatusOK, results)
}

// NOT COMPLETED YET
func (s *APIServer) delBoardFromDb(w http.ResponseWriter, r *http.Request) error {
	parameters := mux.Vars(r)

	id := parameters["id"]

	log.Println(id)

	return nil
}
