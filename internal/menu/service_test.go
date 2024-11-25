package menu_test

import (
	"context"
	menu "menu_manager/internal/menu"
	mocks "menu_manager/internal/menu/mock"
	"menu_manager/internal/oops"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetMeal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockStore(ctrl)
	mockClient := mocks.NewMockClient(ctrl)
	service := menu.NewService(mockStore, mockClient)

	ctx := context.Background()
	userID := "123"
	menuData := []menu.Menu{
		{MealID: "meal1", Time: time.Now().Add(1 * time.Hour), MealType: "lunch"},
		{MealID: "meal2", Time: time.Now().Add(2 * time.Hour), MealType: "dinner"},
	}
	expectedMeal := &menu.Meal{MealID: "meal1", Recipes: []string{"recipe1", "recipe2"}}
	expectedProducts := "product1, product2"

	mockStore.EXPECT().LoadMenu(ctx, userID).Return(menuData, nil)
	mockStore.EXPECT().LoadMeal(ctx, "meal1").Return(expectedMeal, nil)
	mockClient.EXPECT().GetProducts(ctx, expectedMeal.Recipes).Return(expectedProducts, nil)

	meal, products, err := service.GetMeal(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedMeal, meal)
	assert.Equal(t, expectedProducts, products)
}

func TestIsActual(t *testing.T) {
	now := time.Now()
	menuData := []menu.Menu{
		{MealID: "meal1", Time: now, MealType: "breakfast"},
		{MealID: "meal2", Time: now.Add(-24 * time.Hour), MealType: "dinner"},
	}

	assert.True(t, menu.IsActual(menuData))
	assert.False(t, menu.IsActual([]menu.Menu{
		{MealID: "meal3", Time: now.Add(-24 * time.Hour), MealType: "lunch"},
	}))
}

func TestFindClosestMeal(t *testing.T) {
	now := time.Now()
	menuData := []menu.Menu{
		{MealID: "meal1", Time: now.Add(1 * time.Hour), MealType: "lunch"},
		{MealID: "meal2", Time: now.Add(2 * time.Hour), MealType: "dinner"},
	}

	mealID, err := menu.FindClosestMeal(menuData)
	assert.NoError(t, err)
	assert.Equal(t, "meal1", mealID)

	_, err = menu.FindClosestMeal([]menu.Menu{})
	assert.ErrorIs(t, err, oops.ErrInvalidDates)
}

func TestRescheduleMenu(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockStore(ctrl)
	service := menu.NewService(mockStore, nil)

	ctx := context.Background()
	userID := "123"
	menuData := []menu.Menu{
		{MealID: "meal1", Time: time.Now(), MealType: "lunch"},
		{MealID: "meal2", Time: time.Now().Add(1 * time.Hour), MealType: "dinner"},
	}

	mockStore.EXPECT().UpdateMenu(ctx, userID, gomock.Any()).Return(nil)

	updatedMenu, err := service.RescheduleMenu(ctx, menuData, userID)
	assert.NoError(t, err)
	assert.Len(t, updatedMenu, len(menuData))

	for _, v := range updatedMenu {
		assert.True(t, v.Time.After(time.Now()))
	}
}

func TestGetMenu(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockStore(ctrl)
	service := menu.NewService(mockStore, nil)

	ctx := context.Background()
	userID := "123"
	expectedMenu := []menu.Menu{
		{MealID: "meal1", Time: time.Now(), MealType: "breakfast"},
	}

	mockStore.EXPECT().LoadMenu(ctx, userID).Return(expectedMenu, nil)

	menu, err := service.GetMenu(ctx, userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedMenu, menu)
}
