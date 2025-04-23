package main

import (
	"log"

	"github.com/hnsia/chicomm/db"
)

func main() {
	db, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("error opening db: %v", err)
	}
	defer db.Close()
	log.Println("successfully connected to db")

	// do something with the db
	// st := storer.NewMySQLStorer(db.GetDB())
}
