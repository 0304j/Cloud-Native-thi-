package domain

type Cart struct {
	UserID string
	Items  []CartItem
}

type CartItem struct {
	ProductID string
	Qty       int
}
