<template>
  <div class="min-h-screen bg-gradient-to-br from-orange-50 to-white">
    <Navbar />
    <section class="container mx-auto px-4 pt-6 max-w-6xl">
      <div class="flex items-center justify-between mb-6">
        <h1 class="text-3xl font-bold">Top Favorites</h1>
      </div>

      <div v-if="loading" class="text-center text-gray-500 py-10">Loading...</div>
      <div v-else-if="error" class="text-center text-red-600 py-10">{{ error.message }}</div>
      <div v-else>
        <div v-if="recipes.length === 0" class="text-center text-gray-400 py-10">No recipes found.</div>
        <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-6">
          <RecipeCard
            v-for="r in recipes"
            :key="r.id"
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

const { result, loading, error } = useQuery(gql`
  query TopFav {
    recipes(order_by: [{favorites_aggregate: {count: desc}}, {average_rating: desc}, {created_at: desc}], limit: 30) {
      id
      title
      description
      featured_image
      average_rating
      media(order_by: {position: asc}) { url media_type position }
      favorites_aggregate { aggregate { count } }
    }
  }
`)

const recipes = computed(() => result.value?.recipes || [])

const userId = computed(() => {
  if (typeof window === 'undefined') return 0
  const raw = localStorage.getItem('userId') || ''
  const parsed = parseInt(raw, 10)
  return Number.isFinite(parsed) ? parsed : 0
})
const { result: favResult, refetch: refetchFavorites } = useQuery(gql`
  query($uid: Int!) {
    favorites(where: { user_id: { _eq: $uid } }) { recipe { id } }
  }
`, () => ({ uid: userId.value }))
const favorites = computed(() => favResult.value?.favorites || [])
</script>
