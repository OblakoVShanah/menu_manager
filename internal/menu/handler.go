package menu

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Handler обрабатывает HTTP-запросы для работы с меню
type Handler struct {
	router  *chi.Mux
	service Service
}

// NewHandler создает новый обработчик HTTP-запросов
func NewHandler(router *chi.Mux, service Service) *Handler {
	return &Handler{
		router:  router,
		service: service,
	}
}

// Register регистрирует все обработчики маршрутов
func (h *Handler) Register() {
	h.router.Route("/api/v1", func(r chi.Router) {
		r.Get("/menus/getMeal", h.getMeal)
	})
}

// getMeal получает описание следующего приема пиши и список продуктов, которые нужно докупить
func (h *Handler) getMeal(w http.ResponseWriter, r *http.Request) {

	// Извлекаем userID из query-параметров
	userID := string(r.URL.Query().Get("user_id"))
	if userID == "" {
		http.Error(w, "missing user_id parameter", http.StatusBadRequest)
		return
	}

	meal, products, err := h.service.GetMeal(r.Context(), userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Meal         Meal   `json:"meal"`
		ShoppingList string `json:"shopping_list"`
	}{
		Meal:         *meal,
		ShoppingList: products,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	log.Println(response)
}
