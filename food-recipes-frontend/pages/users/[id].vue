<template>
  <div class="min-h-screen bg-gradient-to-br from-orange-50 to-white">
    <Navbar />
    <section class="container mx-auto px-4 pt-6 max-w-6xl">
      <div class="flex items-center justify-between mb-6">
        <div class="flex items-center gap-3">
          <span class="inline-flex items-center justify-center h-14 w-14 rounded-full bg-orange-100 text-orange-700 font-bold text-xl">{{ initials }}</span>
          <div>
            <h1 class="text-2xl md:text-3xl font-bold flex items-center gap-2">
              <span>{{ user.name }}</span>
              <span v-if="user.role === 'chef' && user.is_verified" class="inline-flex items-center text-xs px-2 py-0.5 rounded bg-green-100 text-green-700">Verified</span>
            </h1>
            <p class="text-gray-600" v-if="user.bio">{{ user.bio }}</p>
          </div>
        </div>
        <button
          v-if="canFollow"
          class="px-3 py-2 rounded bg-gray-100"
          :disabled="followLoading"
          @click="toggleFollow"
        >
          {{ isFollowing ? 'Unfollow' : 'Follow' }}
        </button>
      </div>

      <div class="rounded-xl bg-white border border-gray-200 p-4 shadow-sm">
        <h2 class="text-xl font-semibold mb-4">Recipes by {{ user.name }}</h2>
        <div v-if="recipes.length === 0" class="text-center text-gray-400 py-10">No recipes yet.</div>
        <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-6">
          <RecipeCard
            v-for="r in recipes"
            :key="r.id"
            :recipe="r"
            :favorites="favorites"
            @refresh-favorites="refetchFavorites"
          />
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { computed, ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useQuery } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import Navbar from '~/components/Navbar.vue'
import RecipeCard from '~/components/RecipeCard.vue'

const route = useRoute()
const userId = parseInt(route.params.id, 10)

const backend = useRuntimeConfig().public.NUXT_PUBLIC_BACKEND_ENDPOINT
const userRaw = ref({ id: userId, name: '', bio: '', role: '', is_verified: false })
const recipesRaw = ref([])
const pageLoading = ref(true)
const pageError = ref('')

async function fetchProfile() {
  pageLoading.value = true
  pageError.value = ''
  try {
    const [uRes, rRes] = await Promise.all([
      fetch(`${backend}/users/${userId}`),
      fetch(`${backend}/users/${userId}/recipes`)
    ])
    const uData = await uRes.json()
    const rData = await rRes.json()
    if (!uRes.ok) throw new Error(uData?.error || 'Failed to load user')
    if (!rRes.ok) throw new Error(rData?.error || 'Failed to load recipes')
    userRaw.value = uData?.user || userRaw.value
    recipesRaw.value = rData?.recipes || []
  } catch (e) {
    pageError.value = e?.message || 'Failed to load'
  }
  pageLoading.value = false
}

onMounted(fetchProfile)

const user = computed(() => userRaw.value)
const recipes = computed(() => recipesRaw.value)
const initials = computed(() => (user.value.name || '').split(' ').map(p => p[0]).join('').toUpperCase().slice(0,2))

const isLoggedIn = computed(() => {
  if (typeof window === 'undefined') return false
  return !!localStorage.getItem('token')
})
const canFollow = computed(() => isLoggedIn.value && viewerId.value > 0 && viewerId.value !== userId)

const isFollowing = ref(false)
const followLoading = ref(false)

async function refreshFollowStatus() {
  if (!canFollow.value) return
  try {
    const token = localStorage.getItem('token') || ''
    const res = await fetch(`${backend}/users/${userId}/following`, {
      headers: token ? { Authorization: `Bearer ${token}` } : {}
    })
    const data = await res.json()
    if (!res.ok) return
    isFollowing.value = !!data?.following
  } catch {
    // ignore
  }
}

onMounted(refreshFollowStatus)

async function toggleFollow() {
  if (!canFollow.value || followLoading.value) return
  followLoading.value = true
  try {
    const token = localStorage.getItem('token') || ''
    const res = await fetch(`${backend}/users/${userId}/follow`, {
      method: isFollowing.value ? 'DELETE' : 'POST',
      headers: token ? { Authorization: `Bearer ${token}` } : {}
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data?.error || 'Request failed')
    isFollowing.value = !isFollowing.value
  } catch (e) {
    alert(e?.message || 'Request failed')
  }
  followLoading.value = false
}

// Viewer favorites for clickable state
const viewerId = computed(() => {
  if (typeof window === 'undefined') return 0
  const raw = localStorage.getItem('userId') || ''
  const parsed = parseInt(raw, 10)
  return Number.isFinite(parsed) ? parsed : 0
})
const { result: favResult, refetch: refetchFavorites } = useQuery(gql`
  query($uid: Int!) {
    favorites(where: { user_id: { _eq: $uid } }) { recipe { id } }
  }
`, () => ({ uid: viewerId.value }))
const favorites = computed(() => favResult.value?.favorites || [])
</script>