<template>
  <div class="min-h-screen bg-gradient-to-br from-orange-50 to-white">
    <Navbar />
    <section class="container mx-auto px-4 pt-6 max-w-6xl">
      <h1 class="text-3xl md:text-4xl font-bold mb-4">Browse Recipes</h1>
      <div class="rounded-xl bg-white border border-gray-200 p-4 mb-6 shadow-sm">
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <input v-model="search" type="text" placeholder="Search by title" class="rounded-lg border border-gray-200 px-3 py-2" />
          <select v-model="prepFilter" class="rounded-lg border border-gray-200 px-3 py-2">
            <option value="">Prep time</option>
            <option value="15">≤ 15 min</option>
            <option value="30">≤ 30 min</option>
            <option value="60">≤ 60 min</option>
          </select>
          <input v-model="ingredient" type="text" placeholder="Filter by ingredient (client-side)" class="rounded-lg border border-gray-200 px-3 py-2" />
        </div>
        <div class="mt-4">
          <CategoryGrid :categories="categories" @select="selectCategory" />
        </div>
      </div>

      <div v-if="recLoading" class="text-center text-gray-500 py-10">Loading...</div>
      <div v-else-if="recError" class="text-center text-red-600 py-10">{{ recError.message }}</div>
      <div v-else>
        <div class="text-xs text-gray-400 pb-2">Raw: {{ recipes.length }} | Filtered: {{ filteredRecipes.length }}</div>
        <div v-if="filteredRecipes.length === 0" class="text-center text-gray-400 py-10">No recipes found.</div>
        <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-6">
          <RecipeCard
            v-for="r in filteredRecipes"
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
import { ref, computed, watchEffect } from 'vue'
import { useQuery } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import Navbar from '~/components/Navbar.vue'
import RecipeCard from '~/components/RecipeCard.vue'
import CategoryGrid from '~/components/CategoryGrid.vue'

const search = ref('')
const prepFilter = ref('')
const ingredient = ref('')
const selectedCategory = ref(null)

const { result: recResult, loading: recLoading, error: recError } = useQuery(
  gql`query { recipes(order_by: { created_at: desc }) { id title description featured_image prep_time_minutes average_rating media(order_by: {position: asc}) { url media_type position } favorites_aggregate { aggregate { count } } } }`,
  null,
  { fetchPolicy: 'cache-and-network' }
)
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
const { result: catResult } = useQuery(gql`query { categories { id name slug } }`)

const recipes = computed(() => recResult.value?.recipes || [])
const categories = computed(() => catResult.value?.categories || [
  { id: 1, name: 'Breakfast', slug: 'breakfast' },
  { id: 2, name: 'Lunch', slug: 'lunch' },
  { id: 3, name: 'Dinner', slug: 'dinner' },
  { id: 4, name: 'Dessert', slug: 'dessert' },
  { id: 5, name: 'Vegan', slug: 'vegan' },
  { id: 6, name: 'Vegetarian', slug: 'vegetarian' },
  { id: 7, name: 'Gluten Free', slug: 'gluten-free' },
  { id: 8, name: 'Keto', slug: 'keto' }
])

const favorites = computed(() => favResult.value?.favorites || [])
const filteredRecipes = computed(() => {
  const s = search.value.trim().toLowerCase()
  const p = prepFilter.value ? parseInt(prepFilter.value) : null
  const ing = ingredient.value.trim().toLowerCase()
  return recipes.value.filter(r => {
    const matchesSearch = !s || r.title?.toLowerCase().includes(s)
    const matchesCategory = true
    const matchesPrep = !p || (r.prep_time_minutes || 9999) <= p
    const matchesIngredient = !ing || r.description?.toLowerCase().includes(ing)
    return matchesSearch && matchesCategory && matchesPrep && matchesIngredient
  })
})

function selectCategory(slug) { selectedCategory.value = slug }

watchEffect(() => {
  if (typeof window === 'undefined') return
  const n = recResult.value?.recipes?.length
  console.log('[recipes/index] raw recipes length:', n)
  console.log('[recipes/index] filters:', {
    search: search.value,
    prepFilter: prepFilter.value,
    ingredient: ingredient.value,
    selectedCategory: selectedCategory.value,
  })
})
</script>