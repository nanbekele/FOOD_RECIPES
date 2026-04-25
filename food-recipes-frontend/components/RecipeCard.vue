<template>
  <div class="bg-white rounded-xl shadow-lg p-4 flex flex-col relative">
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
    <button
      type="button"
      @click.stop.prevent="onFavoriteClick"
      @keydown.enter.prevent="onFavoriteClick"
      @keydown.space.prevent="onFavoriteClick"
      class="absolute top-2 right-2 z-[999] w-10 h-10 flex items-center justify-center rounded-full bg-white/90 hover:bg-white border border-gray-200 text-yellow-500 shadow cursor-pointer pointer-events-auto"
      :title="isFavorite ? 'Remove from favorites' : 'Add to favorites'"
      aria-label="Toggle favorite"
      tabindex="0"
    >
      <span class="text-lg leading-none" v-if="isFavorite">★</span>
      <span class="text-lg leading-none" v-else>☆</span>
    </button>
  </div>
</template>
<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useMutation } from '@vue/apollo-composable'
import gql from 'graphql-tag'

const emit = defineEmits(['refresh-favorites'])
const props = defineProps({
  recipe: { type: Object, required: true },
  favorites: { type: Array, default: () => [] }, // Pass user's favorites from parent
  toastRef: { type: Object, default: null }
})

const isLoggedIn = ref(false)
const userId = ref(null)
const router = useRouter()

onMounted(() => {
  isLoggedIn.value = !!localStorage.getItem('token')
  userId.value = parseInt(localStorage.getItem('userId') || '', 10)
})

const isFavorite = computed(() =>
  props.favorites.some(fav => fav.recipe.id === props.recipe.id)
)

const favCount = ref(0)
onMounted(() => {
  const initial = props.recipe?.favorites_aggregate?.aggregate?.count
  favCount.value = typeof initial === 'number' ? initial : 0
})

const mediaUrl = computed(() => {
  if (props.recipe?.featured_image) return props.recipe.featured_image
  const list = Array.isArray(props.recipe?.media)
    ? props.recipe.media
    : (Array.isArray(props.recipe?.recipe_media) ? props.recipe.recipe_media : [])
  const first = list[0] || null
  return first?.url || ''
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
function onImgError(e){ if (e && e.target) e.target.src = 'https://placehold.co/400x200?text=Image+Unavailable' }
function onMediaError(){ }

// GraphQL mutations (user_id is set by Hasura using X-Hasura-User-Id)
const ADD_FAVORITE = gql`
  mutation($recipeId: uuid!) {
    insert_favorites_one(object: { recipe_id: $recipeId }) { id }
  }
`
const REMOVE_FAVORITE = gql`
  mutation($recipeId: uuid!) {
    delete_favorites(where: { recipe_id: { _eq: $recipeId } }) { affected_rows }
  }
`

const { mutate: addFavorite } = useMutation(ADD_FAVORITE)
const { mutate: removeFavorite } = useMutation(REMOVE_FAVORITE)

function toggleFavorite() {
  if (!isLoggedIn.value) return
  if (isFavorite.value) {
    removeFavorite({ recipeId: props.recipe.id }, {
      onDone: () => {
        props.toastRef?.value?.showToast?.('Removed from favorites', 'success')
        emit('refresh-favorites')
        if (favCount.value > 0) favCount.value -= 1
        // Optionally emit event to parent to refresh favorites
      }
    })
  } else {
    addFavorite({ recipeId: props.recipe.id }, {
      onDone: () => {
        props.toastRef?.value?.showToast?.('Added to favorites', 'success')
        emit('refresh-favorites')
        favCount.value += 1
        // Optionally emit event to parent to refresh favorites
      }
    })
  }
}

function onFavoriteClick() {
  if (!isLoggedIn.value) {
    if (typeof window !== 'undefined') {
      window.location.href = '/login'
    } else {
      router.push('/login')
    }
    return
  }
  toggleFavorite()
}
</script>
