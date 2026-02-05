package repository

import (
	"context"

	"github.com/ktsoator/connectify/internal/domain"
	"github.com/ktsoator/connectify/internal/repository/dao"
)

type UserRepository struct {
	userDAO *dao.UserDAO
}

func NewUserRepository(userDAO *dao.UserDAO) *UserRepository {
	return &UserRepository{userDAO: userDAO}
}

func (r *UserRepository) Create(ctx context.Context, user domain.User) error {
	return r.userDAO.Insert(ctx, dao.User{
		Email:    user.Email,
		Password: user.Password,
	})
}
