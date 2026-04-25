<template>
  <Toast ref="toastRef" />
  <div class="container mx-auto p-4">
    <h1 class="text-3xl font-bold mb-4">My Favorites</h1>
    <div v-if="!isLoggedIn" class="text-center text-gray-500">Please log in to see your favorites.</div>
    <div v-else>
      <div v-if="loading" class="text-center">Loading...</div>
      <div v-else-if="favorites.length === 0" class="text-center text-gray-400">No favorites yet.</div>
      <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-8">
        <RecipeCard
          v-for="fav in favorites"
          :key="fav.recipe.id"
          :recipe="fav.recipe"
          :favorites="favorites"
          :toastRef="toastRef"
          @refresh-favorites="refetch()"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuery } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import RecipeCard from '~/components/RecipeCard.vue'
import Toast from '~/components/Toast.vue'

const toastRef = ref(null)
const isLoggedIn = computed(() => !!localStorage.getItem('token'))
const userId = parseInt(localStorage.getItem('userId') || '', 10)
const { result, loading, refetch } = useQuery(gql`
  query($userId: Int!) {
    favorites(where: {user_id: {_eq: $userId}}) {
      recipe {
        id
        title
        description
        featured_image
        average_rating
        favorites_aggregate { aggregate { count } }
      }
    }
  }
`, { userId })

const favorites = computed(() => result.value?.favorites || [])

onMounted(() => {
  if (!isLoggedIn.value) window.location.href = '/login'
})
</script>