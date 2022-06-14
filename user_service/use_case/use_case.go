package use_case

import (
	"context"

	"github.com/AndreyAndreevich/articles/user_service/domain"
	"github.com/shopspring/decimal"
)

type Storage interface {
	CreateUser(ctx context.Context, name string) (domain.User, error)
	GetUser(ctx context.Context, id int) (domain.User, error)
	UpdateBalance(ctx context.Context, id int, balance decimal.Decimal) error
}

type Billing interface {
	UserDeposit(ctx context.Context, userID int, amount decimal.Decimal) error
}

type useCase struct {
	storage Storage
	billing Billing
}

func New(storage Storage, billing Billing) *useCase {
	return &useCase{
		storage: storage,
		billing: billing,
	}
}

func (c *useCase) CreateUser(ctx context.Context, name string) (domain.User, error) {
	return c.storage.CreateUser(ctx, name)
}

func (c *useCase) GetUser(ctx context.Context, id int) (domain.User, error) {
	return c.storage.GetUser(ctx, id)
}

func (c *useCase) UpdateBalance(ctx context.Context, id int, amount decimal.Decimal) (domain.User, error) {
	err := c.billing.UserDeposit(ctx, id, amount)
	if err != nil {
		return domain.User{}, err
	}

	user, err := c.storage.GetUser(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	user.Balance = user.Balance.Add(amount)

	err = c.storage.UpdateBalance(ctx, id, user.Balance)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}
