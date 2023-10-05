package main

import (
	api "SSO/api_server"
	"SSO/storage"
	"log"
)

func main() {
	store, err := storage.NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}
	server := api.NewAPIServer(":3001", store)
	server.Run()
}
