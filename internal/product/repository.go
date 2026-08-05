package product

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(product *Product) error
	FindAll() ([]Product, error)
	FindByID(id uuid.UUID) (*Product, error)
	Update(product *Product) error
	Delete(id uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(product *Product) error {
	return r.db.Create(product).Error
}

func (r *repository) FindAll() ([]Product, error) {
	var products []Product

	err := r.db.Find(&products).Error

	return products, err
}

func (r *repository) FindByID(id uuid.UUID) (*Product, error) {
	var product *Product

	err := r.db.First(&product, "id = ?", id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrProductNotFound
	}

	return product, err
}

func (r *repository) Update(product *Product) error {
	return r.db.Model(product).Updates(map[string]any{
		"name":        product.Name,
		"description": product.Description,
		"price":       product.Price,
		"stock":       product.Stock,
	}).Error
}

func (r *repository) Delete(id uuid.UUID) error {
	return r.db.Delete(&Product{}, "id = ?", id).Error
}
