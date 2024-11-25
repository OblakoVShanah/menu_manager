package menu_test

import (
	"encoding/json"
	"errors"
	"menu_manager/internal/menu"
	mocks "menu_manager/internal/menu/mock"
	common "menu_manager/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetMeal_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок для Service
	mockService := mocks.NewMockService(ctrl)

	// Данные для теста
	expectedMeal := menu.Meal{
		MealID:         "meal1",
		DishIDs:        []string{"eggs", "bread"},
		DishNames:      []string{"eggs", "bread"},
		Type:           "Breakfast",
		Recipes:        []string{"eggs", "bread"},
		TotalNutrition: common.NutritionalValueAbsolute{Proteins: 10, Fats: 10, Carbohydrates: 30, Calories: 300},
	}
	expectedProducts := "eggs, bread"
	mockService.EXPECT().GetMeal(gomock.Any(), "123").Return(&expectedMeal, expectedProducts, nil)

	// Создаем HTTP-реквест и респонс
	router := chi.NewRouter()
	handler := menu.NewHandler(router, mockService)
	handler.Register()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/menus/getMeal", nil)
	rec := httptest.NewRecorder()

	// Выполняем запрос
	router.ServeHTTP(rec, req)

	// Проверяем результат
	assert.Equal(t, http.StatusOK, rec.Code)

	var response struct {
		Meal         menu.Meal `json:"meal"`
		ShoppingList string    `json:"shopping_list"`
	}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedMeal, response.Meal)
	assert.Equal(t, expectedProducts, response.ShoppingList)
}

func TestGetMeal_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок для Service
	mockService := mocks.NewMockService(ctrl)

	// Настройка мока для ошибки
	mockService.EXPECT().GetMeal(gomock.Any(), "123").Return(nil, "", errors.New("service error"))

	// Создаем HTTP-реквест и респонс
	router := chi.NewRouter()
	handler := menu.NewHandler(router, mockService)
	handler.Register()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/menus/getMeal", nil)
	rec := httptest.NewRecorder()

	// Выполняем запрос
	router.ServeHTTP(rec, req)

	// Проверяем результат
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "service error")
}
