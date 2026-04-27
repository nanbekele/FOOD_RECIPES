<template>
  <div class="bg-white rounded-xl shadow-lg p-4 flex flex-col relative">
    <!-- Media -->
    <div class="w-full h-40 rounded mb-2 overflow-hidden">
      <template v-if="mediaUrl">
        <video v-if="isVideo"
               :src="mediaUrl"
               class="w-full h-40 object-cover"
               preload="metadata" playsinline muted controls
               @error="onMediaError" />
        <img v-else
             :src="mediaUrl"
             alt="Recipe image"
             class="w-full h-40 object-cover transition-transform duration-300 hover:scale-105 pointer-events-none z-0 select-none"
             @error="onImgError" />
      </template>
      <template v-else>
        <img src="https://placehold.co/400x200?text=No+Image" class="w-full h-40 object-cover" />
      </template>
    </div>

    <!-- Rating + fav count badges -->
    <div class="absolute top-2 left-2 flex items-center gap-1 z-10">
      <span class="bg-black/60 text-white text-xs px-2 py-1 rounded">★ {{ Number(recipe.average_rating ?? 0).toFixed(1) }}</span>
      <span class="bg-black/60 text-white text-xs px-2 py-1 rounded">❤ {{ favCount }}</span>
    </div>

    <h3 class="font-bold text-lg mb-1 truncate">{{ recipe.title }}</h3>
    <p class="text-gray-600 mb-2 line-clamp-2">{{ recipe.description }}</p>
    <NuxtLink
      :to="`/recipes/${recipe.id}`"
      class="mt-auto text-blue-600 hover:underline font-semibold"
    >View Recipe</NuxtLink>

    <!-- Favorite button -->
    <button
      type="button"
      @click.stop.prevent="onFavoriteClick"
      @keydown.enter.prevent="onFavoriteClick"
      @keydown.space.prevent="onFavoriteClick"
      :disabled="toggling"
      class="absolute top-2 right-2 z-[999] w-10 h-10 flex items-center justify-center rounded-full bg-white/95 hover:bg-white border shadow cursor-pointer pointer-events-auto transition-all duration-200"
      :class="isFavorite
        ? 'border-red-400 text-red-500 scale-110'
        : 'border-gray-200 text-gray-300 hover:text-red-400 hover:border-red-300'"
      :title="isFavorite ? 'Remove from favorites' : 'Add to favorites'"
      aria-label="Toggle favorite"
      tabindex="0"
    >
      <!-- Filled heart when favorited -->
      <svg v-if="isFavorite" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" class="w-5 h-5">
        <path d="M11.645 20.91l-.007-.003-.022-.012a15.247 15.247 0 01-.383-.218 25.18 25.18 0 01-4.244-3.17C4.688 15.36 2.25 12.174 2.25 8.25 2.25 5.322 4.714 3 7.688 3A5.5 5.5 0 0112 5.052 5.5 5.5 0 0116.313 3c2.973 0 5.437 2.322 5.437 5.25 0 3.925-2.438 7.111-4.739 9.256a25.175 25.175 0 01-4.244 3.17 15.247 15.247 0 01-.383.219l-.022.012-.007.004-.003.001a.752.752 0 01-.704 0l-.003-.001z" />
      </svg>
      <!-- Outline heart when not favorited -->
      <svg v-else xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" class="w-5 h-5">
        <path stroke-linecap="round" stroke-linejoin="round" d="M21 8.25c0-2.485-2.099-4.5-4.688-4.5-1.935 0-3.597 1.126-4.312 2.733-.715-1.607-2.377-2.733-4.313-2.733C5.1 3.75 3 5.765 3 8.25c0 7.22 9 12 9 12s9-4.78 9-12z" />
      </svg>
    </button>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'

const emit = defineEmits(['refresh-favorites'])
const props = defineProps({
  recipe: { type: Object, required: true },
  favorites: { type: Array, default: () => [] },
  toastRef: { type: Object, default: null }
})

const isLoggedIn = ref(false)
const router = useRouter()
const backendBase = useRuntimeConfig().public.NUXT_PUBLIC_BACKEND_ENDPOINT || 'http://localhost:8081'

onMounted(() => {
  isLoggedIn.value = !!localStorage.getItem('token')
})

// Derive favorited state from parent-supplied favorites list
const propIsFavorite = computed(() =>
  props.favorites.some(fav => fav?.recipe?.id === props.recipe.id)
)

const localIsFavorite = ref(false)
watch(propIsFavorite, (v) => { localIsFavorite.value = !!v }, { immediate: true })
const isFavorite = computed(() => localIsFavorite.value)

// Favorite count — keep in sync with recipe prop changes
const favCount = ref(0)
onMounted(() => {
  const initial = props.recipe?.favorites_aggregate?.aggregate?.count
  favCount.value = typeof initial === 'number' ? initial : 0
})
watch(() => props.recipe?.favorites_aggregate?.aggregate?.count, (v) => {
  if (typeof v === 'number') favCount.value = v
})

// Media helpers
const mediaUrl = computed(() => {
  if (props.recipe?.featured_image) return props.recipe.featured_image
  const list = Array.isArray(props.recipe?.media)
    ? props.recipe.media
    : (Array.isArray(props.recipe?.recipe_media) ? props.recipe.recipe_media : [])
  return list[0]?.url || ''
})
const isVideo = computed(() => {
  const list = Array.isArray(props.recipe?.media)
    ? props.recipe.media
    : (Array.isArray(props.recipe?.recipe_media) ? props.recipe.recipe_media : [])
  const first = list[0] || null
  if (first?.media_type) return first.media_type === 'video'
  const u = mediaUrl.value || ''
  return /\.(mp4|mov|avi|mkv|webm|m4v|3gp|mpeg|mpg)(\?.*)?$/i.test(u)
})
function onImgError(e) { if (e?.target) e.target.src = 'https://placehold.co/400x200?text=Image+Unavailable' }
function onMediaError() {}

const toggling = ref(false)

async function toggleFavorite() {
  // Read userId fresh from localStorage every time (avoids stale ref issues)
  const uid = parseInt(localStorage.getItem('userId') || '', 10)
  if (!uid || !Number.isFinite(uid)) return
  if (toggling.value) return
  toggling.value = true

  // Optimistic UI — flip immediately, rollback if error
  const wasF = localIsFavorite.value
  localIsFavorite.value = !wasF
  favCount.value = wasF ? Math.max(0, favCount.value - 1) : favCount.value + 1

  try {
    const token = localStorage.getItem('token') || ''
    const res = await fetch(`${backendBase}/recipes/${props.recipe.id}/favorite`, {
      method: wasF ? 'DELETE' : 'POST',
      headers: token ? { Authorization: `Bearer ${token}` } : {}
    })
    if (!res.ok) {
      const data = await res.json().catch(() => ({}))
      throw new Error(data?.error || 'Request failed')
    }
    emit('refresh-favorites')
  } catch (e) {
    // Rollback on error
    localIsFavorite.value = wasF
    favCount.value = wasF ? favCount.value + 1 : Math.max(0, favCount.value - 1)
    console.error('Favorite toggle failed', e)
  } finally {
    toggling.value = false
  }
}

function onFavoriteClick() {
  if (!isLoggedIn.value) {
    if (typeof window !== 'undefined') window.location.href = '/login'
    else router.push('/login')
    return
  }
  toggleFavorite()
}
</script>
