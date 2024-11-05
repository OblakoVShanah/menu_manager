package postgres_test

import (
	"context"
	"testing"
	"time"

	"menu_manager/internal/menu"
	"menu_manager/internal/menu/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sqlx.DB {
	db, err := sqlx.Connect("postgres", "postgres://test:test@localhost:5432/testdb?sslmode=disable")
	require.NoError(t, err)
	return db
}

func TestStorage_SaveAndLoadMenu(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	storage := postgres.NewStorage(db)
	ctx := context.Background()

	testMenu := &menu.Menu{
		ID:        "test-id",
		UserID:    "test-user",
		StartDate: time.Now().Truncate(time.Second),
		EndDate:   time.Now().Add(24 * time.Hour).Truncate(time.Second),
		CreatedAt: time.Now().Truncate(time.Second),
		GeneratedBy: "test",
	}

	t.Run("save menu", func(t *testing.T) {
		err := storage.SaveMenu(ctx, testMenu)
		require.NoError(t, err)
	})

	t.Run("load menu", func(t *testing.T) {
		loaded, err := storage.LoadMenu(ctx, testMenu.ID)
		require.NoError(t, err)
		require.NotNil(t, loaded)

		require.Equal(t, testMenu.ID, loaded.ID)
		require.Equal(t, testMenu.UserID, loaded.UserID)
		require.Equal(t, testMenu.StartDate, loaded.StartDate)
		require.Equal(t, testMenu.EndDate, loaded.EndDate)
		require.Equal(t, testMenu.GeneratedBy, loaded.GeneratedBy)
	})

	t.Run("load menus by user", func(t *testing.T) {
		menus, err := storage.LoadMenusByUser(ctx, testMenu.UserID)
		require.NoError(t, err)
		require.NotEmpty(t, menus)

		found := false
		for _, m := range menus {
			if m.ID == testMenu.ID {
				found = true
				require.Equal(t, testMenu.UserID, m.UserID)
				break
			}
		}
		require.True(t, found, "Saved menu not found in user's menus")
	})
}
