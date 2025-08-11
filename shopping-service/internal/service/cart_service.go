package service

import "shopping-service/internal/domain"

type CartRepo interface {
	AddItem(userID, productID string, qty int) error
	GetCart(userID string) (domain.Cart, error)
	ClearCart(userID string) error
}

type CartService struct {
	repo CartRepo
}

func NewCartService(repo CartRepo) *CartService {
	return &CartService{repo: repo}
}

func (s *CartService) AddToCart(userID, productID string, qty int) error {
	return s.repo.AddItem(userID, productID, qty)
}
func (s *CartService) GetCart(userID string) (domain.Cart, error) {
	return s.repo.GetCart(userID)
}
func (s *CartService) ClearCart(userID string) error {
	return s.repo.ClearCart(userID)
}
