package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(req RegisterRequest) (*User, error)
}

type service struct {
	repo UserRepository
}

func NewService(repo UserRepository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Register(req RegisterRequest) (*User, error) {
	_, err := s.repo.FindByEmail(req.Email)
	if err == nil {
		return nil, errors.New("Email already registered")
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}

	user := &User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hash),
		Role:         string(RoleBuyer),
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}
