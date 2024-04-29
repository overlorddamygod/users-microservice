package main

import (
	"log"
	"net"

	"github.com/overlorddamygod/users-microservice/models"
	"github.com/overlorddamygod/users-microservice/users_grpc_server/server"
	pb "github.com/overlorddamygod/users-microservice/users_proto"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	dsn  = "root:very_secure_password@tcp(localhost:3306)/users?parseTime=true"
	port = ":3001"
)

func main() {
	db := InitDB()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, server.NewServer(db))
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.User{})
	return db
}
