package main

import (
	"log"
	"net"

	"github.com/hnsia/chicomm/chicomm-grpc/pb"
	"github.com/hnsia/chicomm/chicomm-grpc/server"
	"github.com/hnsia/chicomm/chicomm-grpc/storer"
	"github.com/hnsia/chicomm/db"
	"github.com/ianschenck/envflag"
	"google.golang.org/grpc"
)

func main() {
	var (
		svcAddr = envflag.String("SVC_ADDR", "0.0.0.0:9091", "address where the chicomm-grpc service is listening on")
	)

	// instantiate db
	db, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("error opening db: %v", err)
	}
	defer db.Close()
	log.Println("successfully connected to db")

	// instantiate server
	st := storer.NewMySQLStorer(db.GetDB())
	srv := server.NewServer(st)

	// register our server with the gRPC server
	grpcSrv := grpc.NewServer()
	pb.RegisterChicommServer(grpcSrv, srv)

	listener, err := net.Listen("tcp", *svcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("server listening at %s", *svcAddr)
	if err := grpcSrv.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
