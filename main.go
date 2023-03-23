package main

import (
	"log"
)

func main() {
	db, err := NewPostgres()
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Init(); err != nil {
		log.Fatal(err)
	}

	server := NewAPIService(":8080", db)
	server.Run()
}
