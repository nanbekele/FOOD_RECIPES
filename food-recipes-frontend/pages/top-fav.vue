<template>
  <div class="min-h-screen bg-gradient-to-br from-orange-50 to-white">
    <Navbar />
    <section class="container mx-auto px-4 pt-6 max-w-6xl">
      <!-- Header -->
      <div class="flex items-center justify-between mb-6">
        <div>
          <h1 class="text-3xl md:text-4xl font-bold">🏆 Top Favorites</h1>
          <p class="text-gray-500 text-sm mt-1">Recipes loved most by the community</p>
        </div>
        <NuxtLink to="/recipes" class="text-orange-600 hover:underline text-sm font-medium">Browse all →</NuxtLink>
      </div>

      <!-- Loading skeletons -->
      <div v-if="loading" class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-6">
        <div v-for="n in 6" :key="n" class="bg-white rounded-xl shadow-lg h-64 animate-pulse" />
      </div>

      <!-- Error -->
      <div v-else-if="error" class="text-center text-red-600 py-10 bg-white rounded-xl shadow p-6">
        <p class="text-lg font-semibold">Failed to load top favorites</p>
        <p class="text-sm text-gray-500 mt-1">{{ error.message }}</p>
      </div>

      <!-- Empty state -->
      <div v-else-if="recipes.length === 0" class="text-center py-20 text-gray-400">
        <div class="text-6xl mb-4">💔</div>
        <p class="text-xl font-semibold">No favorited recipes yet</p>
        <p class="text-sm mt-2">Be the first to favorite a recipe!</p>
        <NuxtLink to="/recipes" class="mt-4 inline-block px-4 py-2 bg-orange-600 text-white rounded-lg hover:bg-orange-700">
          Browse Recipes
        </NuxtLink>
      </div>

      <!-- Recipe grid with rank badges -->
      <div v-else class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-6">
        <div v-for="(r, index) in recipes" :key="r.id" class="relative">
          <!-- Rank badge -->
          <div
            class="absolute -top-3 -left-3 z-20 w-9 h-9 rounded-full flex items-center justify-center text-sm font-bold shadow-md border-2 border-white"
            :class="{
              'bg-yellow-400 text-yellow-900': index === 0,
              'bg-gray-300 text-gray-700':     index === 1,
              'bg-amber-600 text-white':        index === 2,
              'bg-orange-100 text-orange-700':  index >= 3
            }"
          >
            {{ index === 0 ? '🥇' : index === 1 ? '🥈' : index === 2 ? '🥉' : `#${index + 1}` }}
          </div>

          <!-- Fav count pill -->
          <div class="absolute -top-3 right-2 z-20 bg-red-500 text-white text-xs font-semibold px-2 py-0.5 rounded-full shadow">
            ❤ {{ r.favorites_aggregate?.aggregate?.count ?? 0 }}
          </div>

          <RecipeCard
            :recipe="r"
            :favorites="favorites"
            @refresh-favorites="refetchFavorites"
          />
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import Navbar from '~/components/Navbar.vue'
import RecipeCard from '~/components/RecipeCard.vue'
import { computed } from 'vue'
import { useQuery } from '@vue/apollo-composable'
import gql from 'graphql-tag'

// Fetch all recipes sorted by most favorited — filter out zero-fav in JS
// (avoids Hasura aggregate-where compatibility issues)
const { result, loading, error } = useQuery(gql`
  query TopFav {
    recipes(
      order_by: [
        { average_rating: desc }
        { created_at: desc }
      ]
      limit: 50
    ) {
      id
      title
      description
      featured_image
      average_rating
      media(order_by: { position: asc }) { url media_type position }
      favorites_aggregate { aggregate { count } }
    }
  }
`, null, { fetchPolicy: 'cache-and-network' })

// Only show recipes that have at least 1 favorite, sorted by count desc
const recipes = computed(() => {
  const all = result.value?.recipes || []
  return all
    .filter(r => (r.favorites_aggregate?.aggregate?.count ?? 0) > 0)
    .sort((a, b) =>
      (b.favorites_aggregate?.aggregate?.count ?? 0) -
      (a.favorites_aggregate?.aggregate?.count ?? 0)
    )
    .slice(0, 12)
})

// Viewer's own favorites so heart icon reflects their state
const userId = computed(() => {
  if (typeof window === 'undefined') return 0
  const raw = localStorage.getItem('userId') || ''
  const parsed = parseInt(raw, 10)
  return Number.isFinite(parsed) ? parsed : 0
})

const { result: favResult, refetch: refetchFavorites } = useQuery(gql`
  query MyFavs($uid: Int!) {
    favorites(where: { user_id: { _eq: $uid } }) { recipe { id } }
  }
`, () => ({ uid: userId.value }), { fetchPolicy: 'cache-and-network' })

const favorites = computed(() => favResult.value?.favorites || [])
</script>
