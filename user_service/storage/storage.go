package storage

import (
	"context"
	"errors"

	"github.com/AndreyAndreevich/articles/user_service/domain"
	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

var (
	IncorrectQueryResponse = errors.New("incorrect query response")
)

type storage struct {
	db *sqlx.DB
}

func New(dsn string) (*storage, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return &storage{
		db: db,
	}, nil
}

func (s *storage) CreateUser(ctx context.Context, name string) (domain.User, error) {
	query := `INSERT INTO users (name) VALUES ($1) RETURNING id, name, balance, created_at, updated_at`
	res, err := s.db.QueryxContext(ctx, query, name)
	if err != nil {
		return domain.User{}, err
	}

	defer res.Close()

	if !res.Next() {
		return domain.User{}, IncorrectQueryResponse
	}
	var resUser domain.User
	if err := res.StructScan(&resUser); err != nil {
		return domain.User{}, err
	}

	return resUser, nil
}

func (s *storage) GetUser(ctx context.Context, id int) (domain.User, error) {
	query := `SELECT id, name, balance FROM users WHERE id = $1`

	var user domain.User
	err := s.db.GetContext(ctx, &user, query, id)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}
