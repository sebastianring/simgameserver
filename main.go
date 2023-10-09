package main

import "github.com/sebastianring/simgameserver/api"

func main() {
	server := api.NewAPIServer(":8081")
	server.Run()
}
