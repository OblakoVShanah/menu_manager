package mysql_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"menu_manager/internal/menu"
	"menu_manager/internal/menu/mysql"
	common "menu_manager/internal/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestLoadMenu_Success(t *testing.T) {
	// Создаем мок для sql.DB
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Оборачиваем *sql.DB в *sqlx.DB
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	mockRows := sqlmock.NewRows([]string{"meal_id", "eat_date", "meal_type"}).
		AddRow("meal1", time.Now(), "lunch").
		AddRow("meal2", time.Now().Add(1*time.Hour), "dinner")

	mock.ExpectQuery(`SELECT meal_id, eat_date, meal_type FROM menu WHERE user_id = \$1`).
		WithArgs("123").
		WillReturnRows(mockRows)

	storage := mysql.NewStorage(sqlxDB)

	menus, err := storage.LoadMenu(context.Background(), "123")
	assert.NoError(t, err)
	assert.Len(t, menus, 2)
	assert.Equal(t, "meal1", menus[0].MealID)
	assert.Equal(t, "lunch", menus[0].MealType)
}

func TestLoadMenu_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	mock.ExpectQuery(`SELECT meal_id, eat_date, meal_type FROM menu WHERE user_id = \$1`).
		WithArgs("123").
		WillReturnError(sql.ErrConnDone)

	storage := mysql.NewStorage(sqlxDB)

	_, err = storage.LoadMenu(context.Background(), "123")
	assert.Error(t, err)
}

func TestLoadMeal_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	nutrition := common.NutritionalValueAbsolute{
		Proteins:      10,
		Fats:          5,
		Carbohydrates: 20,
		Calories:      200,
	}
	nutritionJSON, _ := json.Marshal(nutrition)

	mockRows := sqlmock.NewRows([]string{"dish_id", "name", "recipe", "total_nutrition"}).
		AddRow("dish1", "Pasta", "recipe1", nutritionJSON).
		AddRow("dish2", "Salad", "recipe2", nutritionJSON)

	mock.ExpectQuery(`SELECT dish_id, name, recipe, total_nutrition FROM dishes WHERE meal_id = \$1`).
		WithArgs("meal1").
		WillReturnRows(mockRows)

	storage := mysql.NewStorage(sqlxDB)

	meal, err := storage.LoadMeal(context.Background(), "meal1")
	assert.NoError(t, err)
	assert.Equal(t, "meal1", meal.MealID)
	assert.Len(t, meal.DishIDs, 2)
	assert.Equal(t, nutrition.AddAbsoluteValue(nutrition), meal.TotalNutrition)
}

func TestLoadMeal_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	mock.ExpectQuery(`SELECT dish_id, name, recipe, total_nutrition FROM dishes WHERE meal_id = \$1`).
		WithArgs("meal1").
		WillReturnError(sql.ErrConnDone)

	storage := mysql.NewStorage(sqlxDB)

	_, err = storage.LoadMeal(context.Background(), "meal1")
	assert.Error(t, err)
}

func TestUpdateMenu_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE menu SET eat_date = \$1 WHERE user_id = \$2 AND meal_id = \$3`).
		WithArgs(sqlmock.AnyArg(), "123", "meal1").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	storage := mysql.NewStorage(sqlxDB)

	menus := []menu.Menu{
		{MealID: "meal1", Time: time.Now()},
	}

	err = storage.UpdateMenu(context.Background(), "123", menus)
	assert.NoError(t, err)
}

func TestUpdateMenu_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE menu SET eat_date = \$1 WHERE userID = \$2 AND meal_id = \$3`).
		WithArgs(sqlmock.AnyArg(), "123", "meal1").
		WillReturnError(sql.ErrConnDone)
	mock.ExpectRollback()

	storage := mysql.NewStorage(sqlxDB)

	menus := []menu.Menu{
		{MealID: "meal1", Time: time.Now()},
	}

	err = storage.UpdateMenu(context.Background(), "123", menus)
	assert.Error(t, err)
}
