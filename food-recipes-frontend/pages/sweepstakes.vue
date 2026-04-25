<template>
  <div class="min-h-screen bg-gradient-to-br from-orange-50 to-white">
    <Navbar />
    <section class="container mx-auto px-4 pt-6 max-w-6xl">
      <h1 class="text-3xl md:text-4xl font-bold mb-4">Sweepstakes</h1>
      <div class="rounded-xl bg-white border border-gray-200 p-6 shadow-sm">
        <p class="text-gray-600 mb-4">Enter for a chance to win cooking gear and premium recipe packs.</p>
        <form @submit.prevent="submitEntry" class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <input v-model="name" placeholder="Full name" class="rounded-lg border border-gray-200 px-3 py-2" required />
          <input v-model="email" type="email" placeholder="Email" class="rounded-lg border border-gray-200 px-3 py-2" required />
          <div class="md:col-span-2">
            <button class="px-4 py-2 rounded bg-orange-600 text-white">Submit Entry</button>
            <span v-if="submitted" class="ml-3 text-green-700">Thanks for entering!</span>
          </div>
        </form>
        <p class="mt-4 text-xs text-gray-500">By entering, you agree to official rules and privacy policy. No purchase necessary.</p>
      </div>
    </section>
  </div>
</template>

<script setup>
import Navbar from '~/components/Navbar.vue'
import { ref } from 'vue'

const name = ref('')
const email = ref('')
const submitted = ref(false)

function submitEntry() {
  submitted.value = true
}
</script>
  