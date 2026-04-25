<template>
  <div class="flex flex-col items-center justify-center min-h-screen bg-gray-50">
    <form @submit.prevent="onSubmit" class="bg-white p-8 rounded shadow-md w-full max-w-sm relative">
      <h2 class="text-2xl font-bold mb-6 text-center">Sign Up</h2>
      <input v-model="name" type="text" placeholder="Name" class="input mb-4" required />
      <input v-model="email" type="email" placeholder="Email" class="input mb-4" required />
      <select v-model="role" class="input mb-4">
        <option value="guest">Guest</option>
        <option value="chef">Chef</option>
      </select>
      <input v-model="password" type="password" placeholder="Password" class="input mb-4" required />
      <button type="submit" class="btn w-full" :disabled="loading">
        <span v-if="loading" class="loader mr-2"></span>
        Sign Up
      </button>
      <div v-if="success" class="text-green-600 mt-2">{{ success }}</div>
      <div v-if="error" class="text-red-500 mt-2">{{ error }}</div>
    </form>
  </div>
</template>

<script setup>
import { ref } from 'vue'
const name = ref('')
const email = ref('')
const role = ref('guest')
const password = ref('')
const error = ref('')
const success = ref('')
const loading = ref(false)
const backend = useRuntimeConfig().public.NUXT_PUBLIC_BACKEND_ENDPOINT

const onSubmit = async () => {
  error.value = ''
  success.value = ''
  loading.value = true
  try {
    const res = await fetch(`${backend}/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ name: name.value, email: email.value, role: role.value, password: password.value }),
    })
    const data = await res.json()
    if (res.ok) {
      success.value = 'Registration successful! Please log in.'
      setTimeout(() => window.location.href = '/login', 1200)
    } else {
      error.value = data.error || 'Signup failed'
    }
  } catch (e) {
    error.value = 'Signup failed: network or server error'
  }
  loading.value = false
}
</script>

<style scoped>
.input { @apply border p-2 rounded w-full focus:outline-none focus:ring-2 focus:ring-blue-400 transition; }
.btn { @apply bg-blue-600 text-white py-2 rounded hover:bg-blue-700 transition; }
.loader {
  border: 3px solid #f3f3f3;
  border-top: 3px solid #3490dc;
  border-radius: 50%;
  width: 18px;
  height: 18px;
  animation: spin 1s linear infinite;
  display: inline-block;
  vertical-align: middle;
}
@keyframes spin {
  0% { transform: rotate(0deg);}
  100% { transform: rotate(360deg);}
}
.fade-enter-active, .fade-leave-active { transition: opacity 0.3s; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
