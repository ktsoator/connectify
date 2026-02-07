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
	ErrUserNotFound          = repository.ErrUserNotFound
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

func (s *UserService) Profile(ctx context.Context, id int64) (domain.User, error) {
	// Attempt to retrieve user data from the repository layer
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		// If the repository reports the user wasn't found, map it to a service-level error
		if errors.Is(err, repository.ErrUserNotFound) {
			return domain.User{}, ErrUserNotFound
		}
		// Return any other unexpected infrastructure or database errors
		return domain.User{}, err
	}
	return user, nil
}

func (s *UserService) Update(ctx context.Context, user domain.User) error {
	// Delegate the update operation to the repository layer
	err := s.repo.Update(ctx, user)
	if err != nil {
		// Specifically handle the case where the user record being updated no longer exists
		if errors.Is(err, repository.ErrUserNotFound) {
			return ErrUserNotFound
		}
		return err
	}
	return nil
}
