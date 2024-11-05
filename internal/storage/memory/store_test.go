package memory

import (
	"context"
	common "menu_manager/internal/models"
	"menu_manager/internal/storage"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestStorage_SaveAndLoadProducts(t *testing.T) {
	store := NewStorage()
	ctx := context.Background()

	// Test product
	product := storage.FoodProduct{
		ID:              "test-id",
		Name:            "Test Product",
		WeightPerPkg:    100,
		Amount:          1,
		PricePerPkg:     9.99,
		ExpirationDate:  time.Now().Add(24 * time.Hour),
		PresentInFridge: true,
		NutritionalValueRelative: common.NutritionalValueRelative{
			Proteins:      10,
			Fats:          20,
			Carbohydrates: 30,
			Calories:      400,
		},
	}

	// Test SaveProduct
	t.Run("save product", func(t *testing.T) {
		id, err := store.SaveProduct(ctx, product)
		if err != nil {
			t.Fatalf("SaveProduct failed: %v", err)
		}
		if id != product.ID {
			t.Errorf("SaveProduct returned wrong ID, want %s, got %s", product.ID, id)
		}
	})

	// Test LoadProducts
	t.Run("load products", func(t *testing.T) {
		products, err := store.LoadProducts(ctx)
		if err != nil {
			t.Fatalf("LoadProducts failed: %v", err)
		}

		if len(products) != 1 {
			t.Fatalf("LoadProducts returned wrong number of products, want 1, got %d", len(products))
		}

		if diff := cmp.Diff(product, products[0]); diff != "" {
			t.Errorf("LoadProducts returned wrong product (-want +got):\n%s", diff)
		}
	})
}
