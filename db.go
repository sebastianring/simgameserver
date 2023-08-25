package main

import (
	"database/sql"
	// _ "database/sql/driver"
	"fmt"
	_ "github.com/lib/pq"
	"os"
	// _ "github.com/lib/pq"
)

func openDbConnection() (*sql.DB, error) {
	prefix := "postgres://"
	user := "sim_game"
	password := os.Getenv("SIM_GAME_DB_PW")
	adress := "5.150.233.156"
	// adress := "192.168.0.130"
	// adress := "localhost"
	port := "5432"

	database_url := prefix + user + ":" +
		password + "@" + adress + ":" +
		port + "/postgres"

	db, err := sql.Open("postgres", database_url)

	if err != nil {
		fmt.Println(err.Error())

		return nil, err
	} else {
		fmt.Println("Database connection secured!")
	}

	// err = db.Ping()
	//
	// if err != nil {
	// 	fmt.Println(err.Error())
	//
	// 	return nil, err
	//
	// } else {
	// 	fmt.Println("Database pinged!")
	// }

	return db, nil
}
