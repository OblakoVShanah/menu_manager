package menu_test

import (
	"context"
	"encoding/json"
	"fmt"
	"menu_manager/internal/menu"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetProducts_Success(t *testing.T) {
	// Создаем мок-сервер
	mockResponse := `"Eggs, Bread, Milk"`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, что запрос отправлен корректно
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/v1/products", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// Проверяем тело запроса
		var recipes []string
		err := json.NewDecoder(r.Body).Decode(&recipes)
		assert.NoError(t, err)
		assert.ElementsMatch(t, []string{"recipe1", "recipe2"}, recipes)

		// Отправляем успешный ответ
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	// Создаем клиента
	client := menu.NewClient(server.URL)

	// Вызываем метод GetProducts
	recipes := []string{"recipe1", "recipe2"}
	products, err := client.GetProducts(context.Background(), recipes)

	// Проверяем результат
	assert.NoError(t, err)
	assert.Equal(t, "Eggs, Bread, Milk", products)
}

func TestGetProducts_MarshalError(t *testing.T) {
	originalMarshal := menu.JsonMarshal
	defer func() { menu.JsonMarshal = originalMarshal }()

	menu.JsonMarshal = func(v interface{}) ([]byte, error) {
		return nil, fmt.Errorf("mock marshal error")
	}

	client := menu.NewClient("http://example.com")

	_, err := client.GetProducts(context.Background(), []string{"valid_string"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to marshal product")
}

func TestGetProducts_BadStatusCode(t *testing.T) {
	// Создаем мок-сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request"))
	}))
	defer server.Close()

	// Создаем клиента
	client := menu.NewClient(server.URL)

	// Вызываем метод GetProducts
	recipes := []string{"recipe1", "recipe2"}
	_, err := client.GetProducts(context.Background(), recipes)

	// Проверяем, что ошибка корректно обработана
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected status code: 400")
}

func TestGetProducts_DecodeError(t *testing.T) {
	// Создаем мок-сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{invalid-json"))
	}))
	defer server.Close()

	// Создаем клиента
	client := menu.NewClient(server.URL)

	// Вызываем метод GetProducts
	recipes := []string{"recipe1", "recipe2"}
	_, err := client.GetProducts(context.Background(), recipes)

	// Проверяем, что ошибка корректно обработана
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to decode response")
}
