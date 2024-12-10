package mysql

import (
	"context"
	"encoding/json"

	"menu_manager/internal/menu"
	common "menu_manager/internal/models"
	"menu_manager/internal/oops"

	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db *sqlx.DB
}

func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{db: db}
}

// LoadMenu возвращает меню из БД со списком id приемов пиши и их запланированного времени
func (s *Storage) LoadMenu(ctx context.Context, userID string) ([]menu.Menu, error) {

	query := `
		SELECT meal_id, eat_date, meal_type
		FROM menu
		WHERE user_id = ?
	`

	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, oops.NewDBError(err, "LoadMenu", userID)
	}
	defer rows.Close()

	var menuList []menu.Menu
	for rows.Next() {
		var m menu.Menu
		err := rows.Scan(
			&m.MealID,
			&m.Time,
			&m.MealType,
		)
		if err != nil {
			return nil, oops.NewDBError(err, "LoadMenu.Scan", userID)
		}
		menuList = append(menuList, m)
	}

	return menuList, nil
}

// LoadMeal возвращает из базы прием пищи с описанием составляющих его блюд и продуктов
func (s *Storage) LoadMeal(ctx context.Context, mealID string) (*menu.Meal, error) {

	// текст запроса
	query := `
		SELECT dish_id, name, recipie, total_nutrition
		FROM dishes
		WHERE meal_id = ?
	`
	// выполняем запрос к БД
	rows, err := s.db.QueryContext(ctx, query, mealID)
	if err != nil {
		return nil, oops.NewDBError(err, "LoadMeal", mealID)
	}

	defer rows.Close()

	// парсим данные
	meal := menu.Meal{
		MealID:    mealID,
		DishIDs:   make([]string, 0, 10),
		DishNames: make([]string, 0, 10),
		Type:      "",
		Recipes:   make([]string, 0, 10),
		TotalNutrition: common.NutritionalValueAbsolute{
			Proteins:      0,
			Fats:          0,
			Carbohydrates: 0,
			Calories:      0,
		},
	}
	for rows.Next() {

		var dishID string
		var dishName string
		var recipeJson string
		var nutritionJson string

		err := rows.Scan(
			&dishID,
			&dishName,
			&recipeJson,
			&nutritionJson,
		)
		if err != nil {
			return nil, oops.NewDBError(err, "LoadMeal.Scan", mealID)
		}

		meal.DishIDs = append(meal.DishIDs, dishID)
		meal.DishNames = append(meal.DishNames, dishName)
		meal.Recipes = append(meal.Recipes, recipeJson)

		// Парсим JSON в структуру
		var nutritionalValue common.NutritionalValueAbsolute
		err = json.Unmarshal([]byte(nutritionJson), &nutritionalValue)
		if err != nil {
			return nil, oops.NewDBError(err, "LoadMeal.JsonUnmarshal", mealID)
		}
		meal.TotalNutrition = meal.TotalNutrition.AddAbsoluteValue(nutritionalValue)
	}

	if err = rows.Err(); err != nil {
		return nil, oops.NewDBError(err, "LoadMeal.Rows", mealID)
	}
	// log.Println(meal.Recipes)
	return &meal, nil
}

// UpdateMenu обновляет время и даты приемов пищи
func (s *Storage) UpdateMenu(ctx context.Context, userID string, menuList []menu.Menu) error {
	// Начинаем транзакцию
	tx, err := s.db.Beginx()
	if err != nil {
		return oops.NewDBError(err, "failed to begin transaction", userID)
	}

	// Обновляем каждую запись
	updateQuery := "UPDATE menu SET eat_date = ? WHERE user_id = ? AND meal_id = ?"
	for _, m := range menuList {
		_, err := tx.Exec(updateQuery, m.Time, userID, m.MealID)
		if err != nil {
			// При ошибке откатываем транзакцию
			tx.Rollback()
			return oops.NewDBError(err, "failed to update menu", userID)
		}
	}

	// Если все прошло успешно, фиксируем изменения
	if err := tx.Commit(); err != nil {
		return oops.NewDBError(err, "failed to commit transaction", userID)
	}

	return nil
}
