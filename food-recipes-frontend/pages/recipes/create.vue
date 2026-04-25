<template>
  <div class="min-h-screen bg-gradient-to-br from-orange-50 to-white">
    <Navbar />
    <section class="container mx-auto px-4 pt-6 max-w-6xl">
      <div class="rounded-2xl bg-white border border-gray-200 p-6 shadow-sm">
        <h1 class="text-3xl md:text-4xl font-bold mb-4">Create a Recipe</h1>
        <form @submit.prevent="submitRecipe" class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label class="block text-sm font-medium">Title</label>
            <input v-model="title" class="w-full rounded border border-gray-200 px-3 py-2" placeholder="e.g., Quick Garlic Shrimp" />
          </div>

          <div>
            <label class="block text-sm font-medium">Prep Time (minutes)</label>
            <input v-model.number="prepTime" type="number" min="0" class="w-full rounded border border-gray-200 px-3 py-2" placeholder="e.g., 20" />
          </div>

          <div class="md:col-span-2">
            <label class="block text-sm font-medium">Description</label>
            <textarea v-model="description" rows="4" class="w-full rounded border border-gray-200 px-3 py-2" placeholder="Describe the dish"></textarea>
          </div>

          <div>
            <label class="block text-sm font-medium">Category</label>
            <select v-model="categoryId" class="w-full rounded border border-gray-200 px-3 py-2">
              <option :value="null">Select a category</option>
              <option v-for="c in categories" :key="c.id" :value="c.id">{{ c.name }}</option>
            </select>
          </div>

          <div>
            <label class="block text-sm font-medium">Media (images or videos)</label>
            <input type="file" multiple accept="image/*,video/*,.mp4,.mov,.avi,.mkv,.webm,.m4v,.3gp,.mpeg,.mpg" @change="onMediaSelected" class="w-full rounded border border-gray-200 px-3 py-2" />
            <p class="text-xs text-gray-400 mt-1">You can upload up to 3 files (images and/or videos, e.g., mp4, mov, avi, mkv, webm).</p>
            <div v-if="mediaItems.length" class="mt-3 grid grid-cols-3 gap-3">
              <div v-for="(m, i) in mediaItems" :key="(m.url || m.preview || i) + '-' + i" class="relative group">
                <div class="w-full h-24 bg-gray-100 rounded flex items-center justify-center overflow-hidden">
                  <template v-if="m.type === 'video'">
                    <video v-if="(m.preview || m.url) && !m.broken" :src="displaySrc(m)" class="w-full h-full object-contain" controls playsinline preload="metadata" @error="onMediaError(i)"></video>
                    <div v-else class="text-xs text-gray-500">Preview unavailable</div>
                  </template>
                  <template v-else>
                    <img v-if="(m.preview || m.url) && !m.broken" :src="displaySrc(m)" class="w-full h-full object-contain" @error="onMediaError(i)" @load="m.broken=false" />
                    <div v-else class="text-xs text-gray-500">Preview unavailable</div>
                  </template>
                </div>
                <button type="button" @click="removeMedia(i)" class="absolute -top-2 -right-2 bg-white border border-gray-300 text-gray-700 w-6 h-6 rounded-full flex items-center justify-center shadow-sm hover:bg-red-600 hover:text-white" aria-label="Remove">
                  ×
                </button>
                <label class="absolute top-2 left-2 bg-white bg-opacity-80 rounded px-2 py-1 text-xs">
                  <input type="radio" name="featured" :value="i" v-model="featuredIndex" /> Featured
                </label>
              </div>
            </div>
          </div>

          <!-- Ingredients Section -->
          <div class="md:col-span-2">
            <label class="block text-sm font-medium mb-2">Ingredients</label>
            <div class="space-y-3">
              <div v-for="(ing, idx) in ingredients" :key="idx" class="grid grid-cols-1 md:grid-cols-3 gap-3 items-center">
                <input v-model="ing.name" placeholder="Ingredient name" class="rounded border border-gray-200 px-3 py-2" />
                <input v-model="ing.quantity" placeholder="Quantity (e.g., 2 cups)" class="rounded border border-gray-200 px-3 py-2" />
                <div class="flex gap-2">
                  <button type="button" @click="removeIngredient(idx)" class="px-3 py-2 bg-gray-200 rounded" v-if="ingredients.length > 1">Remove</button>
                </div>
              </div>
              <button type="button" @click="addIngredient" class="px-3 py-2 bg-blue-600 text-white rounded">Add Ingredient</button>
            </div>
          </div>

          <!-- Steps Section -->
          <div class="md:col-span-2">
            <label class="block text-sm font-medium mb-2">Steps</label>
            <div class="space-y-3">
              <div v-for="(st, idx) in steps" :key="idx" class="grid grid-cols-1 md:grid-cols-8 gap-3 items-start">
                <div class="md:col-span-1">
                  <input :value="idx + 1" disabled class="w-full rounded border border-gray-200 px-3 py-2 text-center bg-gray-50" />
                </div>
                <div class="md:col-span-6">
                  <textarea v-model="st.description" rows="2" placeholder="Describe this step" class="w-full rounded border border-gray-200 px-3 py-2"></textarea>
                </div>
                <div class="md:col-span-1 flex items-center">
                  <button type="button" @click="removeStep(idx)" class="px-3 py-2 bg-gray-200 rounded" v-if="steps.length > 1">Remove</button>
                </div>
              </div>
              <button type="button" @click="addStep" class="px-3 py-2 bg-blue-600 text-white rounded">Add Step</button>
              <p class="text-xs text-gray-400">At least one step is required.</p>
            </div>
          </div>

          <div class="md:col-span-2 flex items-center gap-3">
            <button type="submit" class="px-4 py-2 bg-green-600 text-white rounded">Create Recipe</button>
            <NuxtLink to="/recipes" class="text-orange-700 hover:underline">Cancel</NuxtLink>
          </div>

          <p v-if="error" class="md:col-span-2 mt-3 text-red-600">{{ error }}</p>
        </form>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import Navbar from '~/components/Navbar.vue'
