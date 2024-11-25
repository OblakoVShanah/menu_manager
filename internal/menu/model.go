package menu

import (
	"context"
	common "menu_manager/internal/models"
	"time"
)

// Menu представляет план питания на определенный период
type Menu struct {
	MealID   string
	Time     time.Time // когда надо кушать
	MealType string    // завтрак, обед, ужин
}

// Meal представляет прием пищи
type Meal struct {
	MealID         string                          `json:"id"`
	DishIDs        []string                        `json:"ID_dish"`
	DishNames      []string                        `json:"dishname"`
	Type           MealType                        `json:"type"`   // завтрак, обед, ужин
	Recipes        []string                        `json:"recipe"` // json с рецептом и списком продуктов
	TotalNutrition common.NutritionalValueAbsolute `json:"total_nutrition"`
}

// MealType определяет тип приема пищи
type MealType string

const (
	MealTypeBreakfast MealType = "breakfast"
	MealTypeLunch     MealType = "lunch"
	MealTypeDinner    MealType = "dinner"
	MealTypeSnack     MealType = "snack"
)

// Service определяет интерфейс для работы с меню
type Service interface {
	// GetMeal возвращает прием пищи и его рецепт со списком продуктов, которые нужно докупить
	GetMeal(ctx context.Context, userID string) (*Meal, string, error)
	// rescheduleMenu обновляет время и даты приемов пищи
	RescheduleMenu(ctx context.Context, currentMenu []Menu, userID string) ([]Menu, error)
	// GetMenu возвращает меню по ID
	GetMenu(ctx context.Context, userID string) ([]Menu, error)
}

// Store определяет интерфейс для хранения меню
type Store interface {
	// LoadMenu возвращает меню из БД со списком id приемов пиши и их запланированного времени
	LoadMenu(ctx context.Context, userID string) ([]Menu, error)
	// LoadMeal возвращает из базы прием пищи с описанием составляющих его блюд и продуктов
	LoadMeal(ctx context.Context, MealID string) (*Meal, error)
	// UpdateMenu обновляет время и даты приемов пищи
	UpdateMenu(ctx context.Context, userID string, menuList []Menu) error
}

type Client interface {
	// GetProducts получает список продуктов для покупки у сервиса barn manager
	GetProducts(ctx context.Context, recipes []string) (string, error)
}
