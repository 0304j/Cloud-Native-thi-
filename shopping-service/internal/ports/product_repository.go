package ports

import "shopping-service/internal/domain"

type ProductRepository interface {
	Create(product *domain.Product) error
	FindAll() ([]domain.Product, error)
}
type ProductService interface {
	CreateProduct(product *domain.Product) error
	GetAllProducts() ([]domain.Product, error)
}
