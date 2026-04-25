BEGIN;

-- Media (images/videos) per recipe
CREATE TABLE IF NOT EXISTS recipe_media (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  recipe_id UUID NOT NULL,
  url TEXT NOT NULL,
  media_type TEXT NOT NULL CHECK (media_type IN ('image','video')),
  position INT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT fk_recipe_media_recipe FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_recipe_media_recipe_id ON recipe_media(recipe_id);

COMMIT;
