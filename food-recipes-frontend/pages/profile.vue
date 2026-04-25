<template>
  <div class="min-h-screen bg-gradient-to-br from-orange-50 to-white">
    <Navbar />
    <section class="container mx-auto px-4 pt-6 max-w-6xl">
      <h1 class="text-3xl font-bold mb-4">My Account</h1>
      <div v-if="!isLoggedIn" class="text-center text-gray-500">Please log in to see your profile.</div>
      <div v-else>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div class="md:col-span-2 bg-white p-6 rounded-lg shadow-md">
            <h2 class="text-xl font-semibold mb-4">Profile Information</h2>
            <p class="mb-2"><strong>Name:</strong> {{ userName || 'Not available' }}</p>
            <p class="mb-2"><strong>Email:</strong> {{ userEmail || 'Not available' }}</p>
            <p class="mb-2"><strong>User ID:</strong> {{ userId || 'Not available' }}</p>
            <div class="mt-6 flex gap-3">
              <NuxtLink to="/recipes/create" class="px-3 py-2 rounded bg-orange-600 text-white">Create Recipe</NuxtLink>
              <NuxtLink to="/recipes/my" class="px-3 py-2 rounded bg-gray-100">Manage My Recipes</NuxtLink>
              <NuxtLink to="/favorites" class="px-3 py-2 rounded bg-gray-100">Favorites</NuxtLink>
            </div>
          </div>
          <div class="bg-white p-6 rounded-lg shadow-md">
            <h2 class="text-xl font-semibold mb-4">Stats</h2>
            <p class="mb-2">Total Recipes: <strong>{{ recipeCount }}</strong></p>
            <p class="text-sm text-gray-500">These counts update as you add or delete recipes.</p>
          </div>
        </div>

        <div class="mt-8 bg-white p-6 rounded-lg shadow-md">
          <h2 class="text-xl font-semibold mb-4">Recent Recipes</h2>
          <div v-if="recentRecipes.length === 0" class="text-gray-400">No recipes yet.</div>
          <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-6">
            <div v-for="r in recentRecipes" :key="r.id" class="border rounded-lg p-3">
              <img :src="r.featured_image || 'https://placehold.co/400x200?text=No+Image'" class="w-full h-32 object-cover rounded" />
              <h3 class="mt-2 font-semibold">{{ r.title }}</h3>
              <div class="mt-2 flex gap-2">
                <NuxtLink :to="`/recipes/${r.id}`" class="px-3 py-1 rounded bg-gray-100">View</NuxtLink>
                <NuxtLink :to="`/recipes/edit/${r.id}`" class="px-3 py-1 rounded bg-blue-600 text-white">Edit</NuxtLink>
              </div>
            </div>
          </div>
        </div>

        <div class="mt-8 bg-white p-6 rounded-lg shadow-md">
          <h2 class="text-xl font-semibold mb-4">Change Password</h2>
          <form @submit.prevent="onChangePassword" class="grid grid-cols-1 md:grid-cols-3 gap-3 max-w-2xl">
            <input v-model="oldPassword" type="password" placeholder="Current password" class="rounded border border-gray-200 px-3 py-2" required />
            <input v-model="newPassword" type="password" placeholder="New password (min 6 chars)" class="rounded border border-gray-200 px-3 py-2" required />
            <button type="submit" class="px-3 py-2 bg-blue-600 text-white rounded" :disabled="pwdLoading">Change</button>
          </form>
          <div v-if="pwdSuccess" class="text-green-600 mt-2">{{ pwdSuccess }}</div>
          <div v-if="pwdError" class="text-red-600 mt-2">{{ pwdError }}</div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import Navbar from '~/components/Navbar.vue'
import { useQuery } from '@vue/apollo-composable'
import gql from 'graphql-tag'

const isLoggedIn = ref(false)
const userName = ref('')
const userId = ref(null)
const userEmail = ref('')

onMounted(() => {
  const token = localStorage.getItem('token')
  if (!token) {
    window.location.href = '/login'
    return
  }
  isLoggedIn.value = true
  userName.value = localStorage.getItem('userName') || ''
  const raw = localStorage.getItem('userId') || ''
  const parsed = parseInt(raw, 10)
  userId.value = Number.isFinite(parsed) ? parsed : null
  userEmail.value = localStorage.getItem('userEmail') || ''
})

const { result } = useQuery(gql`
  query($uid: Int!) {
    recipes_aggregate(where: { user_id: { _eq: $uid } }) { aggregate { count } }
    recipes(where: { user_id: { _eq: $uid } }, order_by: { created_at: desc }, limit: 6) { id title featured_image }
  }
`, () => ({ uid: typeof userId.value === 'number' ? userId.value : 0 }))

const recipeCount = computed(() => result.value?.recipes_aggregate?.aggregate?.count ?? 0)
const recentRecipes = computed(() => result.value?.recipes || [])

// Change password
const oldPassword = ref('')
const newPassword = ref('')
const pwdLoading = ref(false)
const pwdError = ref('')
const pwdSuccess = ref('')
const backend = useRuntimeConfig().public.NUXT_PUBLIC_BACKEND_ENDPOINT

async function onChangePassword() {
  pwdError.value = ''
  pwdSuccess.value = ''
  if (!oldPassword.value || !newPassword.value) { pwdError.value = 'Both fields are required.'; return }
  if (newPassword.value.length < 6) { pwdError.value = 'New password must be at least 6 characters'; return }
  try {
    pwdLoading.value = true
    const token = typeof window !== 'undefined' ? localStorage.getItem('token') : ''
    const res = await fetch(`${backend}/password/change`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', ...(token ? { Authorization: `Bearer ${token}` } : {}) },
      body: JSON.stringify({ old_password: oldPassword.value, new_password: newPassword.value })
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.error || 'Failed to change password')
    pwdSuccess.value = 'Password changed successfully.'
    oldPassword.value = ''
    newPassword.value = ''
  } catch (e) {
    pwdError.value = e.message || 'Failed to change password'
  } finally {
    pwdLoading.value = false
  }
}
</script>