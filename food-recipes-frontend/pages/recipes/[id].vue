<template>
  <div class="min-h-screen bg-gradient-to-br from-orange-50 to-white">
    <Navbar />
    <section class="container mx-auto px-4 pt-6 max-w-5xl">
      <div v-if="loading" class="flex justify-center items-center h-64">
        <span class="loader"></span>
      </div>
      <div v-else-if="error" class="text-red-500 text-center">{{ error.message }}</div>
      <div v-else class="space-y-6">
        <div class="rounded-xl overflow-hidden shadow">
          <img :src="recipe.featured_image || 'https://placehold.co/1200x480?text=Recipe'" class="w-full h-64 md:h-80 object-cover" />
        </div>

        <div class="flex flex-col md:flex-row md:items-start md:justify-between gap-4">
          <div>
            <h1 class="text-3xl md:text-4xl font-bold">{{ recipe.title }}</h1>
            <p class="mt-2 text-gray-600 max-w-3xl">{{ recipe.description }}</p>
            <div class="mt-3 flex items-center gap-2">
              <span class="text-yellow-500">★</span>
              <span class="font-semibold">{{ Number(recipe.average_rating ?? 0).toFixed(1) }}</span>
              <span class="text-sm text-gray-500">({{ ratingCount }} ratings)</span>
            </div>
            <div class="mt-2 flex items-center gap-3">
              <button
                type="button"
                class="px-3 py-1 rounded border border-gray-300 text-sm hover:bg-gray-50"
                @click="onFavoriteClick"
              >
                <span v-if="isFavorited">❤ Unfavorite</span>
                <span v-else>♡ Favorite</span>
              </button>
              <span class="text-sm text-gray-500">❤ {{ favoriteCount }}</span>
            </div>
          </div>
          <div class="shrink-0 bg-white border border-gray-200 rounded-lg p-4 shadow-sm min-w-[220px]">
            <div class="text-sm text-gray-500">Category</div>
            <div class="font-semibold">{{ categoryName }}</div>
            <div class="mt-3 text-sm text-gray-500">Prep Time</div>
            <div class="font-semibold">{{ recipe.prep_time_minutes ? recipe.prep_time_minutes + ' mins' : '—' }}</div>
          </div>
        </div>

        <div v-if="isLoggedIn" class="bg-white border border-gray-200 rounded-xl p-4 shadow-sm">
          <div class="flex items-center gap-3">
            <div class="text-sm text-gray-600">Your rating:</div>
            <div class="flex items-center gap-1">
              <button
                v-for="n in 5"
                :key="n"
                type="button"
                class="text-2xl focus:outline-none"
                :class="n <= currentRating ? 'text-yellow-500' : 'text-gray-300'"
                @click="rate(n)"
                :disabled="rateLoading"
                aria-label="Rate"
              >★</button>
            </div>
            <div v-if="rateLoading" class="text-xs text-gray-400">Saving...</div>
          </div>
        </div>

        <div v-if="(recipe.media?.length || 0) > 0" class="bg-white border border-gray-200 rounded-xl p-4 shadow-sm">
          <h2 class="text-xl font-semibold mb-3">Gallery</h2>
          <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-3">
            <div v-for="m in galleryMedia" :key="m.id" class="aspect-video bg-gray-100 rounded overflow-hidden flex items-center justify-center">
              <template v-if="m.media_type === 'video'">
                <video :src="displaySrc(m.url, 'video')" class="w-full h-full object-cover" controls playsinline preload="metadata"></video>
              </template>
              <template v-else>
                <img :src="displaySrc(m.url, 'image')" class="w-full h-full object-cover" @error="imgError" />
              </template>
            </div>
          </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div class="bg-white border border-gray-200 rounded-xl p-5 shadow-sm">
            <h2 class="text-xl font-semibold mb-3">Ingredients</h2>
            <ul class="space-y-2">
              <li v-for="ing in recipe.ingredients" :key="ing.id" class="flex items-center gap-2">
                <span class="w-1.5 h-1.5 bg-orange-500 rounded-full"></span>
                <span>{{ ing.name }}<span v-if="ing.quantity"> — {{ ing.quantity }}</span></span>
              </li>
            </ul>
          </div>
          <div class="bg-white border border-gray-200 rounded-xl p-5 shadow-sm">
            <h2 class="text-xl font-semibold mb-3">Steps</h2>
            <ol class="space-y-3">
              <li v-for="(step, i) in recipe.steps" :key="step.id" class="flex items-start gap-3">
                <div class="w-7 h-7 rounded-full bg-orange-600 text-white flex items-center justify-center text-sm font-semibold mt-0.5">{{ i+1 }}</div>
                <div class="pt-0.5">{{ step.description }}</div>
              </li>
            </ol>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import Navbar from '~/components/Navbar.vue'
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useQuery, useMutation } from '@vue/apollo-composable'
import gql from 'graphql-tag'

