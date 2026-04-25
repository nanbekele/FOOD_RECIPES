<template>
  <div class="min-h-screen bg-gradient-to-br from-orange-50 to-white">
    <Navbar />
    <section class="container mx-auto px-4 pt-6 max-w-5xl">
      <h1 class="text-3xl font-bold mb-4">Edit Recipe</h1>
      <div v-if="loading" class="flex justify-center items-center h-40">
        <span class="loader"></span>
      </div>
      <div v-else-if="error" class="text-red-500">{{ error }}</div>
      <form v-else @submit.prevent="onSubmit" class="bg-white p-6 rounded-xl shadow space-y-6">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium">Title</label>
            <input v-model="title" type="text" class="w-full rounded border border-gray-200 px-3 py-2" required />
          </div>
          <div>
            <label class="block text-sm font-medium">Prep Time (minutes)</label>
            <input v-model.number="prepTime" type="number" min="0" class="w-full rounded border border-gray-200 px-3 py-2" />
          </div>
          <div class="md:col-span-2">
            <label class="block text-sm font-medium">Description</label>
            <textarea v-model="description" rows="4" class="w-full rounded border border-gray-200 px-3 py-2"></textarea>
          </div>
          <div>
            <label class="block text-sm font-medium">Category</label>
            <select v-model="categoryId" class="w-full rounded border border-gray-200 px-3 py-2">
              <option :value="null">Select a category</option>
              <option v-for="c in categories" :key="c.id" :value="c.id">{{ c.name }}</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium">Featured Image</label>
            <input type="file" @change="onFeaturedChange" accept="image/*" class="w-full rounded border border-gray-200 px-3 py-2" />
            <div class="mt-2">
              <img v-if="featuredPreview || featured_image" :src="featuredPreview || featured_image" class="w-40 h-28 object-cover rounded border" />
            </div>
          </div>
        </div>

        <div>
          <label class="block text-sm font-medium">Media (images or videos)</label>
          <input type="file" multiple accept="image/*,video/*,.mp4,.mov,.avi,.mkv,.webm,.m4v,.3gp,.mpeg,.mpg" @change="onMediaSelected" class="w-full rounded border border-gray-200 px-3 py-2" />
          <p class="text-xs text-gray-400 mt-1">You can upload up to 3 files (images and/or videos).</p>
          <div v-if="mediaItems.length" class="mt-3 grid grid-cols-3 gap-3">
            <div v-for="(m, i) in mediaItems" :key="(m.url || m.preview || i) + '-' + i" class="relative group">
              <div class="w-full h-24 bg-gray-100 rounded flex items-center justify-center overflow-hidden">
                <template v-if="m.type === 'video'">
                  <video v-if="(m.preview || m.url) && !m.broken" :src="displaySrc(m)" class="w-full h-full object-contain" controls playsinline preload="metadata" @error="onMediaError(i)"></video>
                  <div v-else class="text-xs text-gray-500">Preview unavailable</div>
                </template>
                <template v-else>
                  <img v-if="(m.preview || m.url) && !m.broken" :src="displaySrc(m)" class="w-full h-full object-contain" @error="onMediaError(i)" />
                  <div v-else class="text-xs text-gray-500">Preview unavailable</div>
                </template>
              </div>
              <button type="button" @click="removeMedia(i)" class="absolute -top-2 -right-2 bg-white border border-gray-300 text-gray-700 w-6 h-6 rounded-full flex items-center justify-center shadow-sm hover:bg-red-600 hover:text-white" aria-label="Remove">×</button>
            </div>
          </div>
        </div>

        <div>
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

        <div>
          <label class="block text-sm font-medium mb-2">Steps</label>
          <div class="space-y-3">
            <div v-for="(step, idx) in steps" :key="idx" class="grid grid-cols-[auto,1fr,auto] gap-3 items-start">
              <div class="w-8 h-8 rounded-full bg-orange-600 text-white flex items-center justify-center">{{ idx + 1 }}</div>
              <textarea v-model="step.description" placeholder="Step description" class="rounded border border-gray-200 px-3 py-2"></textarea>
              <button type="button" @click="removeStep(idx)" class="px-3 py-2 bg-gray-200 rounded" v-if="steps.length > 1">Remove</button>
            </div>
            <button type="button" @click="addStep" class="px-3 py-2 bg-blue-600 text-white rounded">Add Step</button>
          </div>
        </div>

        <div class="flex gap-3">
          <button type="submit" class="btn" :disabled="submitting">
            <span v-if="submitting" class="loader mr-2"></span>
            Save Changes
          </button>
          <div v-if="success" class="text-green-600 self-center">Recipe updated!</div>
          <div v-if="submitError" class="text-red-500 self-center">{{ submitError }}</div>
        </div>
      </form>
    </section>
  </div>
