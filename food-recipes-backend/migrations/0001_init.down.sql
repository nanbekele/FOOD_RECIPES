BEGIN;
DROP TABLE IF EXISTS recipe_steps;
DROP TABLE IF EXISTS recipe_ingredients;
DROP TABLE IF EXISTS recipes;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS users;
-- Note: extension pgcrypto is left installed intentionally
COMMIT;
