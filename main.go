package simgameserver

import (
	"fmt"
	"github.com/sebastianring/simulationgame"
	"net/http"
)

func runServer() {
	http.HandleFunc("/api/new_sim", func(w http.ResponseWriter, r *http.Request) {
		resultBoard := runSimulation(false)
		jsonBytes, err := json.Marshal(resultBoard)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println(err.Error())
			os.Exit(1)
		}

		fmt.Println(jsonBytes)

		w.Header().Set("Content-type", "application/json")
		w.Write(jsonBytes)
	})

	fmt.Println("Server running at port 8080")
	http.ListenAndServe(":8080", nil)
}

