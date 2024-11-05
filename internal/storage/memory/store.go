package memory

import (
	"context"
	"menu_manager/internal/storage"
	"sync"
)

// Storage реализует хранилище продуктов в памяти
type Storage struct {
	products map[string]storage.FoodProduct
	mu       sync.RWMutex
}

// NewStorage создает новое хранилище в памяти
func NewStorage() *Storage {
	return &Storage{
		products: make(map[string]storage.FoodProduct),
	}
}

// LoadProducts загружает все продукты из хранилища в памяти
func (s *Storage) LoadProducts(ctx context.Context) ([]storage.FoodProduct, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	products := make([]storage.FoodProduct, 0, len(s.products))
	for _, product := range s.products {
		products = append(products, product)
	}

	return products, nil
}

// SaveProduct сохраняет продукт в хранилище в памяти
func (s *Storage) SaveProduct(ctx context.Context, product storage.FoodProduct) (id string, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.products[product.ID] = product
	return product.ID, nil
}