</template>

<script setup>
import Navbar from '~/components/Navbar.vue'
import { ref, watchEffect } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useQuery, useMutation } from '@vue/apollo-composable'
import gql from 'graphql-tag'

const route = useRoute()
const router = useRouter()
const id = route.params.id

// Form state
const title = ref('')
const description = ref('')
const featured_image = ref('')
const featuredFile = ref(null)
const featuredPreview = ref('')
const prepTime = ref(null)
const categoryId = ref(null)
const ingredients = ref([{ name: '', quantity: '' }])
const steps = ref([{ description: '' }])
const mediaItems = ref([]) // { url, preview, type }

// UI state
const loading = ref(true)
const error = ref('')
const submitting = ref(false)
const success = ref(false)
const submitError = ref('')
const backend = useRuntimeConfig().public.NUXT_PUBLIC_BACKEND_ENDPOINT

// Categories
const categories = ref([])

// Queries
const { result, onError } = useQuery(gql`
  query($id: uuid!) {
    recipes_by_pk(id: $id) {
      id title description featured_image prep_time_minutes category_id
      ingredients: recipe_ingredients { id name quantity }
      steps(order_by: {position: asc}) { id description position }
      media(order_by: {position: asc}) { id url media_type position }
    }
    categories(order_by: {name: asc}) { id name }
  }
`, { id })

onError((e) => {
  error.value = e.message
  loading.value = false
})

watchEffect(() => {
  const r = result.value
  if (!r) return
  if (r.categories) categories.value = r.categories
  const rec = r.recipes_by_pk
  if (rec) {
    title.value = rec.title
    description.value = rec.description
    featured_image.value = rec.featured_image
    prepTime.value = rec.prep_time_minutes
    categoryId.value = rec.category_id ?? null
    ingredients.value = (rec.ingredients?.length ? rec.ingredients.map(i => ({ name: i.name || '', quantity: i.quantity || '' })) : [{ name: '', quantity: '' }])
    steps.value = (rec.steps?.length ? rec.steps.map(s => ({ description: s.description || '' })) : [{ description: '' }])
    mediaItems.value = (rec.media?.length ? rec.media.map(m => ({ url: m.url, preview: '', type: m.media_type })) : [])
    loading.value = true === false ? true : false
    loading.value = false
  }
})

async function onFeaturedChange(e) {
  const file = e.target.files?.[0]
  if (!file) return
  featuredFile.value = file
  featuredPreview.value = URL.createObjectURL(file)
}

function displaySrc(m) {
  const src = m.preview || m.url || ''
  if (!src) return ''
  const isCloud = src.includes('res.cloudinary.com') && src.includes('/upload/')
  if (!isCloud) return src
  const trans = m.type === 'video' ? 'f_auto,q_auto,vc_auto' : 'f_auto,q_auto'
  return src.replace('/upload/', `/upload/${trans}/`)
}
function onMediaError(i) { if (mediaItems.value[i]) mediaItems.value[i].broken = true }

function addIngredient() { ingredients.value.push({ name: '', quantity: '' }) }
function removeIngredient(i) { ingredients.value.splice(i, 1) }
function addStep() { steps.value.push({ description: '' }) }
function removeStep(i) { steps.value.splice(i, 1) }

function removeMedia(i) {
  const m = mediaItems.value[i]
  if (m && m.preview && typeof m.preview === 'string' && m.preview.startsWith('blob:')) URL.revokeObjectURL(m.preview)
  mediaItems.value.splice(i, 1)
}

async function onMediaSelected(e) {
  const selected = Array.from(e.target.files || [])
  if (!selected.length) return
  const remaining = Math.max(0, 3 - mediaItems.value.length)
  if (remaining <= 0) { alert('You can upload up to 3 files only.'); return }
  const files = selected.slice(0, remaining)

  // Local previews
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
      resolve({ type: 'video', preview: url, file })
    }
  })))
  for (const p of previews) mediaItems.value.push({ url: '', preview: p.preview, type: p.type })

  // Upload
  const form = new FormData()
  for (const f of files) form.append('files', f)
  try {
    const res = await fetch(`${backend}/upload`, { method: 'POST', body: form })
    const data = await res.json()
    const returned = Array.isArray(data.items) ? data.items : (data.url ? [{ url: data.url, type: data.type || 'image' }] : [])
    let toAssign = Math.min(returned.length, previews.length)
    for (const it of mediaItems.value) {
      if (!it.url && toAssign > 0) {
        const idx = returned.length - toAssign
        it.url = returned[idx].url
        it.type = (returned[idx].type === 'video' ? 'video' : 'image')
        toAssign--
      }
    }
  } catch (err) {
    console.error('Upload failed', err)
  }
  if (e && e.target) e.target.value = ''
}

