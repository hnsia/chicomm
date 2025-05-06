package main

import (
	"log"

	"github.com/hnsia/chicomm/chicomm-api/handler"
	"github.com/hnsia/chicomm/chicomm-api/server"
	"github.com/hnsia/chicomm/chicomm-api/storer"
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
	st := storer.NewMySQLStorer(db.GetDB())
	srv := server.NewServer(st)
	hdl := handler.NewHandler(srv)
	handler.RegisterRoutes(hdl)
	handler.Start(":8080")
}
