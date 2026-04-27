# 🍽️ FOOD_RECIPES — Full System Analysis

## 📐 Architecture Overview

This is a **full-stack food recipe web application** built with a modern 4-tier architecture:

```
┌─────────────────────────────────────────────────────────────┐
│                   Browser (Client)                          │
│              Nuxt 3 SPA (port 3000)                         │
│   Vue 3 + Tailwind CSS + Apollo Client (GraphQL)            │
└──────────────┬──────────────────────────┬───────────────────┘
               │  REST API calls           │ GraphQL queries
               ▼                          ▼
┌──────────────────────┐     ┌────────────────────────────────┐
│  Go/Gin Backend       │     │  Hasura GraphQL Engine          │
│  (port 8081)          │     │  (port 8082)                   │
│  - Auth (JWT)         │     │  - Auto-generated GraphQL API  │
│  - Image Upload       │     │  - Row-level security via JWT  │
│  - Admin ops          │     │  - Direct DB access            │
└──────────┬────────────┘     └───────────────┬────────────────┘
           │ Postgres DSN                      │ Postgres DSN
           └──────────────────┬────────────────┘
                              ▼
                  ┌────────────────────────┐
                  │   PostgreSQL 15         │
                  │   (port 5433 → 5432)   │
                  │   Database: foodrecipes │
                  └────────────────────────┘
```

**External Service:**
- **Cloudinary** — Image/video uploads (stored URLs in DB, not files)

---

## 🗄️ Database Schema (8 Migrations)

| Migration | What it adds |
|-----------|-------------|
| `0001_init` | `users`, `categories`, `recipes`, `recipe_ingredients`, `recipe_steps` |
| `0002_favorites_and_userbio` | `favorites` table, `bio` column on users |
| `0003_seed_categories` | Seeds default food categories |
| `0004_recipe_media` | `recipe_media` table (images/videos per recipe) |
| `0005_add_recipe_ratings` | `recipe_ratings` table + DB trigger for `average_rating` auto-update |
| `0006_add_user_role` | `role` + `is_verified` columns on users |
| `0007_follows` | `follows` table (follower ↔ followee) |
| `0008_news` | `news` table |

### Key Tables & Relationships

```
users (INT PK)
  ├── recipes (UUID PK, FK → users, FK → categories)
  │     ├── recipe_ingredients (UUID PK, FK → recipes)
  │     ├── recipe_steps       (UUID PK, FK → recipes)
  │     ├── recipe_media       (UUID PK, FK → recipes)
  │     ├── recipe_ratings     (SERIAL PK, unique per user+recipe, triggers avg update)
  │     └── favorites          (FK → users + recipes)
  ├── follows (follower_id FK → users, followee_id FK → users)
  └── news (UUID PK, FK → users)

categories (INT PK)
  └── recipes FK → categories
```

---

## 🔧 Backend — Go / Gin (`food-recipes-backend`)

### Startup Flow (`main.go`)
1. `config.LoadEnv()` → loads `.env`
2. `db.RunMigrations(dsn)` → applies SQL migrations (golang-migrate)
3. `gql.InitHasuraClient(endpoint, secret)` → creates Hasura admin client
4. Sets up Gin router with CORS, registers all routes
5. Listens on `:8081`

### Package Map

```
food-recipes-backend/
├── main.go            ← Entry point, router setup
├── config/env.go      ← LoadEnv() + GetEnv() helpers
├── db/migrate.go      ← Runs SQL migrations via golang-migrate
├── gql/
│   └── hasura_client.go ← Global *graphql.Client with admin-secret transport
├── middleware/
│   └── auth.go        ← JWT parsing; Authenticated() + RequireRole() guards
├── utils/
│   ├── jwt.go         ← GenerateJWT(userID, role) → Hasura-compatible JWT
│   └── hash.go        ← bcrypt helpers
├── handlers/
│   ├── auth.go        ← Register, Login, ChangePassword
│   ├── users.go       ← PublicGetUser, PublicListUserRecipes, FollowUser, UnfollowUser, GetFollowStatus
│   ├── admin.go       ← PublicListChefs, AdminListChefs, AdminVerifyChef, AdminResetPassword
│   ├── categories.go  ← ListCategories (seeds defaults if empty)
│   ├── news.go        ← PublicListNews, CreateNews
│   ├── upload.go      ← UploadImage (→ Cloudinary, up to 3 files)
│   └── payment.go     ← (stub / future)
└── migrations/        ← 12 SQL files (up + down)
```

### API Routes

