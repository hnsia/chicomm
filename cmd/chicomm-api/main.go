package main

import (
	"log"

	"github.com/hnsia/chicomm/chicomm-api/handler"
	"github.com/hnsia/chicomm/chicomm-api/server"
	"github.com/hnsia/chicomm/chicomm-api/storer"
	"github.com/hnsia/chicomm/db"
	"github.com/ianschenck/envflag"
)

const minSecretKeySize = 32

func main() {
	var secretKey = envflag.String("SECRET_KEY", "01234567890123456789012345678901", "secret key for JWT signing")
	if len(*secretKey) < minSecretKeySize {
		log.Fatalf("SECRET_KEY must be at least %d characters long", minSecretKeySize)
	}

	db, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("error opening db: %v", err)
	}
	defer db.Close()
	log.Println("successfully connected to db")

	// do something with the db
	st := storer.NewMySQLStorer(db.GetDB())
	srv := server.NewServer(st)
	hdl := handler.NewHandler(srv, *secretKey)
	handler.RegisterRoutes(hdl)
	handler.Start(":8080")
}
