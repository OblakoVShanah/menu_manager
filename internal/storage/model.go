package storage

import (
	"context"
	common "menu_manager/internal/models"
	"time"
)

// FoodProduct представляет собой продукт питания в хранилище
type FoodProduct struct {
	ID                       string                          `json:"id"`
	Name                     string                          `json:"name"`
	WeightPerPkg             uint                            `json:"weight_per_pkg"`
	Amount                   uint                            `json:"amount"`
	PricePerPkg              float32                         `json:"price_per_pkg"`
	ExpirationDate           time.Time                       `json:"expiration_date"`
	PresentInFridge          bool                            `json:"present_in_fridge"`
	NutritionalValueRelative common.NutritionalValueRelative `json:"nutritional_value_relative"`
}

// Service определяет интерфейс для работы с продуктами
type Service interface {
	// AvailableProducts возвращает список доступных продуктов
	AvailableProducts(ctx context.Context) ([]FoodProduct, error)
	// PlaceProduct добавляет новый продукт в хранилище
	PlaceProduct(ctx context.Context, product FoodProduct) (id string, err error)
}

// Store определяет интерфейс для хранения продуктов
type Store interface {
	// LoadProducts загружает все продукты из хранилища
	LoadProducts(ctx context.Context) ([]FoodProduct, error)
	// SaveProduct сохраняет продукт в хранилище
	SaveProduct(ctx context.Context, product FoodProduct) (id string, err error)
}
