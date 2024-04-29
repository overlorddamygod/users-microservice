package server

import (
	"context"
	"errors"

	"github.com/overlorddamygod/users-microservice/models"
	pb "github.com/overlorddamygod/users-microservice/users_proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type Server struct {
	db *gorm.DB
	pb.UnimplementedUserServiceServer
}

func NewServer(db *gorm.DB) *Server {
	return &Server{
		db: db,
	}
}

func (s *Server) AddUser(ctx context.Context, in *pb.AddUserRequest) (*pb.User, error) {
	user := models.User{
		Name:  in.Name,
		Email: in.Email,
	}

	if result := s.db.Create(&user); result.Error != nil {
		return nil, result.Error
	}

	return &pb.User{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *Server) GetUserById(ctx context.Context, in *pb.GetUserByIdRequest) (*pb.User, error) {
	var user models.User

	if result := s.db.First(&user, in.Id); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, result.Error
	}

	return &pb.User{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *Server) GetAllUsers(ctx context.Context, in *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse, error) {
	var users []models.User

	if result := s.db.Find(&users); result.Error != nil {
		return nil, result.Error
	}

	var usersProto []*pb.User
	for _, user := range users {
		usersProto = append(usersProto, &pb.User{
			Id:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		})
	}

	return &pb.GetAllUsersResponse{
		Users: usersProto,
	}, nil
}
