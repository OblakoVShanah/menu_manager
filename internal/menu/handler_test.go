package menu_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"menu_manager/internal/menu"
	"menu_manager/internal/menu/mock"
	"github.com/go-chi/chi/v5"
	"github.com/google/go-cmp/cmp"
)

func TestHandler_CreateMenu(t *testing.T) {
	mockStore := mock.NewStore()
	router := chi.NewRouter()
	service := menu.NewService(mockStore)
	handler := menu.NewHandler(router, service)
	handler.Register()

	t.Run("successful create", func(t *testing.T) {
		req := menu.CreateMenuRequest{
			UserID:    "test-user",
			StartDate: time.Now(),
			EndDate:   time.Now().Add(24 * time.Hour),
		}

		body, err := json.Marshal(req)
		if err != nil {
			t.Fatalf("Failed to marshal request: %v", err)
		}

		r := httptest.NewRequest(http.MethodPost, "/api/v1/menus", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, r)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
		}

		var response menu.Menu
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if response.UserID != req.UserID {
			t.Errorf("Expected user ID %s, got %s", req.UserID, response.UserID)
		}
	})

	t.Run("invalid request", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, "/api/v1/menus", bytes.NewBufferString("invalid json"))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, r)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestHandler_GetMenu(t *testing.T) {
	mockStore := mock.NewStore()
	router := chi.NewRouter()
	service := menu.NewService(mockStore)
	handler := menu.NewHandler(router, service)
	handler.Register()

	testMenu := &menu.Menu{
		ID:        "test-id",
		UserID:    "test-user",
		StartDate: time.Now(),
		EndDate:   time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	}

	mockStore.SetMenu(testMenu)

	t.Run("successful get", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/api/v1/menus/"+testMenu.ID, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		var response menu.Menu
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if diff := cmp.Diff(testMenu, &response); diff != "" {
			t.Errorf("Response mismatch (-want +got):\n%s", diff)
		}
	})
}
