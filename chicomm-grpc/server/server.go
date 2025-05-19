package server

import (
	"context"

	"github.com/hnsia/chicomm/chicomm-grpc/pb"
	"github.com/hnsia/chicomm/chicomm-grpc/storer"
)

type Server struct {
	storer *storer.MySQLStorer
	pb.UnimplementedChicommServer
}

func NewServer(storer *storer.MySQLStorer) *Server {
	return &Server{storer: storer}
}

func (s *Server) CreateProduct(ctx context.Context, req *pb.ProductReq) (*pb.ProductRes, error) {
	return &pb.ProductRes{}, nil
}
