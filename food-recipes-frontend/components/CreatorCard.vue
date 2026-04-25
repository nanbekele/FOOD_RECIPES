<template>
  <div class="p-4 bg-white rounded-xl border border-gray-200 shadow-sm hover:shadow-md transition">
    <div class="flex items-center gap-3">
      <div class="h-10 w-10 rounded-full bg-orange-200 text-orange-800 flex items-center justify-center font-bold">
        {{ initials }}
      </div>
      <div>
        <p class="font-semibold text-gray-900 flex items-center gap-2">
          <span>{{ creator.name }}</span>
          <span v-if="verified" class="inline-flex items-center text-xs px-2 py-0.5 rounded bg-green-100 text-green-700">Verified</span>
        </p>
        <p class="text-sm text-gray-500">Recipe creator</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
const props = defineProps({
  creator: { type: Object, required: true }
})

const initials = computed(() => {
  const n = props.creator?.name || 'U'
  return n.split(' ').map(w => w[0]).slice(0,2).join('').toUpperCase()
})
const verified = computed(() => (props.creator?.role === 'chef') && !!props.creator?.is_verified)
</script>