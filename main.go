package main

import (
	"fmt"
	"log"
	"net"

	"github.com/DuvanAlbarracin/movies_auth/pkg/config"
	"github.com/DuvanAlbarracin/movies_auth/pkg/db"
	"github.com/DuvanAlbarracin/movies_auth/pkg/proto"
	"github.com/DuvanAlbarracin/movies_auth/pkg/services"
	"github.com/DuvanAlbarracin/movies_auth/pkg/utils"
	"google.golang.org/grpc"
)

func main() {
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Failed loading config:", err)
	}

	h := db.Init(c.DBUrl)
	jwt := utils.JwtWrapper{
		SecretKey:       c.JWTSecretKey,
		Issuer:          "movies_auth",
		ExpirationHours: 1,
	}

	lis, err := net.Listen("tcp", c.Port)
	if err != nil {
		log.Fatalln("Failed to listening:", err)
	}

	fmt.Println("Auth service on:", c.Port)

	s := services.Server{
		H:   h,
		Jwt: jwt,
	}
	defer s.H.Conn.Close()

	grpcServer := grpc.NewServer()
	proto.RegisterAuthServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}
