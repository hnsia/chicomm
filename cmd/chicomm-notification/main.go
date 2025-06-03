package main

import (
	"context"
	"log"

	"github.com/hnsia/chicomm/chicomm-grpc/pb"
	"github.com/hnsia/chicomm/chicomm-notification/server"
	"github.com/ianschenck/envflag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var (
		svcAddr    = envflag.String("GRPC_SVC_ADDR", "0.0.0.0:9091", "address where the chicomm-grpc services is listening on")
		adminEmail = envflag.String("ADMIN_EMAIL", "hns-dev@gmail.com", "admin email")
		adminPass  = envflag.String("ADMIN_PASS", "", "admin password")
	)

	envflag.Parse()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient(*svcAddr, opts...)
	if err != nil {
		log.Fatalf("error connecting to chicomm-grpc service: %v", err)
	}
	defer conn.Close()

	client := pb.NewChicommClient(conn)
	srv := server.NewServer(client, &server.AdminInfo{
		Email:    *adminEmail,
		Password: *adminPass,
	})

	done := make(chan struct{})
	go func() {
		srv.Run(context.Background())
		done <- struct{}{}
	}()
	<-done
}
