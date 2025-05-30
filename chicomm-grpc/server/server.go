package server

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hnsia/chicomm/chicomm-grpc/pb"
	"github.com/hnsia/chicomm/chicomm-grpc/storer"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func (s *Server) CreateOrder(ctx context.Context, o *pb.OrderReq) (*pb.OrderRes, error) {
	order, err := s.storer.CreateOrder(ctx, toStorerOrder(o))
	if err != nil {
		return nil, err
	}
	order.Status = storer.Pending

	_, err = s.storer.EnqueueNotificationEvent(ctx, &storer.NotificationEvent{
		UserEmail:   o.GetUserEmail(),
		OrderStatus: order.Status,
		OrderID:     order.ID,
		Attempts:    0,
	})
	if err != nil {
		return nil, err
	}

	return toPBOrderRes(order), nil
}

func (s *Server) GetOrder(ctx context.Context, o *pb.OrderReq) (*pb.OrderRes, error) {
	order, err := s.storer.GetOrder(ctx, o.GetUserId())
	if err != nil {
		return nil, err
	}

	return toPBOrderRes(order), nil
}

func (s *Server) ListOrders(ctx context.Context, o *pb.OrderReq) (*pb.ListOrderRes, error) {
	orders, err := s.storer.ListOrders(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]*pb.OrderRes, 0, len(orders))
	for _, o := range orders {
		res = append(res, toPBOrderRes(o))
	}

	return &pb.ListOrderRes{Orders: res}, nil
}

func (s *Server) UpdateOrderStatus(ctx context.Context, o *pb.OrderReq) (*pb.OrderRes, error) {
	// validate the order req
	order, err := s.storer.GetOrderStatusByID(ctx, o.GetId())
	if err != nil {
		return nil, err
	}

	if o.GetUserId() != order.UserID {
		return nil, fmt.Errorf("order %d does not belong to user %d", o.GetId(), o.GetUserId())
	}

	sOrderStatus := storer.OrderStatus(strings.ToLower(o.GetStatus().String()))
	if sOrderStatus == order.Status {
		return nil, fmt.Errorf("order status is already %s", order.Status)
	}

	order.Status = sOrderStatus
	order.UpdatedAt = toTimePtr(time.Now())
	updatedOrder, err := s.storer.UpdateOrderStatus(ctx, order)
	if err != nil {
		return nil, err
	}

	// enqueue notification event
	_, err = s.storer.EnqueueNotificationEvent(ctx, &storer.NotificationEvent{
		UserEmail:   o.GetUserEmail(),
		OrderStatus: order.Status,
		OrderID:     order.ID,
		Attempts:    0,
	})
	if err != nil {
		return nil, err
	}

	return toPBOrderRes(updatedOrder), nil
}

func (s *Server) DeleteOrder(ctx context.Context, o *pb.OrderReq) (*pb.OrderRes, error) {
	err := s.storer.DeleteOrder(ctx, o.GetId())
	if err != nil {
		return nil, err
	}

	return &pb.OrderRes{}, nil
}

func (s *Server) CreateUser(ctx context.Context, u *pb.UserReq) (*pb.UserRes, error) {
	user, err := s.storer.CreateUser(ctx, toStorerUser(u))
	if err != nil {
		return nil, err
	}

	return toPBUserRes(user), nil
}

func (s *Server) GetUser(ctx context.Context, u *pb.UserReq) (*pb.UserRes, error) {
	user, err := s.storer.GetUser(ctx, u.GetEmail())
	if err != nil {
		return nil, err
	}

	return toPBUserRes(user), nil
}

func (s *Server) ListUsers(ctx context.Context, u *pb.UserReq) (*pb.ListUserRes, error) {
	users, err := s.storer.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]*pb.UserRes, 0, len(users))
	for _, u := range users {
		res = append(res, toPBUserRes(u))
	}

	return &pb.ListUserRes{Users: res}, nil
}

func (s *Server) UpdateUser(ctx context.Context, u *pb.UserReq) (*pb.UserRes, error) {
	user, err := s.storer.GetUser(ctx, u.GetEmail())
	if err != nil {
		return nil, err
	}

	patchUserReq(user, u)

	updatedUser, err := s.storer.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return toPBUserRes(updatedUser), nil
}

func (s *Server) DeleteUser(ctx context.Context, u *pb.UserReq) (*pb.UserRes, error) {
	err := s.storer.DeleteUser(ctx, u.GetId())
	if err != nil {
		return nil, err
	}

	return &pb.UserRes{}, nil
}

func (s *Server) CreateSession(ctx context.Context, sr *pb.SessionReq) (*pb.SessionRes, error) {
	session, err := s.storer.CreateSession(ctx, &storer.Session{
		ID:           sr.GetId(),
		UserEmail:    sr.GetUserEmail(),
		RefreshToken: sr.GetRefreshToken(),
		IsRevoked:    sr.GetIsRevoked(),
		ExpiresAt:    sr.GetExpiresAt().AsTime(),
	})
	if err != nil {
		return nil, err
	}

	return &pb.SessionRes{
		Id:           session.ID,
		UserEmail:    session.UserEmail,
		RefreshToken: session.RefreshToken,
		IsRevoked:    session.IsRevoked,
		ExpiresAt:    timestamppb.New(session.ExpiresAt),
	}, nil
}

func (s *Server) GetSession(ctx context.Context, sr *pb.SessionReq) (*pb.SessionRes, error) {
	session, err := s.storer.GetSession(ctx, sr.GetId())
	if err != nil {
		return nil, err
	}

	return &pb.SessionRes{
		Id:           session.ID,
		UserEmail:    session.UserEmail,
		RefreshToken: session.RefreshToken,
		IsRevoked:    session.IsRevoked,
		ExpiresAt:    timestamppb.New(session.ExpiresAt),
	}, nil
}

func (s *Server) RevokeSession(ctx context.Context, sr *pb.SessionReq) (*pb.SessionRes, error) {
	err := s.storer.RevokeSession(ctx, sr.GetId())
	if err != nil {
		return nil, err
	}

	return &pb.SessionRes{}, nil
}

func (s *Server) DeleteSession(ctx context.Context, sr *pb.SessionReq) (*pb.SessionRes, error) {
	err := s.storer.DeleteSession(ctx, sr.GetId())
	if err != nil {
		return nil, err
	}

	return &pb.SessionRes{}, nil
}
