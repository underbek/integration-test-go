package use_case

import (
	"context"

	"github.com/AndreyAndreevich/articles/user_service/domain"
)

type Storage interface {
	CreateUser(ctx context.Context, name string) (domain.User, error)
	GetUser(ctx context.Context, id int) (domain.User, error)
}

type useCase struct {
	storage Storage
}

func New(storage Storage) *useCase {
	return &useCase{
		storage: storage,
	}
}

func (c *useCase) CreateUser(ctx context.Context, name string) (domain.User, error) {
	return c.storage.CreateUser(ctx, name)
}

func (c *useCase) GetUser(ctx context.Context, id int) (domain.User, error) {
	return c.storage.GetUser(ctx, id)
}
