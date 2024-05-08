package main

import (
	"log"

	"github.com/overlorddamygod/users-microservice/users_api_server/server"
	pb "github.com/overlorddamygod/users-microservice/users_proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const addr = "localhost:3001"

func main() {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	s := server.NewServer(client)
	s.SetupRoutes()

	log.Fatal(s.Start(":3000"))
}
