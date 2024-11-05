package mock

import (
	"context"
	"menu_manager/internal/menu"
	"menu_manager/internal/oops"
)

type Store struct {
	menus   map[string]*menu.Menu
	recipes map[string]*menu.Recipe
	err     error
}

func NewStore() *Store {
	return &Store{
		menus:   make(map[string]*menu.Menu),
		recipes: make(map[string]*menu.Recipe),
	}
}

func (s *Store) SetError(err error) {
	s.err = err
}

func (s *Store) SetMenu(m *menu.Menu) {
	s.menus[m.ID] = m
}

func (s *Store) SaveMenu(ctx context.Context, menu *menu.Menu) error {
	if s.err != nil {
		return s.err
	}
	s.menus[menu.ID] = menu
	return nil
}

func (s *Store) LoadMenu(ctx context.Context, id string) (*menu.Menu, error) {
	if s.err != nil {
		return nil, s.err
	}
	menu, exists := s.menus[id]
	if !exists {
		return nil, oops.ErrMenuNotFound
	}
	return menu, nil
}

func (s *Store) LoadMenusByUser(ctx context.Context, userID string) ([]menu.Menu, error) {
	if s.err != nil {
		return nil, s.err
	}
	var menus []menu.Menu
	for _, m := range s.menus {
		if m.UserID == userID {
			menus = append(menus, *m)
		}
	}
	return menus, nil
}

func (s *Store) SaveRecipe(ctx context.Context, recipe *menu.Recipe) error {
	if s.err != nil {
		return s.err
	}
	s.recipes[recipe.ID] = recipe
	return nil
}
