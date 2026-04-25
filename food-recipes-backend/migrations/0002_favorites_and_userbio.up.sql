BEGIN;

-- Ensure pgcrypto exists for UUID if not already
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Favorites table: user_id INT -> users(id), recipe_id UUID -> recipes(id)
CREATE TABLE IF NOT EXISTS favorites (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id INT NOT NULL,
  recipe_id UUID NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT fk_fav_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_fav_recipe FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE,
  CONSTRAINT uq_fav_user_recipe UNIQUE (user_id, recipe_id)
);

CREATE INDEX IF NOT EXISTS idx_fav_user_id ON favorites(user_id);
CREATE INDEX IF NOT EXISTS idx_fav_recipe_id ON favorites(recipe_id);

-- Users.bio for profile pages
ALTER TABLE users ADD COLUMN IF NOT EXISTS bio TEXT;

COMMIT;
