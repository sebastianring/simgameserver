package trash

//
// import (
// 	"encoding/json"
// 	"fmt"
// 	"github.com/gorilla/mux"
// 	sg "github.com/sebastianring/simulationgame"
// 	"html/template"
// 	"log"
// 	"math/rand"
// 	"net/http"
// 	"simgameserver/api"
// 	sc "simgameserver/simconfig"
// 	"strconv"
// 	"time"
// )
//
// func runServer() {
// 	router := mux.NewRouter()
// 	// router.HandleFunc("/api/new_sim", newSimulation).Methods("GET")
// 	router.HandleFunc("/api/get_sim/{id: [0-9a-fA-F-]+}", getBoardFromDb).Methods("GET")
// 	router.HandleFunc("/api/del_sim/{id: [0-9a-fA-F-]+}", delBoardFromDb).Methods("DELETE")
// 	// router.HandleFunc("/api/new_random_sim", newRandomSimulation).Methods("GET")
// 	// router.HandleFunc("/new_sim_form", newSimForm).Methods("GET", "POST")
// 	// router.HandleFunc("/api/new_multiple_sim_conc", newMultipleRandomSimulationsConcurrent).Methods("GET")
// 	router.HandleFunc("/api/new_multiple_sim", newMultipleRandomSimulations).Methods("GET")
//
// 	http.Handle("/", router)
// 	fmt.Println("Server running at port 8081")
// 	http.ListenAndServe(":8081", nil)
// }
//
// func newSimForm(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == "GET" {
// 		t, err := template.ParseFiles("html/new_sim_form.html")
//
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
//
// 		p := struct {
// 			Title string
// 			Text  string
// 		}{
// 			Title: "Simulation configuration",
// 			Text:  "Please add your simulation configuration below",
// 		}
//
// 		t.Execute(w, p)
//
// 	} else if r.Method == "POST" {
// 		sc, err := sc.getSimulationConfigFromUrlValues(r.PostForm)
//
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
//
// 		resultBoard, err := sg.RunSimulation(sc)
//
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
//
// 		roundData, err := api.getRoundData(resultBoard, api.AliveAtEnd)
//
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
//
// 		jsonBytes, err := json.MarshalIndent(roundData, "", " ")
//
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
//
// 		w.Header().Set("Content-type", "application/json")
// 		w.Write(jsonBytes)
// 	}
// }
//
// // func newRandomSimulation(w http.ResponseWriter, r *http.Request) {
// // 	sc, err := getRandomSimulationConfigFromUrl(r)
// //
// // 	if err != nil {
// // 		http.Error(w, err.Error(), http.StatusInternalServerError)
// // 	}
// //
// // 	fmt.Println("Starting simulation with config: ", sc)
// // 	resultBoard, err := sg.RunSimulation(sc)
// //
// // 	if err != nil {
// // 		http.Error(w, err.Error(), http.StatusInternalServerError)
// // 	}
// //
// // 	roundData, err := getRoundData(resultBoard, AliveAtEnd)
// //
// // 	if err != nil {
// // 		http.Error(w, err.Error(), http.StatusInternalServerError)
// // 		panic(err.Error())
// // 	}
// //
// // 	jsonBytes, err := json.MarshalIndent(roundData, "", " ")
// //
// // 	if err != nil {
// // 		http.Error(w, err.Error(), http.StatusInternalServerError)
// // 		return
// // 	}
// //
// // 	// fmt.Println("Respond with json data: \n" + string(jsonBytes))
// //
// // 	w.Header().Set("Content-type", "application/json")
// // 	w.Write(jsonBytes)
// //
// // }
//
// func newSimulation(w http.ResponseWriter, r *http.Request) {
// 	sc, err := sc.getSimulationConfigFromUrlValues(r.URL.Query())
//
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		fmt.Println(err.Error())
// 	}
//
// 	fmt.Println("Starting simulation with config: ")
// 	fmt.Println("sc.creature1: " + strconv.Itoa(int(sc.Creature1)))
//
// 	resultBoard, err := sg.RunSimulation(sc)
//
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
//
// 	jsonBytes, err := json.MarshalIndent(resultBoard, "", " ")
//
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		panic(err.Error())
// 	}
//
// 	w.Header().Set("Content-type", "application/json")
// 	w.Write(jsonBytes)
// }
//
// // Maybe use go init() function instead?
// func InitServer() {
// 	InitRules()
// 	rand.Seed(time.Now().UnixNano())
// 	log.Println("Server initialized.")
// }
//
// func getBoardFromDb(w http.ResponseWriter, r *http.Request) {
// 	id := r.URL.Query().Get("id")
//
// 	if id == "" {
// 		http.Error(w, "No id given, please check parameter id, currently given id: "+id, http.StatusInternalServerError)
// 		return
//
// 	} else {
// 		fmt.Println("This is the ID given, looking for this in the db: " + id)
// 	}
//
// 	db, err := db.openDbConnection()
//
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
//
// 	defer db.Close()
//
// 	query := "SELECT * FROM simulation_game.boards WHERE id = $1"
// 	rows, err := db.Query(query, id)
//
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}
//
// 	var results []db.DBboard
//
// 	for rows.Next() {
// 		dbboard := DBboard{}
//
// 		if err := rows.Scan(&dbboard.Id, &dbboard.Rows, &dbboard.Cols); err != nil {
// 			http.Error(w, "Database scan error", http.StatusInternalServerError)
// 			return
// 		}
//
// 		results = append(results, dbboard)
// 	}
//
// 	jsonResponse, err := json.Marshal(results)
//
// 	if err != nil {
// 		http.Error(w, "Error marshaling json", http.StatusInternalServerError)
// 		return
// 	}
//
// 	w.Header().Set("Content-type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(jsonResponse)
// 	fmt.Println("Responded with json file: " + string(jsonResponse))
// }
//
// // Need to be worked on
// func delBoardFromDb(w http.ResponseWriter, r *http.Request) {
// 	parameters := mux.Vars(r)
//
// 	id := parameters["id"]
//
// 	fmt.Println(id)
// }
//
// // func newMultipleRandomSimulationsConcurrent(w http.ResponseWriter, r *http.Request) {
// // 	vars := mux.Vars(r)
// //
// // 	var iterations uint
// //
// // 	if len(vars["iterations"]) == 0 {
// // 		iterations = 10
// // 	} else {
// // 		temp, err := strconv.Atoi(vars["iterations"])
// //
// // 		if err != nil {
// // 			http.Error(w, "Error converting parameter iterations to uint", http.StatusInternalServerError)
// // 			return
// // 		}
// //
// // 		if temp < 1 || temp > 100 {
// // 			http.Error(w, "Either too few or too many iterations, interval should be between 1-100.", http.StatusInternalServerError)
// // 			return
// // 		}
// //
// // 		iterations = uint(temp)
// // 	}
// //
// // 	boardMap := [][]*simpleRoundData{}
// // 	wg := sync.WaitGroup{}
// //
// // 	for i := uint(0); i < iterations; i++ {
// // 		wg.Add(1)
// // 		go runRandomSimAsGoRoutine(&w, &boardMap, &wg, i)
// // 	}
// //
// // 	wg.Wait()
// //
// // 	jsonBytes, err := json.MarshalIndent(boardMap, "", " ")
// //
// // 	if err != nil {
// // 		http.Error(w, err.Error(), http.StatusInternalServerError)
// // 	}
// //
// // 	w.Header().Set("Content-type", "application/json")
// // 	w.WriteHeader(http.StatusOK)
// // 	w.Write(jsonBytes)
// //
// // 	fmt.Println("Finished running simulations")
// // }
//
// func newMultipleRandomSimulations(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
//
// 	var iterations uint
//
// 	if len(vars["iterations"]) == 0 {
// 		iterations = 10
// 	} else {
// 		temp, err := strconv.Atoi(vars["iterations"])
//
// 		if err != nil {
// 			http.Error(w, "Error converting parameter iterations to uint", http.StatusInternalServerError)
// 			return
// 		}
//
// 		if temp < 1 || temp > 100 {
// 			http.Error(w, "Either too few or too many iterations, interval should be between 1-100.", http.StatusInternalServerError)
// 			return
// 		}
//
// 		iterations = uint(temp)
// 	}
//
// 	boardMap := [][]*simpleRoundData{}
//
// 	for i := uint(0); i < iterations; i++ {
// 		sc, err := getRandomSimulationConfig()
//
// 		fmt.Println("Starting random simulation with this config: ", sc)
//
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
//
// 		resultBoard, err := sg.RunSimulation(sc)
//
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
//
// 		roundData, err := getRoundData(resultBoard, AliveAtEnd)
//
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
//
// 		boardMap = append(boardMap, roundData)
// 		i++
// 	}
//
// 	jsonBytes, err := json.MarshalIndent(boardMap, "", " ")
//
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
//
// 	w.Header().Set("Content-type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(jsonBytes)
//
// 	fmt.Println("Finished running simulations")
// }