| Method | Path | Auth | Handler |
|--------|------|------|---------|
| GET | `/health` | Public | health check |
| POST | `/register` | Public | Register new user |
| POST | `/login` | Public | Login → JWT |
| POST | `/upload` | Public | Upload 1–3 images/videos to Cloudinary |
| GET | `/categories` | Public | List categories (seeds if empty) |
| GET | `/chefs` | Public | List all chefs (optional `?verified=true/false`) |
| GET | `/news` | Public | List all news articles |
| GET | `/users/:id` | Public | Public user profile |
| GET | `/users/:id/recipes` | Public | User's recipes |
| POST | `/news` | 🔐 JWT | Create news article |
| POST | `/password/change` | 🔐 JWT | Change own password |
| GET | `/users/:id/following` | 🔐 JWT | Check follow status |
| POST | `/users/:id/follow` | 🔐 JWT | Follow a user |
| DELETE | `/users/:id/follow` | 🔐 JWT | Unfollow a user |
| GET | `/admin/chefs` | 🔐 admin | Admin list chefs |
| PATCH | `/admin/chefs/:id/verify` | 🔐 admin | Verify/revoke chef |
| PATCH | `/admin/users/:id/password` | 🔐 admin | Reset any user password |

### JWT Structure

```json
{
  "sub": "userId",
  "exp": "now + 72h",
  "https://hasura.io/jwt/claims": {
    "x-hasura-user-id": "userId",
    "x-hasura-allowed-roles": ["public", "guest", "chef", "admin"],
    "x-hasura-default-role": "guest",
    "x-hasura-role": "guest"
  }
}
```
> Signed with HS256 using `JWT_SECRET`. Hasura reads the same secret to validate requests from the frontend.

### User Roles

| Role | Permissions |
|------|------------|
| `guest` | Recipes, favorites, rating, follow, change own password |
| `chef` | Same as guest + `is_verified` flag (set by admin) |
| `admin` | All + admin panel (chef verify, reset any password) |

---

## 🌐 Frontend — Nuxt 3 SPA (`food-recipes-frontend`)

### Configuration
- **SSR: disabled** (pure SPA, `ssr: false`)
- **Styling:** Tailwind CSS
- **GraphQL client:** Apollo Client (via `@vue/apollo-composable`)
- **Auth state:** Stored in `localStorage` (`token`, `userName`, `userId`, `userEmail`, `role`, `is_verified`)

### Plugin — `apollo.js`
- Creates an `ApolloClient` pointed at Hasura (`port 8082`)
- `authLink` automatically injects `Authorization: Bearer <token>` from localStorage
- Provided globally via `DefaultApolloClient`

### Middleware — `auth.global.js` (runs on every route)
- **Public routes:** `/`, `/login`, `/signup`, `/recipes`, `/chefs`, `/news`, `/shows`, `/sweepstakes`
- **Protected routes:** `/profile`, `/favorites`, `/recipes/create`, `/recipes/my`, `/recipes/edit`, `/admin/*`
- Unauthenticated access to protected routes → redirect to `/login`
- Non-admin access to `/admin/*` → redirect to `/`

### Pages Map

```
pages/
├── index.vue          ← Home: hero + search + recipes grid + top creators
├── login.vue          ← Login form → POST /login → stores JWT in localStorage
├── signup.vue         ← Signup form → POST /register
├── profile.vue        ← My Account: info, stats, recent recipes, change password
├── favorites.vue      ← My favorites (GraphQL query by userId)
├── chefs.vue          ← Browse chefs (GET /chefs + search filter)
├── news.vue           ← Redirect wrapper to /news/index
├── top-fav.vue        ← Top favorited recipes
├── shop.vue           ← (stub/future)
├── shows.vue          ← (stub/future)
├── sweepstakes.vue    ← (stub/future)
├── news/
│   ├── index.vue      ← Browse news (GET /news + search)
│   └── create.vue     ← Create news article (POST /news)
├── recipes/
│   ├── index.vue      ← Browse all recipes (GraphQL)
│   ├── [id].vue       ← Recipe detail: gallery, steps, ingredients, rating, favorite
│   ├── create.vue     ← Create recipe form: media upload + GraphQL insert
│   ├── my.vue         ← My recipes list + delete
│   └── edit/[id].vue  ← Edit recipe
├── users/
│   └── [id].vue       ← Public user profile + their recipes + follow/unfollow
└── admin/
    └── chefs.vue      ← Admin: approve/revoke chef verification
```

### Components

| Component | Purpose |
|-----------|---------|
| `Navbar.vue` | Navigation bar with auth-aware dropdown (Login/SignUp or user initials + admin link) |
| `RecipeCard.vue` | Recipe card with favorite toggle (used across home, favorites, user profile) |
| `CategoryGrid.vue` | Clickable category filter pills |
| `CreatorCard.vue` | Chef/creator card (used on home + chefs page) |
| `Modal.vue` | Reusable confirm dialog (used for recipe delete) |
| `Toast.vue` | Notification toast |

---

