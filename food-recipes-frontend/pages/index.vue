<template>
  <div class="min-h-screen bg-gradient-to-br from-orange-50 to-white">
    <Navbar />

    <!-- Hero / Search -->
    <section class="container mx-auto px-4 pt-6 max-w-6xl">
      <div class="rounded-2xl bg-gradient-to-r from-orange-100 via-amber-100 to-yellow-100 p-8 shadow-sm">
        <h1 class="text-4xl md:text-5xl font-extrabold tracking-tight text-gray-900">Cook, Share, and Discover Amazing Recipes</h1>
        <p class="mt-3 text-gray-600">Browse recipes by category, creator, ingredients, or prep time.</p>
        <div class="mt-6 flex gap-3">
          <input v-model="search" type="text" placeholder="Search recipes by title..." class="flex-1 rounded-xl border border-gray-200 bg-white px-4 py-3 shadow-sm focus:outline-none focus:ring-2 focus:ring-orange-300" />
          <button @click="clearSearch" class="rounded-xl px-4 py-3 bg-gray-100 hover:bg-gray-200 text-gray-700">Clear</button>
        </div>
        <div class="mt-4">
          <CategoryGrid :categories="categories" @select="selectCategory" />
        </div>
      </div>
    </section>

    <!-- Featured Recipes -->
    <section class="container mx-auto px-4 py-10 max-w-6xl">
      <div class="flex items-center justify-between mb-6">
        <h2 class="text-2xl md:text-3xl font-bold">Featured Recipes</h2>
        <NuxtLink to="/recipes" class="text-orange-600 hover:underline">Browse all</NuxtLink>
      </div>
      <div v-if="recLoading" class="text-center text-gray-500 py-10">Loading...</div>
      <div v-else-if="recError" class="text-center text-red-600 py-10">{{ recError.message }}</div>
      <div v-else>
        <div v-if="filteredRecipes.length === 0" class="text-center text-gray-400 py-10">No recipes found.</div>
        <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-8">
          <RecipeCard
            v-for="recipe in filteredRecipes"
            :key="recipe.id"
            :recipe="recipe"
            :favorites="favorites"
            @refresh-favorites="refetchFavorites"
          />
        </div>
      </div>
    </section>

    <!-- Top Creators -->
    <section class="container mx-auto px-4 pb-16 max-w-6xl">
      <h2 class="text-2xl md:text-3xl font-bold mb-6">Top Creators</h2>
      <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-6">
        <CreatorCard v-for="u in creators" :key="u.id" :creator="u" />
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuery } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import Navbar from '~/components/Navbar.vue'
import RecipeCard from '~/components/RecipeCard.vue'
import CategoryGrid from '~/components/CategoryGrid.vue'
import CreatorCard from '~/components/CreatorCard.vue'

// Queries (fallback to sample data if unavailable)
// Use newest first and a network-aware fetch policy so new items appear promptly
const { result: recResult, loading: recLoading, error: recError } = useQuery(
  gql`query { recipes(order_by: { created_at: desc }) { id title description featured_image average_rating category_id category { id name slug } media(order_by: {position: asc}) { url media_type position } favorites_aggregate { aggregate { count } } } }`,
  null,
  { fetchPolicy: 'cache-and-network' }
)
const { result: catResult, error: catError } = useQuery(gql`query { categories { id name slug } }`)
const backend = useRuntimeConfig().public.NUXT_PUBLIC_BACKEND_ENDPOINT
const creatorsRaw = ref([])

async function fetchCreators() {
  try {
    const res = await fetch(`${backend}/chefs?verified=true`)
    const data = await res.json()
    if (!res.ok) return
    creatorsRaw.value = data?.chefs || []
  } catch {
    // ignore
  }
}

onMounted(fetchCreators)

const search = ref('')
const selectedCategory = ref(null)

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
const creators = computed(() => creatorsRaw.value || [])

const filteredRecipes = computed(() => {
  const s = search.value.trim().toLowerCase()
  return recipes.value.filter(r => {
    const matchesSearch = !s || r.title?.toLowerCase().includes(s)
    const matchesCategory = !selectedCategory.value || r.category?.slug === selectedCategory.value
    return matchesSearch && matchesCategory
  })
})

function selectCategory(slug) { selectedCategory.value = slug }
function clearSearch() { search.value = '' }

// Favorites for current viewer
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

