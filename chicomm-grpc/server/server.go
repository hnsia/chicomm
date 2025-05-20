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
	pr, err := s.storer.CreateProduct(ctx, toStorerProduct(req))
	if err != nil {
		return nil, err
	}

	return toPBProductRes(pr), nil
}

func (s *Server) GetProduct(ctx context.Context, p *pb.ProductReq) (*pb.ProductRes, error) {
	pr, err := s.storer.GetProduct(ctx, p.GetId())
	if err != nil {
		return nil, err
	}

	return toPBProductRes(pr), nil
}

func (s *Server) ListProducts(ctx context.Context, p *pb.ProductReq) (*pb.ListProductRes, error) {
	products, err := s.storer.ListProducts(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]*pb.ProductRes, 0, len(products))
	for _, p := range products {
		res = append(res, toPBProductRes(p))
	}

	return &pb.ListProductRes{Products: res}, nil
}

func (s *Server) UpdateProduct(ctx context.Context, p *pb.ProductReq) (*pb.ProductRes, error) {
	product, err := s.storer.GetProduct(ctx, p.GetId())
	if err != nil {
		return nil, err
	}

	patchProductReq(product, p)

	pr, err := s.storer.UpdateProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	return toPBProductRes(pr), nil
}

func (s *Server) DeleteProduct(ctx context.Context, p *pb.ProductReq) (*pb.ProductRes, error) {
	err := s.storer.DeleteProduct(ctx, p.GetId())
	if err != nil {
		return nil, err
	}

	return &pb.ProductRes{}, nil
}
