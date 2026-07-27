package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service interface {
	Register(req RegisterRequest) (*User, error)
	Login(req LoginRequest) (*User, error)
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

	// Jika email sudah terdaftar
	if err == nil {
		return nil, errors.New("email already registered")
	}

	// Jika ada error tapi bukan error record not found. Contoh: masalah koneksi database
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
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

func (s *service) Login(req LoginRequest) (*User, error) {
	user, err := s.repo.FindByEmail(req.Email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}

		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.Password),
	)

	if err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}
