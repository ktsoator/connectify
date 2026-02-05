package service

import (
	"context"

	"github.com/ktsoator/connectify/internal/domain"
	"github.com/ktsoator/connectify/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) Signup(ctx context.Context, user domain.User) error {

	return s.repo.Create(ctx, user)
}
