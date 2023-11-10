package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type DBboard struct {
	Id   uuid.UUID `json:"id"`
	Rows int       `json:"rows"`
	Cols int       `json:"cols"`
}

func OpenDbConnection() (*sql.DB, error) {

	prefix := "postgres://"
	user := "sim_game"
	password := os.Getenv("SIM_GAME_DB_PW")
	fmt.Println("Password:", password)
	adress, hit := os.LookupEnv("SIM_GAME_DB_IP")

	if !hit {
		fmt.Println("No environment variable set:", adress)
		adress = "5.150.233.156"
	}

	port := "5432"

	database_url := prefix + user + ":" +
		password + "@" + adress + ":" +
		port + "/postgres"

	fmt.Println("Trying to connect to DB: ", database_url)

	db, err := sql.Open("postgres", database_url)

	if err != nil {
		fmt.Println(err.Error())

		return nil, err

	} else {
		fmt.Println("Database connection secured")
	}

	fmt.Println("Pinging database")
	err = db.Ping()

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	fmt.Println("Database pinged!")

	return db, nil
}
