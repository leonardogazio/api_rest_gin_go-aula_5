package main

import (
	"github.com/guilhermeonrails/api-go-gin/database"
	"github.com/guilhermeonrails/api-go-gin/routes"
)

func main() {
	database.NewRepo(false)
	routes.HandleRequests()
}
