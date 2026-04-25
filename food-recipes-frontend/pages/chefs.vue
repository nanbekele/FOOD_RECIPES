<template>
  <div class="min-h-screen bg-gradient-to-br from-orange-50 to-white">
    <Navbar />
    <section class="container mx-auto px-4 pt-6 max-w-6xl">
      <div class="flex items-center justify-between mb-6">
        <h1 class="text-3xl md:text-4xl font-bold">Chefs & Creators</h1>
        <NuxtLink to="/recipes" class="text-orange-700 hover:underline">Browse Recipes</NuxtLink>
      </div>
      <div class="rounded-xl bg-white border border-gray-200 p-4 shadow-sm mb-6">
        <input v-model="search" type="text" placeholder="Search creators by name" class="w-full rounded-lg border border-gray-200 px-3 py-2" />
      </div>
      <div v-if="loading" class="text-center text-gray-500 py-10">Loading...</div>
      <div v-else-if="error" class="text-center text-red-600 py-10">{{ error.message }}</div>
      <div v-else>
        <div v-if="filteredUsers.length === 0" class="text-center text-gray-400 py-10">No creators found.</div>
        <div v-else class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-6">
          <NuxtLink v-for="u in filteredUsers" :key="u.id" :to="`/users/${u.id}`" class="block">
            <CreatorCard :creator="u" />
          </NuxtLink>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import Navbar from '~/components/Navbar.vue'
import CreatorCard from '~/components/CreatorCard.vue'

const search = ref('')
const backend = useRuntimeConfig().public.NUXT_PUBLIC_BACKEND_ENDPOINT
const loading = ref(true)
const error = ref(null)
const usersRaw = ref([])

async function fetchChefs() {
  loading.value = true
  error.value = null
  try {
    const res = await fetch(`${backend}/chefs`)
    const data = await res.json()
    if (!res.ok) throw new Error(data?.error || 'Failed to load creators')
    usersRaw.value = data?.chefs || []
  } catch (e) {
    error.value = e
  }
  loading.value = false
}

onMounted(fetchChefs)

const users = computed(() => usersRaw.value || [])

const filteredUsers = computed(() => {
  const s = search.value.trim().toLowerCase()
  return users.value.filter(u => !s || u.name?.toLowerCase().includes(s))
})
</script>