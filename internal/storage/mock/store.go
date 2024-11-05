package mock

import (
	"context"
	"menu_manager/internal/storage"
)

type Store struct {
	products []storage.FoodProduct
	err      error
}

func NewStore() *Store {
	return &Store{
		products: make([]storage.FoodProduct, 0),
	}
}

// SetError allows setting an error for testing error scenarios
func (s *Store) SetError(err error) {
	s.err = err
}

// SetProducts allows setting predefined products for testing
func (s *Store) SetProducts(products []storage.FoodProduct) {
	s.products = products
}

func (s *Store) LoadProducts(ctx context.Context) ([]storage.FoodProduct, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.products, nil
}

func (s *Store) SaveProduct(ctx context.Context, product storage.FoodProduct) (id string, err error) {
	if s.err != nil {
		return "", s.err
	}
	s.products = append(s.products, product)
	return product.ID, nil
}
