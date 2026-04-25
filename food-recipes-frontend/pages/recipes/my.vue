<template>
  <div class="min-h-screen bg-gradient-to-br from-orange-50 to-white">
    <Navbar />
    <section class="container mx-auto px-4 pt-6 max-w-6xl">
      <div class="flex items-center justify-between mb-4">
        <h1 class="text-3xl font-bold">My Recipes</h1>
        <NuxtLink to="/recipes/create" class="px-3 py-2 rounded bg-orange-600 text-white">Create New</NuxtLink>
      </div>
      <div v-if="loading" class="text-center text-gray-500 py-10">Loading...</div>
      <div v-else-if="error" class="text-center text-red-600 py-10">{{ error }}</div>
      <div v-else>
        <div v-if="recipes.length === 0" class="text-center text-gray-400 py-10">You have no recipes yet.</div>
        <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-6">
          <div v-for="r in recipes" :key="r.id" class="group bg-white rounded-xl border border-gray-200 overflow-hidden shadow-sm hover:shadow-md transition">
            <div class="relative">
              <img :src="r.featured_image || 'https://placehold.co/600x320?text=No+Image'" class="w-full h-40 object-cover" />
              <div class="absolute inset-x-0 bottom-0 bg-gradient-to-t from-black/50 to-transparent p-3">
                <h3 class="text-white font-semibold truncate">{{ r.title }}</h3>
              </div>
            </div>
            <div class="p-4">
              <p class="text-sm text-gray-600 line-clamp-2 min-h-[2.5rem]">{{ r.description }}</p>
              <div class="mt-4 flex items-center gap-2">
                <NuxtLink :to="`/recipes/${r.id}`" class="px-3 py-1 rounded bg-gray-100 hover:bg-gray-200 transition">View Recipe</NuxtLink>
                <NuxtLink :to="`/recipes/edit/${r.id}`" class="px-3 py-1 rounded bg-blue-600 text-white hover:bg-blue-700 transition">Edit</NuxtLink>
                <button type="button" @click="confirmDelete(r.id)" class="ml-auto px-3 py-1 rounded bg-red-600 text-white hover:bg-red-700 transition" aria-label="Delete recipe">Delete</button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <Modal v-if="showDelete" @close="showDelete=false">
        <template #title>Delete Recipe</template>
        <template #body>
          <p>Are you sure you want to delete this recipe?</p>
        </template>
        <template #footer>
          <button @click="performDelete" class="px-3 py-2 bg-red-600 text-white rounded">Delete</button>
          <button @click="showDelete=false" class="px-3 py-2 bg-gray-100 rounded">Cancel</button>
        </template>
      </Modal>
    </section>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import Navbar from '~/components/Navbar.vue'
import Modal from '~/components/Modal.vue'
import { useQuery, useMutation } from '@vue/apollo-composable'
import gql from 'graphql-tag'

const userId = ref(null)
const showDelete = ref(false)
const deleteId = ref(null)

onMounted(() => {
  const raw = localStorage.getItem('userId') || ''
  const parsed = parseInt(raw, 10)
  userId.value = Number.isFinite(parsed) ? parsed : null
})

const queryEnabled = computed(() => Number.isFinite(userId.value) && userId.value > 0)

const { result, loading, error } = useQuery(gql`
  query($uid: Int!) {
    recipes(where: { user_id: { _eq: $uid } }) {
      id
      title
      description
      featured_image
    }
  }
`, () => ({ uid: userId.value }), { fetchPolicy: 'cache-and-network', enabled: queryEnabled })

const recipes = computed(() => result.value?.recipes || [])

const DELETE_BY_PK = gql`mutation($id: uuid!) { delete_recipes_by_pk(id: $id) { id } }`
const DELETE_GENERIC = gql`mutation($id: uuid!) { delete_recipes(where: { id: { _eq: $id } }) { affected_rows } }`

const { mutate: deleteByPk } = useMutation(DELETE_BY_PK)
const { mutate: deleteGeneric } = useMutation(DELETE_GENERIC)

function confirmDelete(id) {
  deleteId.value = id
  showDelete.value = true
}

async function performDelete() {
  showDelete.value = false
  try {
    const { data } = await deleteByPk({ id: deleteId.value })
    if (!data?.delete_recipes_by_pk?.id) throw new Error('pk delete failed')
  } catch (e) {
    try {
      const { data } = await deleteGeneric({ id: deleteId.value })
      if (!data?.delete_recipes?.affected_rows) throw new Error('generic delete failed')
    } catch (err) {
      alert('Delete failed. Ensure Hasura recipes table exists and permissions allow delete.')
      return
    }
  }
  // Refetch or filter out deleted locally
  const idx = recipes.value.findIndex(r => r.id === deleteId.value)
  if (idx >= 0) recipes.value.splice(idx,1)
}
</script>