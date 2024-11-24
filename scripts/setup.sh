#!/bin/bash

# Run MySQL initialization
mysql -u root -p < scripts/init_db.sql

# Run migrations
migrate -database "mysql://menu_manager:menu_manager@tcp(localhost:3306)/menu" -path migrations up
migrate -database "mysql://menu_manager:menu_manager@tcp(localhost:3306)/menu_test" -path migrations up 