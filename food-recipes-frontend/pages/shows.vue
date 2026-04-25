<template>
  <div class="min-h-screen bg-gradient-to-br from-orange-50 to-white">
    <Navbar />
    <section class="container mx-auto px-4 pt-6 max-w-6xl">
      <h1 class="text-3xl md:text-4xl font-bold mb-4">Cooking Shows</h1>
      <div class="rounded-xl bg-white border border-gray-200 p-4 shadow-sm mb-6">
        <input v-model="search" type="text" placeholder="Search shows" class="w-full rounded-lg border border-gray-200 px-3 py-2" />
      </div>
      <div v-if="filtered.length === 0" class="text-center text-gray-400 py-10">No shows found.</div>
      <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-6">
        <div v-for="s in filtered" :key="s.id" class="bg-white rounded-xl border border-gray-200 p-4 shadow-sm">
          <img :src="s.image" class="w-full h-36 object-cover rounded" />
          <h3 class="mt-2 font-semibold">{{ s.title }}</h3>
          <p class="text-sm text-gray-600">{{ s.description }}</p>
          <div class="mt-2 text-xs text-gray-500">{{ s.category }}</div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import Navbar from '~/components/Navbar.vue'

const search = ref('')
const shows = ref([
  { id: 1, title: '30-Minute Meals', description: 'Quick and tasty weeknight cooking', image: 'https://placehold.co/600x400?text=Show+1', category: 'Quick' },
  { id: 2, title: 'Vegan Delights', description: 'Plant-based recipes and techniques', image: 'https://placehold.co/600x400?text=Show+2', category: 'Vegan' },
  { id: 3, title: 'Dessert Masters', description: 'Sweet treats and baking tips', image: 'https://placehold.co/600x400?text=Dessert', category: 'Dessert' }
])

const filtered = computed(() => {
  const s = search.value.trim().toLowerCase()
  return shows.value.filter(x => !s || x.title.toLowerCase().includes(s))
})
</script>
