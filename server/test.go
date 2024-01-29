package server

import (
	"ae.com/proto-buffers/model"
	"ae.com/proto-buffers/respository"
	"ae.com/proto-buffers/studentpb"
	"ae.com/proto-buffers/testpb"
	"context"
	"io"
	"log"
	"time"
)

type TestServer struct {
	repository respository.Repository
	testpb.UnimplementedTestServiceServer
}

func NewTestServer(repository respository.Repository) *TestServer {
	return &TestServer{repository: repository}
}

func (s *TestServer) GetTest(ctx context.Context, req *testpb.GetTestRequest) (*testpb.Test, error) {
	test, err := s.repository.GetTest(ctx, req.Id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &testpb.Test{
		Id:   test.Id,
		Name: test.Name,
	}, nil
}

func (s *TestServer) SetTest(ctx context.Context, req *testpb.Test) (*testpb.SetTestResponse, error) {
	test := &model.Test{
		Id:   req.GetId(),
		Name: req.GetName(),
	}

	err := s.repository.SetTest(ctx, test)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &testpb.SetTestResponse{Id: test.Id}, nil

}

func (s *TestServer) SetQuestions(stream testpb.TestService_SetQuestionsServer) error {
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(
				&testpb.SetQuestionResponse{Ok: true},
			)
		}
		if err != nil {
			return err
		}

		question := &model.Question{
			Id:       msg.GetId(),
			Question: msg.GetQuestion(),
			Answer:   msg.GetAnswer(),
			TestId:   msg.GetTestId(),
		}

		err = s.repository.
			SetQuestion(context.Background(), question)
		if err != nil {
			return stream.SendAndClose(
				&testpb.SetQuestionResponse{Ok: false},
			)
		}

	}
}

func (s *TestServer) EnrollStudents(stream testpb.TestService_EnrollStudentsServer) error {
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&testpb.SetQuestionResponse{
				Ok: true,
			})
		}
		if err != nil {
			log.Fatalf("Error reading stream: %v", err)
			return err
		}
		enrollment := &model.Enrollment{
			StudentId: msg.GetStudentId(),
			TestId:    msg.GetTestId(),
		}
		err = s.repository.SetEnrollment(context.Background(), enrollment)
		if err != nil {
			return stream.SendAndClose(&testpb.SetQuestionResponse{
				Ok: false,
			})
		}

	}
}

func (s *TestServer) GetStudentsPerTest(req *testpb.GetStudentsPerTestRequest, stream testpb.TestService_GetStudentsPerTestServer) error {
	students, err := s.repository.GetStudentsPerTest(context.Background(), req.GetTestId())
	if err != nil {
		return err
	}
	for _, student := range students {
		student := &studentpb.Student{
			Id:   student.Id,
			Name: student.Name,
			Age:  student.Age,
		}
		err := stream.Send(student) // Send data to the client
		time.Sleep(2 * time.Second)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *TestServer) TakeTest(stream testpb.TestService_TakeTestServer) error {
	questions, err := s.repository.GetQuestionsPerTest(context.Background(), "t1")
	if err != nil {
		return err
	}
	i := 0
	var currentQuestion = &model.Question{}
	for {
		if i < len(questions) {
			currentQuestion = questions[i]
		}

		if i <= len(questions)-1 {
			questionToSend := &testpb.Question{
				Id:       currentQuestion.Id,
				Question: currentQuestion.Question,
			}
			err := stream.Send(questionToSend)
			if err != nil {
				log.Printf("Error sending question: %v", err)
				return err
			}
			i++
		}
		answer, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Printf("Error receiving answer: %v", err)
			return err
		}
		log.Println("Answer for question:", currentQuestion.Question, "is", answer.GetAnswer())
	}
}
