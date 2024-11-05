package common

// NutritionalValueRelative представляет относительную пищевую ценность продукта (на 100г)
type NutritionalValueRelative struct {
	Proteins      uint `json:"proteins"`      // Белки в граммах
	Fats          uint `json:"fats"`          // Жиры в граммах
	Carbohydrates uint `json:"carbohydrates"` // Углеводы в граммах
	Calories      uint `json:"calories"`      // Калории
}

// NutritionalValueAbsolute представляет абсолютную пищевую ценность продукта
type NutritionalValueAbsolute struct {
	Proteins      uint `json:"proteins"`      // Белки в граммах
	Fats          uint `json:"fats"`          // Жиры в граммах
	Carbohydrates uint `json:"carbohydrates"` // Углеводы в граммах
	Calories      uint `json:"calories"`      // Калории
}
