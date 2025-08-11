package service

import (
	"shopping-service/internal/domain"
	"shopping-service/internal/ports"
)

type ProductService struct {
	repo ports.ProductRepository
}

// Removed incomplete NewCartService method causing missing return error.

// NewProductService creates a new ProductService instance.
func NewProductService(r ports.ProductRepository) *ProductService {
	return &ProductService{repo: r}
}

func (s *ProductService) CreateProduct(p *domain.Product) error {
	return s.repo.Create(p)
}

func (s *ProductService) ListProducts() ([]domain.Product, error) {
	return s.repo.FindAll()
}
func (s *ProductService) GetAllProducts() ([]domain.Product, error) {
	products, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	return products, nil
}
