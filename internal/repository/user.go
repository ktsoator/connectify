package repository

import (
	"context"
	"errors"

	"github.com/ktsoator/connectify/internal/domain"
	"github.com/ktsoator/connectify/internal/repository/dao"
)

var (
	ErrDuplicateEmail = dao.ErrDuplicateEmail
	ErrUserNotFound   = dao.ErrRecordNotFound
)

type UserRepository struct {
	userDAO *dao.UserDAO
}

func NewUserRepository(userDAO *dao.UserDAO) *UserRepository {
	return &UserRepository{userDAO: userDAO}
}

func (r *UserRepository) Create(ctx context.Context, user domain.User) error {
	err := r.userDAO.Insert(ctx, dao.UserModel{
		Email:    user.Email,
		Password: user.Password,
	})
	if err != nil {
		if errors.Is(err, ErrDuplicateEmail) {
			return ErrDuplicateEmail
		}
		return err
	}
	return nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.userDAO.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, dao.ErrRecordNotFound) {
			return domain.User{}, ErrUserNotFound
		}
		return domain.User{}, err
	}
	return domain.User{
		ID:       u.ID,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}
