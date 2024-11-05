package storage_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	common "menu_manager/internal/models"
	"menu_manager/internal/storage"
	"menu_manager/internal/storage/mock"

	"github.com/go-chi/chi/v5"
	"github.com/google/go-cmp/cmp"
)

func TestHandler_getProducts(t *testing.T) {
	mockStore := mock.NewStore()
	router := chi.NewRouter()
	service := storage.NewService(mockStore)
	handler := storage.NewHandler(router, service)
	handler.Register()

	testProducts := []storage.FoodProduct{
		{
			ID:              "test-id-1",
			Name:            "Test Product 1",
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
		},
	}

	t.Run("successful get", func(t *testing.T) {
		mockStore.SetProducts(testProducts)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/products", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				rr.Code, http.StatusOK)
		}

		var got []storage.FoodProduct
		err := json.NewDecoder(rr.Body).Decode(&got)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if diff := cmp.Diff(testProducts, got); diff != "" {
			t.Errorf("handler returned wrong body (-want +got):\n%s", diff)
		}
	})

	t.Run("error case", func(t *testing.T) {
		mockStore.SetError(errors.New("database error"))

		req := httptest.NewRequest(http.MethodGet, "/api/v1/products", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v",
				rr.Code, http.StatusInternalServerError)
		}
	})
}
