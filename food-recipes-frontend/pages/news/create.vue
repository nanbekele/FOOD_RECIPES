<template>
  <div class="min-h-screen bg-gradient-to-br from-orange-50 to-white">
    <Navbar />
    <section class="container mx-auto px-4 pt-6 max-w-3xl">
      <div class="rounded-2xl bg-white border border-gray-200 p-6 shadow-sm">
        <h1 class="text-3xl md:text-4xl font-bold mb-4">Create News</h1>

        <form @submit.prevent="submit" class="space-y-4">
          <div>
            <label class="block text-sm font-medium">Title</label>
            <input v-model="title" class="w-full rounded border border-gray-200 px-3 py-2" placeholder="e.g., Sustainable Cooking" />
          </div>

          <div>
            <label class="block text-sm font-medium">Tag</label>
            <input v-model="tag" class="w-full rounded border border-gray-200 px-3 py-2" placeholder="e.g., Trends" />
          </div>

          <div>
            <label class="block text-sm font-medium">Summary</label>
            <textarea v-model="summary" rows="3" class="w-full rounded border border-gray-200 px-3 py-2" placeholder="Short summary..."></textarea>
          </div>

          <div>
            <label class="block text-sm font-medium">Content</label>
            <textarea v-model="content" rows="7" class="w-full rounded border border-gray-200 px-3 py-2" placeholder="Write the full article..."></textarea>
          </div>

          <div>
            <label class="block text-sm font-medium">Image</label>
            <input type="file" accept="image/*" @change="onImageSelected" class="w-full rounded border border-gray-200 px-3 py-2" />
            <div v-if="imagePreview" class="mt-3">
              <img :src="imagePreview" class="w-full max-h-64 object-cover rounded" />
            </div>
          </div>

          <div class="flex items-center gap-3">
            <button type="submit" class="px-4 py-2 bg-green-600 text-white rounded" :disabled="submitting || uploading">
              {{ submitting ? 'Creating...' : (uploading ? 'Uploading...' : 'Create News') }}
            </button>
            <NuxtLink to="/news" class="text-orange-700 hover:underline">Cancel</NuxtLink>
          </div>

          <p v-if="error" class="text-red-600">{{ error }}</p>
        </form>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, computed, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import Navbar from '~/components/Navbar.vue'

const backend = useRuntimeConfig().public.NUXT_PUBLIC_BACKEND_ENDPOINT

const title = ref('')
const tag = ref('')
const summary = ref('')
const content = ref('')
const imageUrl = ref('')
const imagePreview = ref('')

const error = ref('')
const uploading = ref(false)
const submitting = ref(false)

const router = useRouter()

function withBase(path) {
  const base = router.currentRoute.value?.path?.startsWith('/app') ? '/app' : ''
  return `${base}${path}`
}

const isLoggedIn = computed(() => {
  if (typeof window === 'undefined') return false
  return !!localStorage.getItem('token')
})

async function onImageSelected(e) {
  const file = (e.target.files && e.target.files[0]) ? e.target.files[0] : null
  if (!file) return

  if (imagePreview.value && String(imagePreview.value).startsWith('blob:')) {
    URL.revokeObjectURL(imagePreview.value)
  }
  imagePreview.value = URL.createObjectURL(file)
  uploading.value = true
  error.value = ''

  const form = new FormData()
  form.append('files', file)
  try {
    const res = await fetch(`${backend}/upload`, { method: 'POST', body: form })
    const data = await res.json()
    if (!res.ok) throw new Error(data?.error || 'Upload failed')

    const item = Array.isArray(data.items) ? data.items[0] : (data.item || null)
    const url = item?.url || data.url || (Array.isArray(data.urls) ? data.urls[0] : '')
    if (!url) throw new Error('Upload failed')

    imageUrl.value = url
  } catch (e2) {
    error.value = e2?.message || 'Upload failed'
    imageUrl.value = ''
  }
  uploading.value = false
}

onBeforeUnmount(() => {
  if (imagePreview.value && String(imagePreview.value).startsWith('blob:')) {
    URL.revokeObjectURL(imagePreview.value)
  }
})

async function submit() {
  error.value = ''
  if (!isLoggedIn.value) {
    await navigateTo(withBase('/login'))
    return
  }
  if (!title.value.trim()) {
    error.value = 'Title is required'
    return
  }

  submitting.value = true
  try {
    const token = localStorage.getItem('token') || ''
    const res = await fetch(`${backend}/news`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...(token ? { Authorization: `Bearer ${token}` } : {})
      },
      body: JSON.stringify({
        title: title.value,
        tag: tag.value,
        summary: summary.value,
        content: content.value,
        image: imageUrl.value
      })
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data?.error || 'Create failed')

    await navigateTo(withBase('/news'))
  } catch (e3) {
    error.value = e3?.message || 'Create failed'
  }
  submitting.value = false
}
</script>
