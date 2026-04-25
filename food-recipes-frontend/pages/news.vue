<template>
  <div class="min-h-screen bg-gradient-to-br from-orange-50 to-white">
    <Navbar />
    <section class="container mx-auto px-4 pt-6 max-w-6xl">
      <h1 class="text-3xl md:text-4xl font-bold mb-4">Food News</h1>
      <div class="rounded-xl bg-white border border-gray-200 p-4 shadow-sm mb-6">
        <input v-model="search" type="text" placeholder="Search articles" class="w-full rounded-lg border border-gray-200 px-3 py-2" />
      </div>
      <div v-if="filtered.length === 0" class="text-center text-gray-400 py-10">No articles found.</div>
      <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-6">
        <article v-for="a in filtered" :key="a.id" class="bg-white rounded-xl border border-gray-200 p-4 shadow-sm">
          <img :src="a.image" class="w-full h-36 object-cover rounded" />
          <h3 class="mt-2 font-semibold">{{ a.title }}</h3>
          <p class="text-sm text-gray-600 line-clamp-3">{{ a.summary }}</p>
          <div class="mt-2 text-xs text-gray-500">{{ a.tag }}</div>
        </article>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import Navbar from '~/components/Navbar.vue'

const search = ref('')
const articles = ref([
  { id: 1, title: 'New Superfood Trends', summary: 'Discover the latest superfoods taking the culinary world by storm.', image: 'https://placehold.co/600x400?text=News+1', tag: 'Trends' },
  { id: 2, title: 'Chef Spotlight: Layla', summary: 'An in-depth interview with Chef Layla about modern cuisine.', image: 'https://placehold.co/600x400?text=News+2', tag: 'Interview' },
  { id: 3, title: 'Sustainable Cooking', summary: 'How to cook sustainably with local, seasonal ingredients.', image: 'https://placehold.co/600x400?text=News+3', tag: 'Sustainability' }
])

const filtered = computed(() => {
  const s = search.value.trim().toLowerCase()
  return articles.value.filter(x => !s || x.title.toLowerCase().includes(s))
})
</script>
  