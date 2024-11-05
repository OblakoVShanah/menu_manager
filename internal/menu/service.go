package menu

import (
	"context"
	"menu_manager/internal/oops"
	"time"

	"github.com/google/uuid"
)

// AppService реализует бизнес-логику работы с меню
type AppService struct {
	storage Store
	// Здесь можно добавить клиент для ChatGPT
	// chatGPT ChatGPTClient
}

// NewService создает новый экземпляр сервиса
func NewService(storage Store) Service {
	return &AppService{
		storage: storage,
	}
}

// GetMenu возвращает меню по ID
func (s *AppService) GetMenu(ctx context.Context, id string) (*Menu, error) {
	menu, err := s.storage.LoadMenu(ctx, id)
	if err != nil {
		return nil, err
	}
	return menu, nil
}

// CreateMenu создает новое меню
func (s *AppService) CreateMenu(ctx context.Context, userID string, startDate, endDate time.Time) (*Menu, error) {
	if startDate.After(endDate) {
		return nil, oops.NewValidationError("dates", oops.ErrInvalidDates)
	}

	menu := &Menu{
		ID:        uuid.New().String(),
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
		CreatedAt: time.Now(),
	}

	if err := s.storage.SaveMenu(ctx, menu); err != nil {
		return nil, err
	}

	return menu, nil
}

// GenerateMenu генерирует меню с помощью ChatGPT
func (s *AppService) GenerateMenu(ctx context.Context, userID string, preferences map[string]interface{}) (*Menu, error) {
	// TODO: Реализовать генерацию меню через ChatGPT
	return nil, oops.ErrNotImplemented
}

// AddRecipe добавляет рецепт в меню
func (s *AppService) AddRecipe(ctx context.Context, menuID string, recipe Recipe) error {
	menu, err := s.storage.LoadMenu(ctx, menuID)
	if err != nil {
		return err
	}

	if err := s.storage.SaveRecipe(ctx, &recipe); err != nil {
		return err
	}

	// Добавить рецепт к соответствующему приему пищи
	// TODO: Реализовать логику добавления рецепта к конкретному приему пищи

	return s.storage.SaveMenu(ctx, menu)
}
