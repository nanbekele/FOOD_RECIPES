<template>
  <nav class="flex items-center bg-white shadow px-6 py-3 justify-between">
    <div class="flex items-center">
      <!-- Hamburger Menu -->
      <button class="mr-4 focus:outline-none md:hidden">
        <svg class="w-7 h-7" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" d="M4 6h16M4 12h16M4 18h16"/>
        </svg>
      </button>

      <!-- Logo -->
      <NuxtLink to="/" class="flex items-center mr-8">
        <img src="/logo.jpg" alt="Logo" class="h-10 w-10 mr-2" />
      </NuxtLink>

      <!-- Navigation Links -->
      <div class="flex space-x-8">
        <NuxtLink
          v-for="item in navItems"
          :key="item.name"
          :to="item.to"
          class="relative font-semibold text-lg"
          :class="{
            'text-black': $route.path !== item.to,
            'text-black font-bold': $route.path === item.to
          }"
        >
          {{ item.name }}
          <span
            v-if="$route.path === item.to"
            class="absolute left-0 -bottom-1 w-full h-1 bg-red-600 rounded"
          ></span>
        </NuxtLink>
      </div>
    </div>

    <!-- Auth/User Links -->
    <div class="flex items-center space-x-4">
      <template v-if="isLoggedIn">
        <!-- Dropdown Menu -->
        <div class="relative" @keydown.esc="showMenu = false">
          <button @click="toggleMenu" class="flex items-center focus:outline-none menu-toggle-button">
            <span class="inline-flex items-center justify-center h-10 w-10 rounded-full bg-gray-200 text-gray-700 font-bold text-lg">
              <template v-if="userInitials">{{ userInitials }}</template>
              <template v-else>
                <svg class="h-6 w-6 text-gray-400" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24" stroke-linecap="round" stroke-linejoin="round">
                  <circle cx="12" cy="7" r="4"/>
                  <path d="M5.5 21a8.38 8.38 0 0 1 13 0"/>
                </svg>
              </template>
            </span>
            <svg class="ml-2 h-4 w-4 text-gray-500" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path d="M19 9l-7 7-7-7"/>
            </svg>
          </button>

          <!-- Dropdown Content -->
          <div v-show="showMenu" ref="menuRef" class="absolute right-0 mt-2 w-48 bg-white border border-gray-200 rounded shadow-lg z-50">
            <NuxtLink to="/profile" @click="closeMenuAfterNavigation" class="block px-4 py-2 text-gray-700 hover:bg-gray-100">
              My Account
            </NuxtLink>
            <NuxtLink v-if="isAdmin" to="/admin/chefs" @click="closeMenuAfterNavigation" class="block px-4 py-2 text-gray-700 hover:bg-gray-100">
              Admin Dashboard
            </NuxtLink>
            <NuxtLink to="/recipes/create" @click="closeMenuAfterNavigation" class="block px-4 py-2 text-gray-700 hover:bg-gray-100">
              Create Recipe
            </NuxtLink>
            <NuxtLink to="/favorites" @click="closeMenuAfterNavigation" class="block px-4 py-2 text-gray-700 hover:bg-gray-100">
              Favorites
            </NuxtLink>
            <button @click="logout" class="w-full text-left block px-4 py-2 text-red-600 hover:bg-gray-100">
              Logout
            </button>
          </div>
        </div>
      </template>
      <template v-else>
        <NuxtLink to="/login" class="text-blue-600 hover:underline">Login</NuxtLink>
        <NuxtLink to="/signup" class="text-blue-600 hover:underline">Sign Up</NuxtLink>
      </template>
    </div>
  </nav>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, watchEffect, nextTick } from 'vue'

const navItems = [
  { name: 'Recipes', to: '/recipes' },
  { name: 'Top fav', to: '/top-fav' },
  { name: 'Chefs', to: '/chefs' },
  { name: 'Food News', to: '/news' }
]

const isLoggedIn = ref(false)
const showMenu = ref(false)
const menuRef = ref(null)
const userInitials = ref('')
const role = ref('')
const isAdmin = computed(() => role.value === 'admin')

onMounted(() => {
  updateAuthState()
  document.addEventListener('click', handleClickOutside)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', handleClickOutside)
})

const updateAuthState = () => {
  if (process.client) {
    isLoggedIn.value = !!localStorage.getItem('token')
    const name = localStorage.getItem('userName')
    if (name) {
      const parts = name.split(' ')
      userInitials.value = parts.map(p => p[0]).join('').toUpperCase().slice(0, 2)
    } else {
      userInitials.value = ''
    }
    role.value = localStorage.getItem('role') || ''
  }
}

// Watch for login state changes
watchEffect(() => {
  if (process.client) {
    isLoggedIn.value = !!localStorage.getItem('token')
    role.value = localStorage.getItem('role') || ''
  }
})

const toggleMenu = () => {
  showMenu.value = !showMenu.value
}

const handleClickOutside = (event) => {
  const isClickInside = menuRef.value && menuRef.value.contains(event.target)
  const isToggleButton = event.target.closest('.menu-toggle-button')
  if (showMenu.value && !isClickInside && !isToggleButton) {
    showMenu.value = false
  }
}

// 🔧 Fix for dropdown navigation issue
const closeMenuAfterNavigation = async () => {
  await nextTick()
  showMenu.value = false
}

const logout = () => {
  if (process.client) {
    localStorage.removeItem('token')
    localStorage.removeItem('userName')
    localStorage.removeItem('userEmail')
    localStorage.removeItem('userId')
    localStorage.removeItem('role')
    localStorage.removeItem('is_verified')
  }
  window.location.href = '/'
}
</script>
