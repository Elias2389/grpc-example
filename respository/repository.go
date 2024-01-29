package respository

import (
	"ae.com/proto-buffers/model"
	"context"
)

type Repository interface {
	GetStudent(ctx context.Context, id string) (*model.Student, error)
	SetStudent(ctx context.Context, student *model.Student) error
	GetTest(ctx context.Context, id string) (*model.Test, error)
	SetTest(ctx context.Context, test *model.Test) error
	SetQuestion(ctx context.Context, test *model.Question) error
}

var impl Repository

func SetRepository(repository Repository) {
	impl = repository
}

func GetStudent(ctx context.Context, id string) (*model.Student, error) {
	return impl.GetStudent(ctx, id)
}

func SetStudent(ctx context.Context, student *model.Student) error {
	return impl.SetStudent(ctx, student)
}

func GetTest(ctx context.Context, id string) (*model.Test, error) {
	return impl.GetTest(ctx, id)
}

func SetTest(ctx context.Context, test *model.Test) error {
	return impl.SetTest(ctx, test)
}

func SetQuestion(ctx context.Context, test *model.Question) error {
	return impl.SetQuestion(ctx, test)
}