import { useQuery, useMutation } from '@vue/apollo-composable'
import gql from 'graphql-tag'

const title = ref('')
const description = ref('')
const prepTime = ref(null)
const categoryId = ref(null)
const mediaItems = ref([]) // [{ url, preview?, type: 'image'|'video' }]
const featuredIndex = ref(0)
const error = ref('')
const backend = useRuntimeConfig().public.NUXT_PUBLIC_BACKEND_ENDPOINT
const ingredients = ref([{ name: '', quantity: '' }])
const steps = ref([{ description: '' }])
const catError = ref('')
const uploading = ref(false)

const userId = ref(null)
onMounted(() => {
  const raw = typeof window !== 'undefined' ? localStorage.getItem('userId') || '' : ''
  const parsed = parseInt(raw, 10)
  userId.value = Number.isFinite(parsed) ? parsed : null
})

// Keep featured index valid after removals
watch(mediaItems, () => {
  if (featuredIndex.value >= mediaItems.value.length) featuredIndex.value = 0
})

// Categories (id type may be string or int depending on DB). Use as-is, with backend fallback.
const categoriesRef = ref([])
const { result: catResult, onError: onCatError } = useQuery(gql`query { categories { id name } }`)
onCatError(async () => {
  catError.value = 'Unable to load categories'
  await fetchCategoriesFallback()
})
watch(catResult, async (val) => {
  const rows = val?.categories || []
  if (rows.length) {
    categoriesRef.value = rows
    catError.value = ''
  } else {
    await fetchCategoriesFallback()
  }
}, { immediate: true })
const categories = computed(() => categoriesRef.value)

async function fetchCategoriesFallback() {
  try {
    const res = await fetch(`${backend}/categories`)
    const data = await res.json()
    categoriesRef.value = data?.categories || []
    if (!categoriesRef.value.length) {
      catError.value = 'No categories found'
    } else {
      catError.value = ''
    }
  } catch (err) {
    catError.value = 'Unable to load categories'
  }
}