// Mutations
const UPDATE_RECIPE = gql`
  mutation UpdateRecipe($id: uuid!, $title: String!, $description: String, $featured_image: String, $prep: Int, $cat: Int) {
    update_recipes_by_pk(pk_columns: {id: $id}, _set: { title: $title, description: $description, featured_image: $featured_image, prep_time_minutes: $prep, category_id: $cat }) { id }
  }
`
const DELETE_INGS = gql`mutation($id: uuid!) { delete_recipe_ingredients(where: { recipe_id: { _eq: $id } }) { affected_rows } }`
const INSERT_INGS = gql`mutation($objects: [recipe_ingredients_insert_input!]!) { insert_recipe_ingredients(objects: $objects) { affected_rows } }`
const DELETE_STEPS = gql`mutation($id: uuid!) { delete_recipe_steps(where: { recipe_id: { _eq: $id } }) { affected_rows } }`
const INSERT_STEPS = gql`mutation($objects: [recipe_steps_insert_input!]!) { insert_recipe_steps(objects: $objects) { affected_rows } }`
const DELETE_MEDIA = gql`mutation($id: uuid!) { delete_recipe_media(where: { recipe_id: { _eq: $id } }) { affected_rows } }`
const INSERT_MEDIA = gql`mutation($objects: [recipe_media_insert_input!]!) { insert_recipe_media(objects: $objects) { affected_rows } }`

const { mutate: doUpdate } = useMutation(UPDATE_RECIPE)
const { mutate: delIngs } = useMutation(DELETE_INGS)
const { mutate: insIngs } = useMutation(INSERT_INGS)
const { mutate: delSteps } = useMutation(DELETE_STEPS)
const { mutate: insSteps } = useMutation(INSERT_STEPS)
const { mutate: delMedia } = useMutation(DELETE_MEDIA)
const { mutate: insMedia } = useMutation(INSERT_MEDIA)

async function onSubmit() {
  submitting.value = true
  submitError.value = ''

  // Upload featured if changed
  let imageUrl = featured_image.value
  if (featuredFile.value) {
    const fd = new FormData()
    fd.append('image', featuredFile.value)
    try {
      const res = await fetch(`${backend}/upload`, { method: 'POST', body: fd })
      const data = await res.json()
      if (res.ok && data.url) imageUrl = data.url
    } catch {}
  }

  try {
    const { data } = await doUpdate({ id, title: title.value, description: description.value || null, featured_image: imageUrl || null, prep: prepTime.value || null, cat: categoryId.value ?? null })
    if (!data?.update_recipes_by_pk?.id) throw new Error('Update failed')

    // Replace ingredients
    await delIngs({ id })
    const ingObjs = ingredients.value
      .filter(i => i.name && i.name.trim().length)
      .map(i => ({ recipe_id: id, name: i.name.trim(), quantity: i.quantity ? i.quantity.trim() : null }))
    if (ingObjs.length) await insIngs({ objects: ingObjs })

    // Replace steps
    await delSteps({ id })
    const stepObjs = steps.value
      .filter(s => s.description && s.description.trim().length)
      .map((s, idx) => ({ recipe_id: id, position: idx + 1, description: s.description.trim() }))
    if (stepObjs.length) await insSteps({ objects: stepObjs })

    // Replace media
    await delMedia({ id })
    if (mediaItems.value.length) {
      const mediaObjs = mediaItems.value.map((m, idx) => ({ recipe_id: id, url: m.url, media_type: m.type, position: idx + 1 }))
      await insMedia({ objects: mediaObjs })
    }

    success.value = true
    setTimeout(() => router.push(`/recipes/${id}`), 800)
  } catch (err) {
    console.error(err)
    submitError.value = 'Update failed. Ensure Hasura permissions and schema.'
  }

  submitting.value = false
}
</script>

<style scoped>
.btn { @apply bg-blue-600 text-white py-2 px-4 rounded hover:bg-blue-700 transition; }
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
@keyframes spin { 0% { transform: rotate(0deg);} 100% { transform: rotate(360deg);} }
</style>



