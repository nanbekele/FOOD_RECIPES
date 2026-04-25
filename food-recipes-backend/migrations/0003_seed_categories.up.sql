BEGIN;
INSERT INTO categories (name, slug) VALUES
  ('Breakfast','breakfast'),
  ('Lunch','lunch'),
  ('Dinner','dinner'),
  ('Dessert','dessert'),
  ('Vegan','vegan'),
  ('Vegetarian','vegetarian'),
  ('Gluten Free','gluten-free'),
  ('Keto','keto')
ON CONFLICT (slug) DO NOTHING;
COMMIT;
