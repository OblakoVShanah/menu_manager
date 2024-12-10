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

// Product представляет собой объект продукт питания с соответствующими характеристиками
type Product struct {
	ID                       string                   `json:"id"`
	Name                     string                   `json:"name"`
	WeightPerPkg             int                      `json:"weight_per_pkg"`
	Amount                   int                      `json:"amount"`
	PricePerPkg              int                      `json:"price_per_pkg"`
	ExpirationDate           string                   `json:"expiration_date"`
	PresentInFridge          bool                     `json:"present_in_fridge"`
	NutritionalValueRelative NutritionalValueRelative `json:"nutritional_value_relative"`
}

// AddAbsoluteValue добавляет пищевую ценность переданного класса к текущему, возвращает новый экземляр класса
func (nv_left NutritionalValueAbsolute) AddAbsoluteValue(nv_right NutritionalValueAbsolute) NutritionalValueAbsolute {
	return NutritionalValueAbsolute{
		Proteins:      nv_left.Proteins + nv_right.Proteins,
		Fats:          nv_left.Fats + nv_right.Fats,
		Carbohydrates: nv_left.Carbohydrates + nv_right.Carbohydrates,
		Calories:      nv_left.Calories + nv_right.Calories,
	}
}
