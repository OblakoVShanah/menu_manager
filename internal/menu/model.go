package menu

import (
	"context"
	common "menu_manager/internal/models"
	"time"
)

// Menu представляет план питания на определенный период
type Menu struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Meals       []Meal    `json:"meals"`
	CreatedAt   time.Time `json:"created_at"`
	GeneratedBy string    `json:"generated_by"` // например, "chatgpt"
}

// Meal представляет прием пищи
type Meal struct {
	ID             string                          `json:"id"`
	MenuID         string                          `json:"menu_id"`
	Type           MealType                        `json:"type"` // завтрак, обед, ужин
	Time           time.Time                       `json:"time"`
	Recipes        []Recipe                        `json:"recipes"`
	TotalNutrition common.NutritionalValueAbsolute `json:"total_nutrition"`
}

// Recipe представляет рецепт блюда
type Recipe struct {
	ID          string                          `json:"id"`
	Name        string                          `json:"name"`
	Description string                          `json:"description"`
	Steps       []string                        `json:"steps"`
	Ingredients []Ingredient                    `json:"ingredients"`
	Nutrition   common.NutritionalValueAbsolute `json:"nutrition"`
	CookingTime int                             `json:"cooking_time"` // в минутах
}

// Ingredient представляет ингредиент в рецепте
type Ingredient struct {
	ProductID string `json:"product_id"`
	Name      string `json:"name"`
	Amount    uint   `json:"amount"`
	Unit      string `json:"unit"` // грамм, штук и т.д.
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
	// GetMenu возвращает меню по ID
	GetMenu(ctx context.Context, id string) (*Menu, error)
	// CreateMenu создает новое меню
	CreateMenu(ctx context.Context, userID string, startDate, endDate time.Time) (*Menu, error)
	// GenerateMenu генерирует меню с помощью ChatGPT
	GenerateMenu(ctx context.Context, userID string, preferences map[string]interface{}) (*Menu, error)
	// AddRecipe добавляет рецепт в меню
	AddRecipe(ctx context.Context, menuID string, recipe Recipe) error
}

// Store определяет интерфейс для хранения меню
type Store interface {
	// SaveMenu сохраняет меню
	SaveMenu(ctx context.Context, menu *Menu) error
	// LoadMenu загружает меню по ID
	LoadMenu(ctx context.Context, id string) (*Menu, error)
	// LoadMenusByUser загружает все меню пользователя
	LoadMenusByUser(ctx context.Context, userID string) ([]Menu, error)
	// SaveRecipe сохраняет рецепт
	SaveRecipe(ctx context.Context, recipe *Recipe) error
}
