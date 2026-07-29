package product

import (
	"github.com/google/uuid"
)

type Service interface {
	Create(req CreateProductRequest, sellerID uuid.UUID) (*Product, error)
	List() ([]Product, error)
	GetByID(id uuid.UUID) (*Product, error)
}

type service struct {
	repo ProductRepository
}

func NewService(repo ProductRepository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Create(req CreateProductRequest, sellerID uuid.UUID) (*Product, error) {
	product := &Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		SellerID:    sellerID,
		Stock:       req.Stock,
	}

	err := s.repo.Create(product)

	return product, err
}

func (s *service) List() ([]Product, error) {
	return s.repo.FindAll()
}

func (s *service) GetByID(id uuid.UUID) (*Product, error) {
	return s.repo.FindByID(id)
}
