package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/DuvanAlbarracin/movies_auth/pkg/db"
	"github.com/DuvanAlbarracin/movies_auth/pkg/proto"
	"github.com/DuvanAlbarracin/movies_auth/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	H   db.Handler
	Jwt utils.JwtWrapper
	proto.UnimplementedAuthServiceServer
}

func (s *Server) Register(ctx context.Context,
	req *proto.RegisterRequest) (*proto.RegisterResponse, error) {

	user, err := db.FindUserByEmail(s.H.Conn, req.Email)
	if err != nil {
		st := status.New(codes.AlreadyExists, "User already exists")
		return nil, st.Err()
	}

	user.Username = req.Username
	user.Email = req.Email
	user.Password = utils.HashPassword(req.Password)

	err = db.CreateUser(s.H.Conn, &user)
	if err != nil {
		st := status.New(codes.Internal, "Cannot create the user")
		return nil, st.Err()
	}

	return &proto.RegisterResponse{
		Status: http.StatusCreated,
	}, nil
}

func (s *Server) Login(ctx context.Context,
	req *proto.LoginRequest) (*proto.LoginResponse, error) {
	user, err := db.FindUserByEmail(s.H.Conn, req.Email)
	if err != nil {
		st := status.New(codes.NotFound,
			"There is no User registered with that email",
		)
		return nil, st.Err()
	}

	match := utils.CheckPasswordHash(req.Password, user.Password)
	if !match {
		st := status.New(codes.InvalidArgument,
			"Password incorrect",
		)
		return nil, st.Err()
	}

	token, err := s.Jwt.GenerateToken(user)
	if err != nil {
		st := status.New(codes.Internal,
			"Error generating token:",
		)
		fmt.Println("ERROR EN TOKEN:", err.Error())
		return nil, st.Err()
	}
	fmt.Println("TOKENN EN AUTH:", token)

	return &proto.LoginResponse{
		Status: http.StatusOK,
		Token:  token,
	}, nil
}

func (s *Server) Validate(ctx context.Context,
	req *proto.ValidateRequest) (*proto.ValidateResponse, error) {
	claims, err := s.Jwt.ValidateToken(req.Token)
	if err != nil {
		st := status.New(codes.InvalidArgument,
			"Invalid token",
		)
		return nil, st.Err()

	}

	user, err := db.FindUserByEmail(s.H.Conn, claims.Email)
	if err.Error() == "no rows in result set" {
		st := status.New(codes.NotFound,
			"There is no User registered with that email",
		)
		return nil, st.Err()
	}

	return &proto.ValidateResponse{
		Status: http.StatusOK,
		Id:     user.Id,
	}, nil
}
