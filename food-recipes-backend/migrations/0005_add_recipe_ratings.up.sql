BEGIN;

-- 1) Add average_rating to recipes
ALTER TABLE recipes
  ADD COLUMN IF NOT EXISTS average_rating NUMERIC(2,1) NOT NULL DEFAULT 0;

-- 2) Create recipe_ratings table (1-5 stars per user per recipe)
CREATE TABLE IF NOT EXISTS recipe_ratings (
  id SERIAL PRIMARY KEY,
  user_id INT NOT NULL,
  recipe_id UUID NOT NULL,
  rating INT NOT NULL CHECK (rating >= 1 AND rating <= 5),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT fk_rr_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_rr_recipe FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE,
  CONSTRAINT uq_user_recipe_rating UNIQUE (user_id, recipe_id)
);

CREATE INDEX IF NOT EXISTS idx_recipe_ratings_recipe_id ON recipe_ratings(recipe_id);

-- 3) Trigger/function to update recipes.average_rating when ratings change
CREATE OR REPLACE FUNCTION public.update_recipe_average_rating() RETURNS trigger AS $$
DECLARE
  rec_id uuid;
  avg_val NUMERIC(2,1);
BEGIN
  rec_id := COALESCE(NEW.recipe_id, OLD.recipe_id);
  SELECT COALESCE(ROUND(AVG(rating)::numeric, 1), 0)::numeric(2,1) INTO avg_val
  FROM recipe_ratings
  WHERE recipe_id = rec_id;

  UPDATE recipes SET average_rating = COALESCE(avg_val, 0) WHERE id = rec_id;
  RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_update_avg_rating_insupd ON recipe_ratings;
CREATE TRIGGER trg_update_avg_rating_insupd
AFTER INSERT OR UPDATE OF rating ON recipe_ratings
FOR EACH ROW EXECUTE FUNCTION public.update_recipe_average_rating();

DROP TRIGGER IF EXISTS trg_update_avg_rating_del ON recipe_ratings;
CREATE TRIGGER trg_update_avg_rating_del
AFTER DELETE ON recipe_ratings
FOR EACH ROW EXECUTE FUNCTION public.update_recipe_average_rating();

COMMIT;
