CREATE DATABASE IF NOT EXISTS menu;
CREATE DATABASE IF NOT EXISTS menu_test;

-- Create user if not exists (MySQL 8.0+)
CREATE USER IF NOT EXISTS 'menu_manager'@'localhost' IDENTIFIED BY 'menu_manager';

-- Grant privileges
GRANT ALL PRIVILEGES ON menu.* TO 'menu_manager'@'localhost';
GRANT ALL PRIVILEGES ON menu_test.* TO 'menu_manager'@'localhost';
