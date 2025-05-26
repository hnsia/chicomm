package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/hnsia/chicomm/chicomm-grpc/pb"
	"github.com/hnsia/chicomm/token"
	"github.com/hnsia/chicomm/util"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type handler struct {
	ctx        context.Context
	client     pb.ChicommClient
	TokenMaker *token.JWTMaker
}

func NewHandler(client pb.ChicommClient, secretKey string) *handler {
	return &handler{
		ctx:        context.Background(),
		client:     client,
		TokenMaker: token.NewJWTMaker(secretKey),
	}
}

func (h *handler) createProduct(w http.ResponseWriter, r *http.Request) {
	var p ProductReq
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	product, err := h.client.CreateProduct(h.ctx, toPBProductReq(p))
	if err != nil {
		http.Error(w, "error creating product", http.StatusInternalServerError)
		return
	}

	res := toProductRes(product)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (h *handler) getProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "error parsing product id", http.StatusBadRequest)
		return
	}

	product, err := h.client.GetProduct(h.ctx, &pb.ProductReq{Id: i})
	if err != nil {
		http.Error(w, "error getting product", http.StatusInternalServerError)
		return
	}

	res := toProductRes(product)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *handler) listProducts(w http.ResponseWriter, r *http.Request) {
	lpres, err := h.client.ListProducts(h.ctx, &pb.ProductReq{})
	if err != nil {
		http.Error(w, "error listing products", http.StatusInternalServerError)
		return
	}

	var res []ProductRes
	for _, p := range lpres.GetProducts() {
		res = append(res, toProductRes(p))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *handler) updateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "error parsing product id", http.StatusBadRequest)
		return
	}

	var p ProductReq
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}
	p.ID = i

	updatedProduct, err := h.client.UpdateProduct(h.ctx, toPBProductReq(p))
	if err != nil {
		http.Error(w, "error updating product", http.StatusInternalServerError)
		return
	}

	res := toProductRes(updatedProduct)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *handler) deleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "error parsing product id", http.StatusBadRequest)
		return
	}

	if _, err := h.client.DeleteProduct(h.ctx, &pb.ProductReq{Id: i}); err != nil {
		http.Error(w, "error deleting product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Order related handler methods

func (h *handler) createOrder(w http.ResponseWriter, r *http.Request) {
	var o OrderReq
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	claims := r.Context().Value(authKey{}).(*token.UserClaims)
	po := toPBOrderReq(o)
	po.UserId = claims.ID

	createdOrder, err := h.client.CreateOrder(h.ctx, po)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := toOrderRes(createdOrder)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (h *handler) getOrder(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(authKey{}).(*token.UserClaims)

	order, err := h.client.GetOrder(h.ctx, &pb.OrderReq{UserId: claims.ID})
	if err != nil {
		http.Error(w, "error getting order", http.StatusInternalServerError)
		return
	}

	res := toOrderRes(order)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *handler) listOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.client.ListOrders(h.ctx, &pb.OrderReq{})
	if err != nil {
		http.Error(w, "error listing orders", http.StatusInternalServerError)
		return
	}

	var res []OrderRes
	for _, order := range orders.GetOrders() {
		res = append(res, toOrderRes(order))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *handler) updateOrderStatus(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(authKey{}).(*token.UserClaims)

	var o OrderReq
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	status, err := toPBOrderStatus(OrderStatus(o.Status))
	if err != nil {
		http.Error(w, "invalid status", http.StatusBadRequest)
		return
	}

	res, err := h.client.UpdateOrderStatus(h.ctx, &pb.OrderReq{
		Id:        o.ID,
		UserId:    claims.ID,
		UserEmail: claims.Email,
		Status:    status,
	})
	if err != nil {
		http.Error(w, "error updating order status", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *handler) deleteOrder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "error parsing order id", http.StatusBadRequest)
		return
	}

	if _, err := h.client.DeleteOrder(h.ctx, &pb.OrderReq{Id: i}); err != nil {
		http.Error(w, "error deleting order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) createUser(w http.ResponseWriter, r *http.Request) {
	var u UserReq
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	// hash password
	hashedPassword, err := util.HashPassword(u.Password)
	if err != nil {
		http.Error(w, "error hashing password", http.StatusInternalServerError)
		return
	}
	u.Password = hashedPassword

	createdUser, err := h.client.CreateUser(h.ctx, toPBUserReq(u))
	if err != nil {
		http.Error(w, "error creating user", http.StatusInternalServerError)
		return
	}

	res := toUserRes(createdUser)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (h *handler) listUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.client.ListUsers(h.ctx, &pb.UserReq{})
	if err != nil {
		http.Error(w, "error listing users", http.StatusInternalServerError)
		return
	}

	var res ListUserRes
	for _, u := range users.GetUsers() {
		res.Users = append(res.Users, toUserRes(u))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *handler) updateUser(w http.ResponseWriter, r *http.Request) {
	var u UserReq
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	claims := r.Context().Value(authKey{}).(*token.UserClaims)
	u.Email = claims.Email

	updatedUser, err := h.client.UpdateUser(h.ctx, toPBUserReq(u))
	if err != nil {
		http.Error(w, "error updating user", http.StatusInternalServerError)
		return
	}

	res := toUserRes(updatedUser)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "error parsing user id", http.StatusBadRequest)
		return
	}

	if _, err := h.client.DeleteUser(h.ctx, &pb.UserReq{Id: i}); err != nil {
		http.Error(w, "error deleting user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) loginUser(w http.ResponseWriter, r *http.Request) {
	var u LoginUserReq
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	user, err := h.client.GetUser(h.ctx, &pb.UserReq{Email: u.Email})
	if err != nil {
		http.Error(w, "error getting user", http.StatusInternalServerError)
		return
	}

	if err := util.CheckPassword(u.Password, user.Password); err != nil {
		http.Error(w, "wrong password", http.StatusUnauthorized)
		return
	}

	// create a json web token (JWT) and return it as response
	accessToken, accessClaims, err := h.TokenMaker.CreateToken(user.Id, user.Email, user.GetIsAdmin(), 15*time.Minute)
	if err != nil {
		http.Error(w, "error creating access token", http.StatusInternalServerError)
		return
	}

	refreshToken, refreshClaims, err := h.TokenMaker.CreateToken(user.Id, user.Email, user.GetIsAdmin(), 24*time.Hour)
	if err != nil {
		http.Error(w, "error creating refresh token", http.StatusInternalServerError)
		return
	}

	session, err := h.client.CreateSession(h.ctx, &pb.SessionReq{
		Id:           refreshClaims.RegisteredClaims.ID,
		UserEmail:    user.GetEmail(),
		RefreshToken: refreshToken,
		IsRevoked:    false,
		ExpiresAt:    timestamppb.New(refreshClaims.RegisteredClaims.ExpiresAt.Time),
	})
	if err != nil {
		http.Error(w, "error creating session", http.StatusInternalServerError)
		return
	}

	res := LoginUserRes{
		SessionID:             session.GetId(),
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessClaims.RegisteredClaims.ExpiresAt.Time,
		RefreshTokenExpiresAt: refreshClaims.RegisteredClaims.ExpiresAt.Time,
		User:                  toUserRes(user),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *handler) logoutUser(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(authKey{}).(*token.UserClaims)

	if _, err := h.client.DeleteSession(h.ctx, &pb.SessionReq{Id: claims.RegisteredClaims.ID}); err != nil {
		http.Error(w, "error deleting session", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) renewAccessToken(w http.ResponseWriter, r *http.Request) {
	var req RenewAccessTokenReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	refreshClaims, err := h.TokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		http.Error(w, "error verifying token", http.StatusUnauthorized)
		return
	}

	session, err := h.client.GetSession(h.ctx, &pb.SessionReq{Id: refreshClaims.RegisteredClaims.ID})
	if err != nil {
		http.Error(w, "error getting session", http.StatusInternalServerError)
		return
	}

	if session.IsRevoked {
		http.Error(w, "session is revoked", http.StatusUnauthorized)
		return
	}

	if session.GetUserEmail() != refreshClaims.Email {
		http.Error(w, "invalid session", http.StatusUnauthorized)
		return
	}

	accessToken, accessClaims, err := h.TokenMaker.CreateToken(refreshClaims.ID, refreshClaims.Email, refreshClaims.IsAdmin, 15*time.Minute)
	if err != nil {
		http.Error(w, "error creating access token", http.StatusInternalServerError)
		return
	}

	res := RenewAccessTokenRes{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessClaims.RegisteredClaims.ExpiresAt.Time,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *handler) revokeSession(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(authKey{}).(*token.UserClaims)

	if _, err := h.client.RevokeSession(h.ctx, &pb.SessionReq{Id: claims.RegisteredClaims.ID}); err != nil {
		http.Error(w, "error revoking session", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
