BEGIN;

-- Drop triggers
DROP TRIGGER IF EXISTS trg_update_avg_rating_insupd ON recipe_ratings;
DROP TRIGGER IF EXISTS trg_update_avg_rating_del ON recipe_ratings;

-- Drop function
DROP FUNCTION IF EXISTS public.update_recipe_average_rating();

-- Drop ratings table
DROP TABLE IF EXISTS recipe_ratings;

-- Drop average_rating column from recipes
ALTER TABLE recipes DROP COLUMN IF EXISTS average_rating;

COMMIT;
