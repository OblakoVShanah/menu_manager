package postgres

import (
	"context"
	"database/sql"
	"encoding/json"

	"menu_manager/internal/menu"
	"menu_manager/internal/oops"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db *sqlx.DB
}

func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) SaveMenu(ctx context.Context, menu *menu.Menu) error {
	query := `
		INSERT INTO menus (id, user_id, start_date, end_date, meals, created_at, generated_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE
		SET user_id = EXCLUDED.user_id,
			start_date = EXCLUDED.start_date,
			end_date = EXCLUDED.end_date,
			meals = EXCLUDED.meals,
			generated_by = EXCLUDED.generated_by
	`

	mealsJSON, err := json.Marshal(menu.Meals)
	if err != nil {
		return oops.NewDBError(err, "SaveMenu.Marshal", menu.ID)
	}

	_, err = s.db.ExecContext(ctx, query,
		menu.ID,
		menu.UserID,
		menu.StartDate,
		menu.EndDate,
		mealsJSON,
		menu.CreatedAt,
		menu.GeneratedBy,
	)

	if err != nil {
		return oops.NewDBError(err, "SaveMenu", menu.ID)
	}

	return nil
}

func (s *Storage) LoadMenu(ctx context.Context, id string) (*menu.Menu, error) {
	query := `
		SELECT id, user_id, start_date, end_date, meals, created_at, generated_by
		FROM menus
		WHERE id = $1
	`

	var m menu.Menu
	var mealsJSON []byte

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&m.ID,
		&m.UserID,
		&m.StartDate,
		&m.EndDate,
		&mealsJSON,
		&m.CreatedAt,
		&m.GeneratedBy,
	)

	if err == sql.ErrNoRows {
		return nil, oops.ErrMenuNotFound
	}
	if err != nil {
		return nil, oops.NewDBError(err, "LoadMenu", id)
	}

	if err := json.Unmarshal(mealsJSON, &m.Meals); err != nil {
		return nil, oops.NewDBError(err, "LoadMenu.Unmarshal", id)
	}

	return &m, nil
}

func (s *Storage) LoadMenusByUser(ctx context.Context, userID string) ([]menu.Menu, error) {
	query := `
		SELECT id, user_id, start_date, end_date, meals, created_at, generated_by
		FROM menus
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, oops.NewDBError(err, "LoadMenusByUser", userID)
	}
	defer rows.Close()

	var menus []menu.Menu
	for rows.Next() {
		var m menu.Menu
		var mealsJSON []byte

		err := rows.Scan(
			&m.ID,
			&m.UserID,
			&m.StartDate,
			&m.EndDate,
			&mealsJSON,
			&m.CreatedAt,
			&m.GeneratedBy,
		)
		if err != nil {
			return nil, oops.NewDBError(err, "LoadMenusByUser.Scan", userID)
		}

		if err := json.Unmarshal(mealsJSON, &m.Meals); err != nil {
			return nil, oops.NewDBError(err, "LoadMenusByUser.Unmarshal", userID)
		}

		menus = append(menus, m)
	}

	if err = rows.Err(); err != nil {
		return nil, oops.NewDBError(err, "LoadMenusByUser.Rows", userID)
	}

	return menus, nil
}

func (s *Storage) SaveRecipe(ctx context.Context, recipe *menu.Recipe) error {
	query := `
		INSERT INTO recipes (id, name, description, steps, ingredients, nutrition, cooking_time)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE
		SET name = EXCLUDED.name,
			description = EXCLUDED.description,
			steps = EXCLUDED.steps,
			ingredients = EXCLUDED.ingredients,
			nutrition = EXCLUDED.nutrition,
			cooking_time = EXCLUDED.cooking_time
	`

	stepsJSON, err := json.Marshal(recipe.Steps)
	if err != nil {
		return oops.NewDBError(err, "SaveRecipe.MarshalSteps", recipe.ID)
	}

	ingredientsJSON, err := json.Marshal(recipe.Ingredients)
	if err != nil {
		return oops.NewDBError(err, "SaveRecipe.MarshalIngredients", recipe.ID)
	}

	nutritionJSON, err := json.Marshal(recipe.Nutrition)
	if err != nil {
		return oops.NewDBError(err, "SaveRecipe.MarshalNutrition", recipe.ID)
	}

	_, err = s.db.ExecContext(ctx, query,
		recipe.ID,
		recipe.Name,
		recipe.Description,
		stepsJSON,
		ingredientsJSON,
		nutritionJSON,
		recipe.CookingTime,
	)

	if err != nil {
		return oops.NewDBError(err, "SaveRecipe", recipe.ID)
	}

	return nil
}
