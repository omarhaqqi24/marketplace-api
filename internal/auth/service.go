package auth

import (
	"errors"

	"github.com/omarhaqqi24/marketplace-api/internal/config"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service interface {
	Register(req RegisterRequest) (*User, error)
	Login(req LoginRequest) (string, error)
}

type service struct {
	repo UserRepository
	cfg  *config.Config
}

func NewService(repo UserRepository, cfg *config.Config) Service {
	return &service{
		repo: repo,
		cfg:  cfg,
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

func (s *service) Login(req LoginRequest) (string, error) {
	user, err := s.repo.FindByEmail(req.Email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrInvalidCredentials
		}

		return "", err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.Password),
	)

	if err != nil {
		return "", ErrInvalidCredentials
	}

	return GenerateToken(user, s.cfg.JWTSecret)
}