## 🔄 Complete Data Flow Diagrams

### 1. User Registration / Login

```
Browser signup.vue
  → POST /register (Go backend)
    → Hash password (bcrypt)
    → Check email duplicate (Hasura GraphQL)
    → Insert user (Hasura GraphQL)
    → Return { id }
  → Redirect to /login

Browser login.vue
  → POST /login (Go backend)
    → Query user by email (Hasura GraphQL)
    → bcrypt.CompareHashAndPassword
    → GenerateJWT(userId, role)
    → Return { token, name, email, id, role }
  → Store in localStorage
  → window.location.href = '/'
```

### 2. Recipe Creation

```
Browser recipes/create.vue
  1. User picks files → onMediaSelected()
       → POST /upload (Go backend)
         → Cloudinary upload (up to 3 files)
         → Return { url, type } for each
       → Preview displayed, URL stored in mediaItems[]

  2. User submits form → submitRecipe()
       → Apollo useMutation: insert_recipes_one (Hasura via JWT)
         → Returns recipe UUID
       → Apollo useMutation: insert_recipe_ingredients
       → Apollo useMutation: insert_recipe_steps
       → Apollo useMutation: insert_recipe_media
       → Navigate to /recipes/my
```

### 3. Recipe Rating (with DB Trigger)

```
Browser recipes/[id].vue  
  → User clicks star (1–5)
  → Apollo useMutation: insert_recipe_ratings_one
    → If conflict (already rated): useMutation update_recipe_ratings
  → PostgreSQL trigger: trg_update_avg_rating_insupd
    → Recalculates AVG(rating) for that recipe
    → Updates recipes.average_rating in place
  → Refetch recipe to show new rating
```

### 4. Follow / Unfollow

```
Browser users/[id].vue
  → GET /users/:id/following (Go backend, JWT)  ← check status
  → Click Follow/Unfollow
  → POST/DELETE /users/:id/follow (Go backend, JWT)
    → Hasura mutation: insert_follows / delete_follows
  → isFollowing toggled locally
```

### 5. Admin Chef Verification

```
Browser admin/chefs.vue
  → GET /admin/chefs (Go backend, RequireRole("admin"))
  → Admin clicks Approve/Revoke
  → PATCH /admin/chefs/:id/verify (Go backend, RequireRole("admin"))
    → Hasura mutation: update_users_by_pk → is_verified = true/false
  → Refetch chef list
```

---

## ⚙️ Infrastructure (Docker Compose)

```yaml
Services:
  postgres   → port 5433:5432, volume db_data
  hasura     → port 8082:8080, depends_on postgres (healthy)
  backend    → port 8081:8081, depends_on hasura
  frontend   → port 3000:3000, depends_on backend
```

**Startup order:** `postgres` → `hasura` → `backend` (runs migrations) → `frontend`

---

## 🔑 Environment Variables

| Variable | Used By | Value |
|----------|---------|-------|
| `HASURA_GRAPHQL_ENDPOINT` | Backend | `http://hasura:8080/v1/graphql` |
| `HASURA_GRAPHQL_ADMIN_SECRET` | Backend + Hasura | `myhasurasecret` |
| `JWT_SECRET` | Backend (sign) + Hasura (verify) | `s9d8f7g6h5j4k3l2q1w0e9r8t7y6u5i4o3p2a1s0` |
| `POSTGRES_DSN` | Backend | `postgres://postgres:mysecretpassword@postgres:5432/foodrecipes` |
| `CLOUDINARY_URL` | Backend | from `.env` |
| `CORS_ORIGINS` | Backend | `http://localhost:3000,http://127.0.0.1:3000` |
| `NUXT_PUBLIC_HASURA_ENDPOINT` | Frontend | `http://localhost:8082/v1/graphql` |
| `NUXT_PUBLIC_BACKEND_ENDPOINT` | Frontend | `http://localhost:8081` |

---

## 📝 Notable Observations

> [!NOTE]
> The frontend queries Hasura **directly** via Apollo for recipes, favorites, ratings, and categories. It only calls the Go backend for auth, uploads, chefs listing, news, user profiles, and admin operations.

> [!TIP]
> The `average_rating` is maintained by a **PostgreSQL trigger** — no application-layer calculation needed. It auto-updates whenever a rating is inserted, updated, or deleted.

> [!WARNING]
> Auth state (`token`, `role`, `userId`) lives **entirely in localStorage**. There's no server-side session — security relies 100% on JWT expiry (72 hours).

> [!NOTE]
> The `payment.go` handler exists but appears to be a stub for future payment/subscription functionality.

> [!WARNING]
> `app.vue` has a hardcoded `<NuxtLink to="/profile">Go to Profile</NuxtLink>` that appears on every page. This is likely unintentional and should be removed or styled as a proper layout element.
