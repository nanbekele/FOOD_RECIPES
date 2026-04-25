<template>
  <div class="flex flex-col items-center justify-center min-h-screen bg-gray-50">
    <form @submit.prevent="onSubmit" class="bg-white p-8 rounded shadow-md w-full max-w-sm">
      <h2 class="text-2xl font-bold mb-6 text-center">Login</h2>
      <input v-model="email" type="email" placeholder="Email" class="input mb-4" required />
      <input v-model="password" type="password" placeholder="Password" class="input mb-4" required />
      <button type="submit" class="btn w-full">Login</button>
      <div v-if="error" class="text-red-500 mt-2">{{ error }}</div>
      <div class="text-center mt-4">
        <NuxtLink to="/signup" class="text-blue-600 hover:underline">Don't have an account? Sign up</NuxtLink>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const email = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)
const backend = useRuntimeConfig().public.NUXT_PUBLIC_BACKEND_ENDPOINT

const onSubmit = async () => {
  error.value = ''
  loading.value = true
  try {
    const emailNormalized = (email.value || '').trim().toLowerCase()
    const res = await fetch(`${backend}/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email: emailNormalized, password: password.value }),
    })
    const data = await res.json()
    if (res.ok) {
      localStorage.setItem('token', data.token)
      localStorage.setItem('userName', data.name)
      localStorage.setItem('userEmail', data.email)
      localStorage.setItem('userId', data.id)
      if (data.role) localStorage.setItem('role', data.role)
      if (typeof data.is_verified !== 'undefined') localStorage.setItem('is_verified', String(!!data.is_verified))
      sessionStorage.setItem('justLoggedIn', 'true')
      window.location.href = '/'
    } else {
      error.value = data.error
    }
  } catch (e) {
    error.value = 'Login failed'
  }
  loading.value = false
}
</script>

<style scoped>
.input { @apply border p-2 rounded w-full focus:outline-none focus:ring-2 focus:ring-blue-400 transition; }
.btn { @apply bg-blue-600 text-white py-2 rounded hover:bg-blue-700 transition; }
</style>
