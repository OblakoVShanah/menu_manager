package postgres

import (
	"context"
	"database/sql"
	"fmt"

	common "menu_manager/internal/models"
	"menu_manager/internal/storage"

	"github.com/jmoiron/sqlx"
)

type dbFoodProduct struct {
	ID                    sql.NullString  `db:"id"`
	Name                  sql.NullString  `db:"name"`
	WeightPerPkg          sql.NullInt64   `db:"weight_per_pkg"`
	Amount                sql.NullInt64   `db:"amount"`
	PricePerPkg           sql.NullFloat64 `db:"price_per_pkg"`
	ExpirationDate        sql.NullTime    `db:"expiration_date"`
	PresentInFridge       sql.NullBool    `db:"present_in_fridge"`
	ProteinRelative       sql.NullInt64   `db:"protein_relative"`
	FatRelative           sql.NullInt64   `db:"fat_relative"`
	CarbohydratesRelative sql.NullInt64   `db:"carbohydrates_relative"`
}

type Storage struct {
	db *sqlx.DB
}

func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) LoadProducts(ctx context.Context) ([]storage.FoodProduct, error) {
	var dbProducts []dbFoodProduct
	query := `
		SELECT id, name, weight_per_pkg, amount, price_per_pkg, expiration_date, 
		       present_in_fridge, protein_relative, fat_relative, carbohydrates_relative
		FROM food_products
	`
	err := s.db.SelectContext(ctx, &dbProducts, query)
	if err != nil {
		return nil, fmt.Errorf("failed to select products: %w", err)
	}

	products := make([]storage.FoodProduct, 0, len(dbProducts))
	for _, dbp := range dbProducts {
		product, err := convertToFoodProduct(dbp)
		if err != nil {
			return nil, fmt.Errorf("failed to convert product: %w", err)
		}
		products = append(products, product)
	}

	return products, nil
}

func (s *Storage) SaveProduct(ctx context.Context, product storage.FoodProduct) (id string, err error) {
	dbp, err := convertToDBFoodProduct(product)
	if err != nil {
		return "", fmt.Errorf("failed to convert product to DB format: %w", err)
	}

	query := `
		INSERT INTO food_products (
			id, name, weight_per_pkg, amount, price_per_pkg, expiration_date,
			present_in_fridge, protein_relative, fat_relative, carbohydrates_relative
		) VALUES (
			:id, :name, :weight_per_pkg, :amount, :price_per_pkg, :expiration_date,
			:present_in_fridge, :protein_relative, :fat_relative, :carbohydrates_relative
		) RETURNING id
	`

	rows, err := s.db.NamedQueryContext(ctx, query, dbp)
	if err != nil {
		return "", fmt.Errorf("failed to insert product: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return "", fmt.Errorf("failed to scan returned id: %w", err)
		}
	} else {
		return "", fmt.Errorf("no id returned after insert")
	}

	return id, nil
}

func convertToFoodProduct(dbp dbFoodProduct) (storage.FoodProduct, error) {
	if !dbp.ID.Valid {
		return storage.FoodProduct{}, fmt.Errorf("invalid ID")
	}

	weightPerPkg := uint(dbp.WeightPerPkg.Int64)
	if !dbp.WeightPerPkg.Valid || weightPerPkg == 0 {
		return storage.FoodProduct{}, fmt.Errorf("invalid WeightPerPkg")
	}

	amount := uint(dbp.Amount.Int64)
	if !dbp.Amount.Valid {
		return storage.FoodProduct{}, fmt.Errorf("invalid Amount")
	}

	if !dbp.PricePerPkg.Valid {
		return storage.FoodProduct{}, fmt.Errorf("invalid PricePerPkg")
	}

	if !dbp.ExpirationDate.Valid {
		return storage.FoodProduct{}, fmt.Errorf("invalid ExpirationDate")
	}

	if !dbp.PresentInFridge.Valid {
		return storage.FoodProduct{}, fmt.Errorf("invalid PresentInFridge")
	}

	if !dbp.ProteinRelative.Valid || !dbp.FatRelative.Valid || !dbp.CarbohydratesRelative.Valid {
		return storage.FoodProduct{}, fmt.Errorf("invalid NutritionalValueRelative")
	}

	return storage.FoodProduct{
		ID:              dbp.ID.String,
		Name:            dbp.Name.String,
		WeightPerPkg:    weightPerPkg,
		Amount:          amount,
		PricePerPkg:     float32(dbp.PricePerPkg.Float64),
		ExpirationDate:  dbp.ExpirationDate.Time,
		PresentInFridge: dbp.PresentInFridge.Bool,
		NutritionalValueRelative: common.NutritionalValueRelative{
			Proteins:      uint(dbp.ProteinRelative.Int64),
			Fats:          uint(dbp.FatRelative.Int64),
			Carbohydrates: uint(dbp.CarbohydratesRelative.Int64),
		},
	}, nil
}

func convertToDBFoodProduct(fp storage.FoodProduct) (dbFoodProduct, error) {
	if fp.ID == "" {
		return dbFoodProduct{}, fmt.Errorf("invalid ID")
	}

	if fp.WeightPerPkg == 0 {
		return dbFoodProduct{}, fmt.Errorf("invalid WeightPerPkg")
	}

	return dbFoodProduct{
		ID:                    sql.NullString{String: fp.ID, Valid: true},
		Name:                  sql.NullString{String: fp.Name, Valid: true},
		WeightPerPkg:          sql.NullInt64{Int64: int64(fp.WeightPerPkg), Valid: true},
		Amount:                sql.NullInt64{Int64: int64(fp.Amount), Valid: true},
		PricePerPkg:           sql.NullFloat64{Float64: float64(fp.PricePerPkg), Valid: true},
		ExpirationDate:        sql.NullTime{Time: fp.ExpirationDate, Valid: true},
		PresentInFridge:       sql.NullBool{Bool: fp.PresentInFridge, Valid: true},
		ProteinRelative:       sql.NullInt64{Int64: int64(fp.NutritionalValueRelative.Proteins), Valid: true},
		FatRelative:           sql.NullInt64{Int64: int64(fp.NutritionalValueRelative.Fats), Valid: true},
		CarbohydratesRelative: sql.NullInt64{Int64: int64(fp.NutritionalValueRelative.Carbohydrates), Valid: true},
	}, nil
}
