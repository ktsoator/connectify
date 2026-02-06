package service

import (
	"context"
	"errors"

	"github.com/ktsoator/connectify/internal/domain"
	"github.com/ktsoator/connectify/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail        = repository.ErrDuplicateEmail
	ErrInvalidUserOrPassword = errors.New("invalid email or password")
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
	// Hash the password using bcrypt before storing it
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)

	err = s.repo.Create(ctx, user)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateEmail) {
			return ErrDuplicateEmail
		}
		return err
	}
	return nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (domain.User, error) {
	// 1. Find user by email
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		// If the user is not found, we return a generic "invalid user or password" error.
		// This is a security best practice to prevent user enumeration attacks.
		if errors.Is(err, repository.ErrUserNotFound) {
			return domain.User{}, ErrInvalidUserOrPassword
		}

		// Return other errors
		return domain.User{}, err
	}

	// 2. Compare the provided password with the stored hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		// If the password does not match, we also return the same generic error
		// to ensure the error message is consistent regardless of whether the email or password was incorrect.
		return domain.User{}, ErrInvalidUserOrPassword
	}

	return user, nil
}
