package server

import (
	"ae.com/proto-buffers/model"
	"ae.com/proto-buffers/respository"
	"ae.com/proto-buffers/studentpb"
	"context"
	"log"
)

type Server struct {
	repository respository.Repository
	studentpb.UnimplementedStudentServiceServer
}

func NewStudentServer(repository respository.Repository) *Server {
	return &Server{repository: repository}
}

func (s *Server) GetStudent(ctx context.Context, req *studentpb.GetStudentRequest) (*studentpb.Student, error) {
	student, err := s.repository.GetStudent(ctx, req.Id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &studentpb.Student{
		Id:   student.Name,
		Name: student.Name,
		Age:  student.Age,
	}, nil
}

func (s *Server) SetStudent(ctx context.Context, req *studentpb.Student) (*studentpb.SetStudentResponse, error) {
	student := &model.Student{
		Id:   req.GetId(),
		Name: req.GetName(),
		Age:  req.GetAge(),
	}

	err := s.repository.SetStudent(ctx, student)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &studentpb.SetStudentResponse{Id: student.Id}, nil

}
