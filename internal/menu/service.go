package menu

import (
	"context"
	"math/rand"
	"menu_manager/internal/oops"
	"time"
)

// AppService реализует бизнес-логику работы с меню
type AppService struct {
	storage Store
	client  Client
}

// NewService создает новый экземпляр сервиса
func NewService(storage Store, client Client) Service {
	return &AppService{
		storage: storage,
		client:  client,
	}
}

func (s *AppService) GetMeal(ctx context.Context, userID string) (*Meal, string, error) {

	// получаем меню
	menu, err := s.getMenu(ctx, userID)
	if err != nil {
		return nil, "", err // can't get menu
	}

	// проверка на актуальность меню
	if !isActual(menu) {
		// если устарело, то обновляем меню
		menu, err = s.rescheduleMenu(ctx, menu, userID)
		if err != nil {
			return nil, "", err // can't get menu
		}
	}

	// выбираем прием пищи ближайщий по времени
	mealID, err := findClosestMeal(menu)
	if err != nil {
		return nil, "", err
	}

	// получаем прием пищи
	meal, err := s.storage.LoadMeal(ctx, mealID)
	if err != nil {
		return nil, "", err
	}

	// запрос продуктов в barn manager
	products, err := s.getProducts(ctx, meal.Recipes)
	if err != nil {
		return nil, "", err
	}

	return meal, products, nil
}

// isActual проверяет наличие блюд на сегодняшний день в меню
func isActual(menu []Menu) bool {
	// Получаем текущую дату
	now := time.Now()
	for _, v := range menu {
		// Сравниваем год, месяц и день
		if v.Time.Year() == now.Year() && v.Time.Month() == now.Month() && v.Time.Day() == now.Day() {
			return true
		}
	}
	return false
}

// findClosestMeal возвращает id ближайшего приема пищи
func findClosestMeal(menu []Menu) (string, error) {

	// Получаем текущую дату
	now := time.Now()
	minDist := 24
	minID := ""
	for _, v := range menu {
		// Сравниваем с точностью до часа
		if v.Time.Year() == now.Year() && v.Time.Month() == now.Month() && v.Time.Day() == now.Day() && (v.Time.Hour()-now.Hour()) < minDist {
			minDist = v.Time.Hour() - now.Hour()
			minID = v.MealID
		}
	}
	if minID != "" {
		return minID, nil
	}
	return "", oops.ErrInvalidDates
}

// GetMenu возвращает меню по ID
func (s *AppService) getMenu(ctx context.Context, userID string) ([]Menu, error) {

	// получил блюдо
	menu, err := s.storage.LoadMenu(ctx, userID)

	if err != nil {
		return nil, err
	}
	return menu, nil
}

// rescheduleMenu обновляет время и даты приемов пищи
func (s *AppService) rescheduleMenu(ctx context.Context, currentMenu []Menu, userID string) ([]Menu, error) {

	// Перемешиваем время приемов пищи случайным образом
	rand.Shuffle(len(currentMenu), func(i, j int) {
		currentMenu[i].Time, currentMenu[j].Time = currentMenu[j].Time, currentMenu[i].Time
	})

	// Обновляем даты на неделю вперед
	for _, v := range currentMenu {
		v.Time = v.Time.Add(7 * 24 * time.Hour)
	}

	if err := s.storage.UpdateMenu(ctx, userID, currentMenu); err != nil {
		return nil, err
	}
	return currentMenu, nil
}

// GetMenu возвращает меню по ID
func (s *AppService) getProducts(ctx context.Context, recipes []string) (string, error) {

	// отсылаем рецепты в barn manager и получаем список продуктов
	products, err := s.client.GetProducts(ctx, recipes)

	if err != nil {
		return "", err
	}
	return products, nil
}
