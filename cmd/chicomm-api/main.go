package main

import (
	"log"

	"github.com/hnsia/chicomm/chicomm-api/handler"
	"github.com/hnsia/chicomm/chicomm-grpc/pb"
	"github.com/ianschenck/envflag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const minSecretKeySize = 32

func main() {
	var (
		secretKey = envflag.String("SECRET_KEY", "01234567890123456789012345678901", "secret key for JWT signing")
		svcAddr   = envflag.String("GRPC_SVC_ADDR", "0.0.0.0:9091", "address where the chicomm-grpc service is listening on")
	)
	envflag.Parse()

	if len(*secretKey) < minSecretKeySize {
		log.Fatalf("SECRET_KEY must be at least %d characters long", minSecretKeySize)
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient(*svcAddr, opts...)
	if err != nil {
		log.Fatalf("error connecting to chicomm-grpc service: %v", err)
	}
	defer conn.Close()

	client := pb.NewChicommClient(conn)

	hdl := handler.NewHandler(client, *secretKey)
	handler.RegisterRoutes(hdl)
	handler.Start(":8080")
}
