<template>
  <div class="min-h-screen bg-gradient-to-br from-orange-50 to-white">
    <Navbar />
    <section class="container mx-auto px-4 pt-6 max-w-6xl">
      <div class="flex items-center justify-between mb-4">
        <h1 class="text-3xl font-bold">Chef Approvals</h1>
        <div class="flex gap-2">
          <button class="px-3 py-2 rounded bg-gray-100" @click="filter='all'" :class="{ 'bg-gray-200': filter==='all' }">All</button>
          <button class="px-3 py-2 rounded bg-gray-100" @click="filter='pending'" :class="{ 'bg-gray-200': filter==='pending' }">Pending</button>
          <button class="px-3 py-2 rounded bg-gray-100" @click="filter='verified'" :class="{ 'bg-gray-200': filter==='verified' }">Verified</button>
        </div>
      </div>

      <div class="bg-white border border-gray-200 rounded-xl p-4 shadow-sm">
        <div v-if="loading" class="text-center text-gray-500 py-10">Loading...</div>
        <div v-else-if="error" class="text-center text-red-600 py-10">{{ error }}</div>
        <div v-else>
          <div v-if="rows.length === 0" class="text-center text-gray-400 py-10">No chefs found.</div>
          <table v-else class="w-full text-left">
            <thead>
              <tr class="text-sm text-gray-500">
                <th class="py-2">ID</th>
                <th class="py-2">Name</th>
                <th class="py-2">Email</th>
                <th class="py-2">Verified</th>
                <th class="py-2">Action</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="u in rows" :key="u.id" class="border-t">
                <td class="py-2">{{ u.id }}</td>
                <td class="py-2">{{ u.name }}</td>
                <td class="py-2">{{ u.email }}</td>
                <td class="py-2">
                  <span v-if="u.is_verified" class="inline-flex items-center text-xs px-2 py-0.5 rounded bg-green-100 text-green-700">Verified</span>
                  <span v-else class="inline-flex items-center text-xs px-2 py-0.5 rounded bg-yellow-100 text-yellow-700">Pending</span>
                </td>
                <td class="py-2">
                  <button v-if="!u.is_verified" class="px-3 py-1 rounded bg-green-600 text-white" @click="setVerify(u.id,true)">Approve</button>
                  <button v-else class="px-3 py-1 rounded bg-gray-200" @click="setVerify(u.id,false)">Revoke</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import Navbar from '~/components/Navbar.vue'

const backend = useRuntimeConfig().public.NUXT_PUBLIC_BACKEND_ENDPOINT
const loading = ref(true)
const error = ref('')
const filter = ref('pending') // all|pending|verified
const chefs = ref([])

async function fetchChefs() {
  loading.value = true
  error.value = ''
  try {
    const token = typeof window !== 'undefined' ? localStorage.getItem('token') : ''
    const qs = filter.value === 'all' ? '' : (filter.value === 'pending' ? '?verified=false' : '?verified=true')
    const res = await fetch(`${backend}/admin/chefs${qs}`, {
      headers: token ? { Authorization: `Bearer ${token}` } : {}
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.error || 'Failed to load')
    chefs.value = data.chefs || []
  } catch (e) {
    error.value = e.message || 'Failed to load chefs'
  }
  loading.value = false
}

async function setVerify(id, v) {
  try {
    const token = typeof window !== 'undefined' ? localStorage.getItem('token') : ''
    const res = await fetch(`${backend}/admin/chefs/${id}/verify`, {
      method: 'PATCH',
      headers: {
        'Content-Type': 'application/json',
        ...(token ? { Authorization: `Bearer ${token}` } : {})
      },
      body: JSON.stringify({ is_verified: !!v })
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.error || 'Failed to update')
    // refresh list
    await fetchChefs()
  } catch (e) {
    alert(e.message || 'Update failed')
  }
}

const rows = computed(() => chefs.value)

onMounted(fetchChefs)
watch(filter, fetchChefs)
</script>
