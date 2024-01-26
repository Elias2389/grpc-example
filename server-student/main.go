package main

import (
	"ae.com/proto-buffers/database"
	"ae.com/proto-buffers/server"
	"ae.com/proto-buffers/studentpb"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	list, err := net.Listen("tcp", ":5060")

	if err != nil {
		log.Fatal(err)
	}

	repo, err := database.NewPostgresRepository()

	server := server.NewStudentServer(repo)

	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	studentpb.RegisterStudentServiceServer(s, server)

	reflection.Register(s)

	fmt.Println("Connected to port :5060")

	if err := s.Serve(list); err != nil {
		log.Fatalf("Error serving: %s", err.Error())
	}

}
