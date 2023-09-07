package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sebastianring/simulationgame"
	"html/template"
	"net/http"
	"strconv"
)

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
	router := mux.NewRouter()
	router.HandleFunc("/api/new_sim", newSimulation).Methods("GET")
	router.HandleFunc("/api/get_sim/{id: [0-9a-fA-F-]+}", getBoardFromDb).Methods("GET")
	router.HandleFunc("/api/new_random_sim", newRandomSimulation).Methods("GET")
	router.HandleFunc("/new_sim_form", newSimForm).Methods("GET", "POST")
	// http.HandleFunc("/api/new_sim", newSimulation)
	// http.HandleFunc("/api/sim", getBoardFromDb)

	http.Handle("/", router)
	fmt.Println("Server running at port 8080")
	http.ListenAndServe(":8080", nil)
}

func newSimForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, err := template.ParseFiles("html/new_sim_form.html")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		p := struct {
			Title string
			Text  string
		}{
			Title: "Simulation configuration",
			Text:  "Please add your simulation configuration below",
		}

		t.Execute(w, p)
	} else {
		fmt.Println("YO - non GET method was passed")
		fmt.Println(r.Form)
		fmt.Println("foods: ", r.FormValue("foods"))
		fmt.Println("all: ", r.PostForm)
	}
}

func newRandomSimulation(w http.ResponseWriter, r *http.Request) {
	sc, err := getRandomSimulationConfigFromUrl(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Println("Starting simulation with config: ", sc)
	resultBoard, err := simulationgame.RunSimulation(sc)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

func newSimulation(w http.ResponseWriter, r *http.Request) {
	sc, err := getSimulationConfigFromUrl(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err.Error())
	}

	fmt.Println("Starting simulation with config: ")
	fmt.Println("sc.creature1: " + strconv.Itoa(int(sc.Creature1)))

	resultBoard, err := simulationgame.RunSimulation(sc)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	initRules()
	fmt.Println("Server initialized.")
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
