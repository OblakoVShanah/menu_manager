package oops

import (
	"errors"
	"fmt"
)

var (
	// Ошибки базы данных
	ErrNoData       = errors.New("нет данных в базе данных")
	ErrDuplicateKey = errors.New("дублирование ключа")
	ErrDBConnection = errors.New("ошибка подключения к базе данных")

	// Ошибки бизнес-логики
	ErrMenuNotFound   = errors.New("меню не найдено")
	ErrRecipeNotFound = errors.New("рецепт не найден")
	ErrInvalidDates   = errors.New("некорректные даты")
	ErrNotImplemented = errors.New("функционал не реализован")
)

// ValidationError представляет ошибку валидации
type ValidationError struct {
	Field string
	Err   error
}

// DBError представляет ошибку базы данных
type DBError struct {
	Err error
	ID  string
	Op  string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("ошибка валидации поля '%s': %v", e.Field, e.Err)
}

func (e *DBError) Error() string {
	if e.ID != "" {
		return fmt.Sprintf("операция БД '%s' для ID '%s': %v", e.Op, e.ID, e.Err)
	}
	return fmt.Sprintf("операция БД '%s': %v", e.Op, e.Err)
}

func NewValidationError(field string, err error) *ValidationError {
	return &ValidationError{
		Field: field,
		Err:   err,
	}
}

func NewDBError(err error, op string, id string) *DBError {
	return &DBError{
		Err: err,
		ID:  id,
		Op:  op,
	}
}
