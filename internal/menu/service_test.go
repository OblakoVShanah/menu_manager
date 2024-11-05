package menu_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"menu_manager/internal/menu"
	"menu_manager/internal/menu/mock"
	"menu_manager/internal/oops"
)

func TestService_CreateMenu(t *testing.T) {
	store := mock.NewStore()
	service := menu.NewService(store)
	ctx := context.Background()

	t.Run("successful create", func(t *testing.T) {
		userID := "test-user"
		startDate := time.Now()
		endDate := startDate.Add(24 * time.Hour)

		menu, err := service.CreateMenu(ctx, userID, startDate, endDate)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if menu.UserID != userID {
			t.Errorf("Expected user ID %s, got %s", userID, menu.UserID)
		}

		if menu.StartDate != startDate {
			t.Errorf("Expected start date %v, got %v", startDate, menu.StartDate)
		}

		if menu.EndDate != endDate {
			t.Errorf("Expected end date %v, got %v", endDate, menu.EndDate)
		}
	})

	t.Run("invalid dates", func(t *testing.T) {
		userID := "test-user"
		startDate := time.Now()
		endDate := startDate.Add(-24 * time.Hour) // endDate before startDate

		_, err := service.CreateMenu(ctx, userID, startDate, endDate)
		if err == nil {
			t.Fatal("Expected error for invalid dates, got nil")
		}

		var validationErr *oops.ValidationError
		if !errors.As(err, &validationErr) {
			t.Errorf("Expected ValidationError, got %T", err)
		}
	})
}

func TestService_GetMenu(t *testing.T) {
	store := mock.NewStore()
	service := menu.NewService(store)
	ctx := context.Background()

	testMenu := &menu.Menu{
		ID:        "test-id",
		UserID:    "test-user",
		StartDate: time.Now(),
		EndDate:   time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	}

	store.SetMenu(testMenu)

	t.Run("successful get", func(t *testing.T) {
		menu, err := service.GetMenu(ctx, testMenu.ID)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if menu.ID != testMenu.ID {
			t.Errorf("Expected menu ID %s, got %s", testMenu.ID, menu.ID)
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := service.GetMenu(ctx, "non-existent-id")
		if err == nil {
			t.Fatal("Expected error for non-existent menu, got nil")
		}

		if !errors.Is(err, oops.ErrMenuNotFound) {
			t.Errorf("Expected ErrMenuNotFound, got %v", err)
		}
	})
}
