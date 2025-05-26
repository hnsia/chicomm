package handler

import (
	"fmt"
	"strings"

	"github.com/hnsia/chicomm/chicomm-grpc/pb"
)

func toPBProductReq(p ProductReq) *pb.ProductReq {
	return &pb.ProductReq{
		Id:           p.ID,
		Name:         p.Name,
		Image:        p.Image,
		Category:     p.Category,
		Description:  p.Description,
		Rating:       p.Rating,
		NumReviews:   p.NumReviews,
		Price:        p.Price,
		CountInStock: p.CountInStock,
	}
}

func toProductRes(p *pb.ProductRes) ProductRes {
	return ProductRes{
		Name:         p.Name,
		Image:        p.Image,
		Category:     p.Category,
		Description:  p.Description,
		Rating:       p.Rating,
		NumReviews:   p.NumReviews,
		Price:        p.Price,
		CountInStock: p.CountInStock,
	}
}

func toPBOrderReq(o OrderReq) *pb.OrderReq {
	return &pb.OrderReq{
		PaymentMethod: o.PaymentMethod,
		TaxPrice:      o.TaxPrice,
		ShippingPrice: o.ShippingPrice,
		TotalPrice:    o.TotalPrice,
		Items:         toPBOrderItems(o.Items),
	}
}

func toPBOrderItems(items []*OrderItem) []*pb.OrderItem {
	var res []*pb.OrderItem
	for _, item := range items {
		res = append(res, &pb.OrderItem{
			Name:      item.Name,
			Quantity:  item.Quantity,
			Image:     item.Image,
			Price:     item.Price,
			ProductId: item.ProductID,
		})
	}
	return res
}

type OrderStatus string

const (
	Pending   OrderStatus = "pending"
	Shipped   OrderStatus = "shipped"
	Delivered OrderStatus = "delivered"
)

func toPBOrderStatus(s OrderStatus) (pb.OrderStatus, error) {
	switch s {
	case Pending:
		return pb.OrderStatus_PENDING, nil
	case Shipped:
		return pb.OrderStatus_SHIPPED, nil
	case Delivered:
		return pb.OrderStatus_DELIVERED, nil
	default:
		return 0, fmt.Errorf("invalid order status: %s", s)
	}
}

func toOrderRes(o *pb.OrderRes) OrderRes {
	return OrderRes{
		ID:            o.Id,
		PaymentMethod: o.PaymentMethod,
		TaxPrice:      o.TaxPrice,
		ShippingPrice: o.ShippingPrice,
		TotalPrice:    o.TotalPrice,
		Items:         toOrderItems(o.Items),
		Status:        strings.ToLower(o.GetStatus().String()),
	}
}

func toOrderItems(items []*pb.OrderItem) []*OrderItem {
	var res []*OrderItem
	for _, item := range items {
		res = append(res, &OrderItem{
			Name:      item.Name,
			Quantity:  item.Quantity,
			Image:     item.Image,
			Price:     item.Price,
			ProductID: item.ProductId,
		})
	}
	return res
}

func toPBUserReq(u UserReq) *pb.UserReq {
	return &pb.UserReq{
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
		IsAdmin:  u.IsAdmin,
	}
}

func toUserRes(u *pb.UserRes) UserRes {
	return UserRes{
		Name:    u.Name,
		Email:   u.Email,
		IsAdmin: u.IsAdmin,
	}
}
