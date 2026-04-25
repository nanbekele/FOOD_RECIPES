<template>
  <div class="min-h-screen bg-gradient-to-br from-orange-50 to-white">
    <Navbar />
    <section class="container mx-auto px-4 pt-6 max-w-6xl">
      <div class="flex items-center justify-between mb-6">
        <div class="flex items-center gap-3">
          <span class="inline-flex items-center justify-center h-14 w-14 rounded-full bg-orange-100 text-orange-700 font-bold text-xl">{{ initials }}</span>
          <div>
            <h1 class="text-2xl md:text-3xl font-bold flex items-center gap-2">
              <span>{{ user.name }}</span>
              <span v-if="user.role === 'chef' && user.is_verified" class="inline-flex items-center text-xs px-2 py-0.5 rounded bg-green-100 text-green-700">Verified</span>
            </h1>
            <p class="text-gray-600" v-if="user.bio">{{ user.bio }}</p>
          </div>
        </div>
        <button class="px-3 py-2 rounded bg-gray-100">Follow</button>
      </div>

      <div class="rounded-xl bg-white border border-gray-200 p-4 shadow-sm">
        <h2 class="text-xl font-semibold mb-4">Recipes by {{ user.name }}</h2>
        <div v-if="recipes.length === 0" class="text-center text-gray-400 py-10">No recipes yet.</div>
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
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useQuery } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import Navbar from '~/components/Navbar.vue'
import RecipeCard from '~/components/RecipeCard.vue'

const route = useRoute()
const userId = parseInt(route.params.id, 10)

const { result } = useQuery(gql`
  query($id: Int!) {
    users_by_pk(id: $id) { id name bio role is_verified }
    recipes(where: { user_id: { _eq: $id } }) {
      id
      title
      description
      featured_image
      average_rating
      favorites_aggregate { aggregate { count } }
    }
  }
`, { id: userId })

const user = computed(() => result.value?.users_by_pk || { id: userId, name: 'Chef', bio: '' })
const recipes = computed(() => result.value?.recipes || [])
const initials = computed(() => (user.value.name || '').split(' ').map(p => p[0]).join('').toUpperCase().slice(0,2))

// Viewer favorites for clickable state
const viewerId = computed(() => {
  if (typeof window === 'undefined') return 0
  const raw = localStorage.getItem('userId') || ''
  const parsed = parseInt(raw, 10)
  return Number.isFinite(parsed) ? parsed : 0
})
const { result: favResult, refetch: refetchFavorites } = useQuery(gql`
  query($uid: Int!) {
    favorites(where: { user_id: { _eq: $uid } }) { recipe { id } }
  }
`, () => ({ uid: viewerId.value }))
const favorites = computed(() => favResult.value?.favorites || [])
</script>