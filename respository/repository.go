package respository

import (
	"ae.com/proto-buffers/model"
	"context"
)

type Repository interface {
	GetStudent(ctx context.Context, id string) (*model.Student, error)
	SetStudent(ctx context.Context, student *model.Student) error
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
