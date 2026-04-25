<template>
  <div class="min-h-screen bg-gradient-to-br from-orange-50 to-white">
    <Navbar />
    <section class="container mx-auto px-4 pt-6 max-w-6xl">
      <div class="flex items-center justify-between mb-4">
        <h1 class="text-3xl md:text-4xl font-bold">Shop</h1>
        <NuxtLink to="/recipes" class="text-orange-700 hover:underline">Browse Recipes</NuxtLink>
      </div>
      <div class="rounded-xl bg-white border border-gray-200 p-4 shadow-sm mb-6">
        <p class="text-gray-600">Explore cooking tools, premium recipes, and bundles.</p>
      </div>
      <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-6">
        <div v-for="p in products" :key="p.id" class="bg-white rounded-xl border border-gray-200 p-4 shadow-sm">
          <img :src="p.image" class="w-full h-36 object-cover rounded" />
          <h3 class="mt-2 font-semibold">{{ p.title }}</h3>
          <p class="text-sm text-gray-600">{{ p.description }}</p>
          <div class="mt-3 flex items-center justify-between">
            <span class="font-semibold">${{ p.price }}</span>
            <button class="px-3 py-1 rounded bg-orange-600 text-white">Buy</button>
          </div>
        </div>
      </div>

      <div class="mt-8 rounded-xl bg-gradient-to-r from-orange-100 to-yellow-100 p-6">
        <h2 class="text-2xl font-bold mb-2">Premium Recipe Packs</h2>
        <p class="text-gray-700 mb-3">Unlock curated, step-by-step recipes from top chefs.</p>
        <NuxtLink to="/recipes" class="px-4 py-2 rounded bg-orange-600 text-white">Browse Premium Recipes</NuxtLink>
      </div>
    </section>
  </div>
</template>

<script setup>
import Navbar from '~/components/Navbar.vue'
import { ref } from 'vue'

const products = ref([
  { id: 1, title: 'Chef Knife', description: 'Precision stainless steel knife.', image: 'https://placehold.co/600x400?text=Knife', price: 39 },
  { id: 2, title: 'Cast Iron Pan', description: 'For perfect searing and baking.', image: 'https://placehold.co/600x400?text=Pan', price: 59 },
  { id: 3, title: 'Premium Recipe Pack', description: '10 exclusive recipes by top chefs.', image: 'https://placehold.co/600x400?text=Recipes', price: 19 }
])
</script>
  