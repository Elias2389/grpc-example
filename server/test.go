package server

import (
	"ae.com/proto-buffers/model"
	"ae.com/proto-buffers/respository"
	"ae.com/proto-buffers/testpb"
	"context"
	"io"
	"log"
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