async function onMediaSelected(e) {
  uploading.value = true
  const selected = Array.from(e.target.files || [])
  if (!selected.length) return
  const remaining = Math.max(0, 3 - mediaItems.value.length)
  if (remaining <= 0) { alert('You can upload up to 3 files only.'); return }
  if (selected.length > remaining) {
    alert(`You can add only ${remaining} more file(s). Extra files will be ignored.`)
  }
  const files = selected.slice(0, remaining)

  // Add local previews first (FileReader for images; object URL for videos)
  const previews = await Promise.all(files.map(file => new Promise(resolve => {
    const lower = file.name ? file.name.toLowerCase() : ''
    const isVideo = (file.type || '').startsWith('video') || /\.(mp4|mov|avi|mkv|webm|m4v|3gp|mpeg|mpg)$/i.test(lower)
    const isImage = (file.type || '').startsWith('image') || /\.(png|jpg|jpeg|gif|webp|bmp|heic|heif|tiff|svg)$/i.test(lower)
    if (isImage && !isVideo) {
      const reader = new FileReader()
      reader.onload = () => resolve({ type: 'image', preview: reader.result, file })
      reader.onerror = () => resolve({ type: 'image', preview: '', file })
      reader.readAsDataURL(file)
    } else {
      const url = URL.createObjectURL(file)
      resolve({ type: isVideo ? 'video' : 'image', preview: url, file })
    }
  })))
  for (const p of previews) mediaItems.value.push({ url: '', preview: p.preview, type: p.type })

  // Upload and update the just-added placeholders in order
  const form = new FormData()
  for (const f of files) form.append('files', f)
  try {
    const res = await fetch(`${backend}/upload`, { method: 'POST', body: form })
    const data = await res.json()
    if (!res.ok) throw new Error(data?.error || 'Upload failed')
    const returned = Array.isArray(data.items)
      ? data.items
      : data.item ? [data.item] : (Array.isArray(data.urls) ? data.urls.map(u => ({ url: u, type: 'image' })) : (data.url ? [{ url: data.url, type: 'image' }] : []))
    let toAssign = returned.length
    for (const it of mediaItems.value) {
      if (!it.url && toAssign > 0) {
        const idx = returned.length - toAssign
        it.url = returned[idx].url
        it.type = (returned[idx].type === 'video' ? 'video' : 'image')
        toAssign--
      }
    }
    if (!returned.length) {
      error.value = 'Upload failed'
      mediaItems.value = mediaItems.value.filter(m => !!m.url)
    }
  } catch (err) {
    console.error('Upload failed', err)
    error.value = 'Upload failed'
    mediaItems.value = mediaItems.value.filter(m => !!m.url)
  } finally {
    uploading.value = false
    // Clear the input value to allow selecting the same files again
    if (e && e.target) e.target.value = ''
  }
}

function removeMedia(i) {
  const m = mediaItems.value[i]
  if (m && m.preview && typeof m.preview === 'string' && m.preview.startsWith('blob:')) URL.revokeObjectURL(m.preview)
  mediaItems.value.splice(i, 1)
  if (featuredIndex.value >= mediaItems.value.length) featuredIndex.value = 0
}

function onMediaError(i) {
  if (mediaItems.value[i]) mediaItems.value[i].broken = true
}

function displaySrc(m) {
  if (!m) return ''
  // Use the original URLs without adding transformations to prevent 400/403 from strict Cloudinary settings
  if (m.type === 'video') {
    return m.url || m.preview || ''
  }
  return m.preview || m.url || ''
}

// Create recipe mutation. category_id is optional Int. user_id is set by Hasura from JWT.
const CREATE_RECIPE = gql`
  mutation(
    $title: String!,
    $description: String,
    $featured_image: String,
    $prep_time_minutes: Int,
    $category_id: Int
  ) {
    insert_recipes_one(object: {
      title: $title,
      description: $description,
      featured_image: $featured_image,
      prep_time_minutes: $prep_time_minutes,
      category_id: $category_id
    }) { id }
  }
`
const { mutate: createRecipe } = useMutation(CREATE_RECIPE)

const INSERT_INGREDIENTS = gql`
  mutation($objects: [recipe_ingredients_insert_input!]!) {
    insert_recipe_ingredients(objects: $objects) { affected_rows }
  }
`
const { mutate: insertIngredients } = useMutation(INSERT_INGREDIENTS)

