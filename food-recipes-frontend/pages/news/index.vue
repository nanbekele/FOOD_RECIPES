<template>
  <div class="min-h-screen bg-gradient-to-br from-orange-50 to-white">
    <Navbar />
    <section class="container mx-auto px-4 pt-6 max-w-6xl">
      <div class="flex items-center justify-between mb-4">
        <h1 class="text-3xl md:text-4xl font-bold">Food News</h1>
        <a v-if="isLoggedIn" :href="`${basePath}/news/create`" class="px-4 py-2 rounded bg-green-600 text-white no-underline">Create News</a>
      </div>
      <div class="rounded-xl bg-white border border-gray-200 p-4 shadow-sm mb-6">
        <input v-model="search" type="text" placeholder="Search articles" class="w-full rounded-lg border border-gray-200 px-3 py-2" />
      </div>
      <div v-if="loading" class="text-center text-gray-500 py-10">Loading...</div>
      <div v-else-if="error" class="text-center text-red-600 py-10">{{ error }}</div>
      <div v-else>
        <div v-if="filtered.length === 0" class="text-center text-gray-400 py-10">No articles found.</div>
        <div v-else class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-6">
          <article v-for="a in filtered" :key="a.id" class="bg-white rounded-xl border border-gray-200 p-4 shadow-sm">
            <img v-if="a.image" :src="a.image" class="w-full h-36 object-cover rounded" />
            <div v-else class="w-full h-36 bg-gray-100 rounded"></div>
            <h3 class="mt-2 font-semibold">{{ a.title }}</h3>
            <p class="text-sm text-gray-600 line-clamp-3">{{ a.summary || a.content }}</p>
            <div class="mt-2 text-xs text-gray-500">{{ a.tag }}</div>
          </article>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import Navbar from '~/components/Navbar.vue'

const search = ref('')
const backend = useRuntimeConfig().public.NUXT_PUBLIC_BACKEND_ENDPOINT
const loading = ref(true)
const error = ref('')
const articles = ref([])

const basePath = computed(() => {
  if (typeof window === 'undefined') return ''
  return window.location.pathname.startsWith('/app') ? '/app' : ''
})

const isLoggedIn = computed(() => {
  if (typeof window === 'undefined') return false
  return !!localStorage.getItem('token')
})

async function fetchNews() {
  loading.value = true
  error.value = ''
  try {
    const res = await fetch(`${backend}/news`)
    const data = await res.json()
    if (!res.ok) throw new Error(data?.error || 'Failed to load news')
    articles.value = data?.news || []
  } catch (e) {
    error.value = e?.message || 'Failed to load news'
  }
  loading.value = false
}

onMounted(fetchNews)

const filtered = computed(() => {
  const s = search.value.trim().toLowerCase()
  return articles.value.filter(x => !s || x.title?.toLowerCase().includes(s))
})
</script>
