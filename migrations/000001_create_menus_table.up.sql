-- migrate -database "mysql://menu_manager:menu_manager@tcp(localhost:3306)/menu_test" -path /Users/doblakov/magistratura/hw/go/havchik_podbirator/menu_manager/migrations up

CREATE TABLE menu (
    meal_id VARCHAR(36) PRIMARY KEY,
    eat_date TIMESTAMP NOT NULL,
    user_id VARCHAR(36) NOT NULL
);

CREATE TABLE dishes (
    meal_id VARCHAR(36) NOT NULL,
    dish_id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(255) NOT NULL,
    recipie JSON NOT NULL,
    total_nutrition JSON NOT NULL
);


-- Insert example data
INSERT INTO menu (meal_id, eat_date, user_id) VALUES
('1', '2024-03-20 08:00:00', 'kolya'),
('2', '2024-03-20 13:00:00', 'kolya'),
('3', '2024-03-20 19:00:00', 'dan');

INSERT INTO dishes (meal_id, dish_id, name, type, recipie, total_nutrition) VALUES
('1', '1', 'Овсяная каша', 'breakfast', 
    JSON_OBJECT(
        'ingredients', JSON_ARRAY(
            JSON_OBJECT(
                'product_id', 'овсяные_хлопья',
                'amount', 100,
                'unit', 'г'
            ),
            JSON_OBJECT(
                'product_id', 'молоко',
                'amount', 200,
                'unit', 'мл'
            )
        ),
        'steps', JSON_ARRAY(
            'Вскипятить молоко',
            'Добавить хлопья',
            'Варить 5 минут'
        )
    ),
    JSON_OBJECT(
        'calories', 350,
        'proteins', 12,
        'fats', 7,
        'carbohydrates', 55
    )
),
('2', '2', 'Куриный суп', 'lunch',
    JSON_OBJECT(
        'ingredients', JSON_ARRAY(
            JSON_OBJECT(
                'product_id', 'куриное_филе',
                'amount', 200,
                'unit', 'г'
            ),
            JSON_OBJECT(
                'product_id', 'морковь',
                'amount', 100,
                'unit', 'г'
            )
        ),
        'steps', JSON_ARRAY(
            'Сварить бульон',
            'Добавить овощи',
            'Варить до готовности'
        )
    ),
    JSON_OBJECT(
        'calories', 450,
        'proteins', 35,
        'fats', 12,
        'carbohydrates', 25
    )
),
('2', '4', 'Рататуй', 'lunch',
    JSON_OBJECT(
        'ingredients', JSON_ARRAY(
            JSON_OBJECT(
                'product_id', 'крыса',
                'amount', 200,
                'unit', 'г'
            ),
            JSON_OBJECT(
                'product_id', 'помидор',
                'amount', 100,
                'unit', 'г'
            )
        ),
        'steps', JSON_ARRAY(
            'Сварить крысы',
            'Добавить помидор',
            'Варить до вечера'
        )
    ),
    JSON_OBJECT(
        'calories', 40,
        'proteins', 335,
        'fats', 121,
        'carbohydrates', 259
    )
),
('3', '3', 'Сила Земли', 'dinner',
    JSON_OBJECT(
        'ingredients', JSON_ARRAY(
            JSON_OBJECT(
                'product_id', 'огурец',
                'amount', 200,
                'unit', 'г'
            ),
            JSON_OBJECT(
                'product_id', 'морковь',
                'amount', 100,
                'unit', 'г'
            )
        ),
        'steps', JSON_ARRAY(
            'Берем молоденький огурец',
            'Надкусываем и смачиваем слюной',
            'Не отрывая от ботвы',
            'Засунуть в ...'
        )
    ),
    JSON_OBJECT(
        'calories', 100500,
        'proteins', 42,
        'fats', -10,
        'carbohydrates', 25
    )
);