package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/overlorddamygod/users-microservice/models"
	pb "github.com/overlorddamygod/users-microservice/users_proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	router     *gin.Engine
	grpcClient pb.UserServiceClient
}

func NewServer(grpcClient pb.UserServiceClient) *Server {
	return &Server{
		router:     gin.Default(),
		grpcClient: grpcClient,
	}
}

func (s *Server) Start(port string) error {
	return s.router.Run(port)
}

func (s *Server) SetupRoutes() {
	s.router.POST("/users", s.HandleAddUser)
	s.router.GET("/users/:id", s.HandleGetUserById)
	s.router.GET("/users", s.HandleGetAllUsers)
}

func (s *Server) HandleAddUser(c *gin.Context) {
	user := models.User{}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	r, err := s.grpcClient.AddUser(c, &pb.AddUserRequest{
		Name:  user.Name,
		Email: user.Email,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error adding user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    r.GetId(),
		"name":  r.GetName(),
		"email": r.GetEmail(),
	})
}

func (s *Server) HandleGetUserById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
		})
		return
	}

	r, err := s.grpcClient.GetUserById(c, &pb.GetUserByIdRequest{
		Id: int32(id),
	})

	if err != nil {
		if statusErr, ok := status.FromError(err); ok {
			switch statusErr.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"message": statusErr.Message(),
				})
				return
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": statusErr.Message(),
				})
			}
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":    r.GetId(),
			"name":  r.GetName(),
			"email": r.GetEmail(),
		},
	})
}

func (s *Server) HandleGetAllUsers(c *gin.Context) {
	r, err := s.grpcClient.GetAllUsers(c, &pb.GetAllUsersRequest{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error getting users",
		})
		return
	}

	if r.Users == nil {
		c.JSON(http.StatusOK, gin.H{
			"users": []string{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": r.GetUsers(),
	})
}
