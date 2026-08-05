package product

import (
	"github.com/google/uuid"
)

type Service interface {
	Create(req CreateProductRequest, sellerID uuid.UUID) (*Product, error)
	List() ([]Product, error)
	GetByID(id uuid.UUID) (*Product, error)
	Update(id uuid.UUID, userID uuid.UUID, req UpdateProductRequest) (*Product, error)
	Delete(id uuid.UUID, userID uuid.UUID) error
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

func (s *service) Update(id uuid.UUID, userID uuid.UUID, req UpdateProductRequest) (*Product, error) {

	product, err := s.repo.FindByID(id)

	if err != nil {
		return nil, err
	}

	if isOwner := product.SellerID == userID; !isOwner {
		return nil, ErrForbiddenAccess
	}

	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.Stock = req.Stock

	err = s.repo.Update(product)

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *service) Delete(id uuid.UUID, userID uuid.UUID) error {

	product, err := s.repo.FindByID(id)

	if err != nil {
		return err
	}

	if product.SellerID == userID {
		return ErrForbiddenAccess
	}

	return s.repo.Delete(id)
}