const INSERT_STEPS = gql`
  mutation($objects: [recipe_steps_insert_input!]!) {
    insert_recipe_steps(objects: $objects) { affected_rows }
  }
`
const { mutate: insertSteps } = useMutation(INSERT_STEPS)

// Recipe media (images/videos)
const INSERT_MEDIA = gql`
  mutation($objects: [recipe_media_insert_input!]!) {
    insert_recipe_media(objects: $objects) { affected_rows }
  }
`
const { mutate: insertMedia } = useMutation(INSERT_MEDIA)

function addIngredient() {
  ingredients.value.push({ name: '', quantity: '' })
}
function removeIngredient(i) {
  ingredients.value.splice(i, 1)
}
function addStep() {
  steps.value.push({ description: '' })
}
function removeStep(i) {
  steps.value.splice(i, 1)
}

async function submitRecipe() {
  error.value = ''
  if (uploading.value) { error.value = 'Please wait for media uploads to finish.'; return }
  if (!title.value) { error.value = 'Title is required'; return }
  if (!userId.value) { error.value = 'You must be logged in to create recipes.'; return }
  if (!steps.value.length || !steps.value.some(s => s.description && s.description.trim().length)) {
    error.value = 'Please add at least one step.'; return
  }
  // If the user attempted to add media, ensure we have at least one uploaded URL
  const attemptedMedia = mediaItems.value.length > 0
  const uploadedCount = mediaItems.value.filter(m => !!m.url).length
  if (attemptedMedia && uploadedCount === 0) {
    error.value = 'Media upload failed. Please try uploading again.'
    return
  }
  // Use a valid uploaded URL for featured if available; otherwise null
  let featured = mediaItems.value[featuredIndex.value]?.url || ''
  if (!featured) {
    const firstWithUrl = mediaItems.value.find(m => m && m.url)
    featured = firstWithUrl?.url || ''
  }

  // Try to coerce categoryId to Int if possible; else null
  let catVar = null
  if (categoryId.value !== null && categoryId.value !== undefined) {
    const maybe = parseInt(String(categoryId.value), 10)
    catVar = Number.isFinite(maybe) ? maybe : null
  }

  try {
    const { data } = await createRecipe({
      title: title.value,
      description: description.value || null,
      featured_image: featured || null,
      prep_time_minutes: prepTime.value || null,
      category_id: catVar,
    })
    const id = data?.insert_recipes_one?.id
    if (!id) throw new Error('Recipe create failed')

    // Bulk insert ingredients
    const ingObjects = ingredients.value
      .filter(i => i.name && i.name.trim().length)
      .map(i => ({ recipe_id: id, name: i.name.trim(), quantity: i.quantity ? i.quantity.trim() : null }))
    if (ingObjects.length) {
      await insertIngredients({ objects: ingObjects })
    }

    // Bulk insert steps with positions
    const stepObjects = steps.value
      .filter(s => s.description && s.description.trim().length)
      .map((s, idx) => ({ recipe_id: id, position: idx + 1, description: s.description.trim() }))
    if (stepObjects.length) {
      await insertSteps({ objects: stepObjects })
    }

    // Insert media (images/videos) with positions if available
    if (mediaItems.value.length) {
      const mediaObjects = mediaItems.value
        .filter(m => !!m.url)
        .map((m, idx) => ({
          recipe_id: id,
          url: m.url,
          media_type: m.type,
          position: idx + 1,
        }))
      if (mediaObjects.length) {
        const mediaRes = await insertMedia({ objects: mediaObjects })
        const affected = mediaRes?.data?.insert_recipe_media?.affected_rows ?? 0
        const gqlErrors = mediaRes?.errors?.length ? mediaRes.errors : []
        if (gqlErrors.length) {
          throw new Error(gqlErrors[0]?.message || 'Failed to save media records')
        }
        if (affected < mediaObjects.length) {
          throw new Error('Failed to save all media records')
        }
      }
    }

    navigateTo('/recipes/my');
    return
  } catch (e) {
    console.error(e)
    error.value = e?.message || 'Error creating recipe.'
  }
}

</script>