function displaySrc(url, kind) {
  // Use the original Cloudinary URL without adding transformations to avoid 400s when strict transformations are enabled
  return typeof url === 'string' ? url : ''
}

function imgError(e) {
  // replace failed images with a subtle placeholder
  if (e && e.target) e.target.src = 'https://placehold.co/600x400?text=Image+Unavailable'
}

function onFavoriteClick() {
  if (!isLoggedIn.value) {
    router.push('/login')
    return
  }
  toggleFavorite()
}
const route = useRoute()
const router = useRouter()
const isLoggedIn = ref(false)
const userId = ref(null)

onMounted(() => {
  if (typeof window !== 'undefined') {
    isLoggedIn.value = !!localStorage.getItem('token')
    const raw = localStorage.getItem('userId') || ''
    const parsed = parseInt(raw, 10)
    userId.value = Number.isFinite(parsed) ? parsed : null
  }
})

const { result, loading, error, refetch } = useQuery(gql`
  query($id: uuid!, $uid: Int = 0) {
    recipes_by_pk(id: $id) {
      id
      title
      description
      featured_image
      prep_time_minutes
      category_id
      ingredients: recipe_ingredients { id name quantity }
      steps(order_by: {position: asc}) { id description position }
      media(order_by: {position: asc}) { id url media_type position }
      average_rating
      recipe_ratings_aggregate { aggregate { count } }
      my_rating: recipe_ratings(where: { user_id: { _eq: $uid } }, limit: 1) { rating }
      favorites_aggregate { aggregate { count } }
      my_favorite: favorites(where: { user_id: { _eq: $uid } }, limit: 1) { id }
    }
  }
`, () => ({ id: route.params.id, uid: userId.value ?? 0 }))

const { result: catResult } = useQuery(gql`query { categories { id name } }`)

const recipe = computed(() => result.value?.recipes_by_pk || {})
const categoryName = computed(() => {
  const cid = recipe.value?.category_id
  if (!cid) return 'Uncategorized'
  const cats = catResult.value?.categories || []
  const found = cats.find(c => c?.id === cid)
  return found?.name || 'Uncategorized'
})
const galleryMedia = computed(() => (recipe.value?.media || []).filter(m => !!m?.url))
const ratingCount = computed(() => recipe.value?.recipe_ratings_aggregate?.aggregate?.count ?? 0)
const currentRating = computed(() => {
  const r = recipe.value?.my_rating?.[0]?.rating
  return typeof r === 'number' ? r : 0
})
const favoriteCount = computed(() => recipe.value?.favorites_aggregate?.aggregate?.count ?? 0)
const isFavorited = computed(() => (recipe.value?.my_favorite?.length ?? 0) > 0)

const rateLoading = ref(false)

async function rate(n) {
  if (!isLoggedIn.value || !userId.value) return
  rateLoading.value = true
  try {
    const token = localStorage.getItem('token') || ''
    const backendBaseUrl = useRuntimeConfig().public.NUXT_PUBLIC_BACKEND_ENDPOINT || 'http://localhost:8081'
    const res = await fetch(`${backendBaseUrl}/recipes/${route.params.id}/rate`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...(token ? { Authorization: `Bearer ${token}` } : {})
      },
      body: JSON.stringify({ rating: n })
    })
    if (!res.ok) {
      const data = await res.json().catch(() => ({}))
      throw new Error(data?.error || 'Rating request failed')
    }
  } catch (e) {
    console.error('rating failed', e)
  }
  rateLoading.value = false
  try { await refetch() } catch {}
}

const backendBase = useRuntimeConfig().public.NUXT_PUBLIC_BACKEND_ENDPOINT || 'http://localhost:8081'
const favToggling = ref(false)

async function toggleFavorite() {
  if (!isLoggedIn.value || !userId.value) return
  if (favToggling.value) return
  favToggling.value = true
  try {
    const token = localStorage.getItem('token') || ''
    const res = await fetch(`${backendBase}/recipes/${route.params.id}/favorite`, {
      method: isFavorited.value ? 'DELETE' : 'POST',
      headers: token ? { Authorization: `Bearer ${token}` } : {}
    })
    if (!res.ok) {
      const data = await res.json().catch(() => ({}))
      throw new Error(data?.error || 'Request failed')
    }
    try { await refetch() } catch {}
  } catch (e) {
    console.error('favorite failed', e)
  } finally {
    favToggling.value = false
  }
}
</script>

<style scoped>
.loader {
  border: 4px solid #f3f3f3;
  border-top: 4px solid #3490dc;
  border-radius: 50%;
  width: 40px;
  height: 40px;
  animation: spin 1s linear infinite;
}
@keyframes spin {
  0% { transform: rotate(0deg);}
  100% { transform: rotate(360deg);}
}
</style>
