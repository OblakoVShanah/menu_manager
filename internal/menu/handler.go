package menu

import (
	"encoding/json"
	"net/http"
	"time"

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
		r.Get("/menus/{id}", h.getMenu)
		r.Post("/menus", h.createMenu)
		r.Post("/menus/generate", h.generateMenu)
		r.Post("/menus/{menuID}/recipes", h.addRecipe)
	})
}

// CreateMenuRequest представляет запрос на создание меню
type CreateMenuRequest struct {
	UserID    string    `json:"user_id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

// getMenu получает меню из БД и записывает ответ в http.ResponseWriter
func (h *Handler) getMenu(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	menu, err := h.service.GetMenu(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(menu)
}

// createMenu создает новое меню в базе
func (h *Handler) createMenu(w http.ResponseWriter, r *http.Request) {
	var req CreateMenuRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	menu, err := h.service.CreateMenu(r.Context(), req.UserID, req.StartDate, req.EndDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(menu)
}

// generateMenu генерирует меню через chatGPT
func (h *Handler) generateMenu(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID      string                 `json:"user_id"`
		Preferences map[string]interface{} `json:"preferences"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	menu, err := h.service.GenerateMenu(r.Context(), req.UserID, req.Preferences)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(menu)
}

// addRecipe добавляет новый рецепт в меню
func (h *Handler) addRecipe(w http.ResponseWriter, r *http.Request) {
	menuID := chi.URLParam(r, "menuID")

	var recipe Recipe
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.AddRecipe(r.Context(), menuID, recipe); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
