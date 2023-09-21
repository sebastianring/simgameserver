package main

import "github.com/sebastianring/simgameserver/api"

func main() {

	// go runServer()

	server := api.NewAPIServer(":8080")
	server.Run()
}
