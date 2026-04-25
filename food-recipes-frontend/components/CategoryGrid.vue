<template>
  <div class="flex flex-wrap gap-2">
    <button
      v-for="c in categories"
      :key="c.id"
      @click="$emit('select', c.slug || c.name)"
      class="px-3 py-2 rounded-full bg-white border border-orange-200 text-orange-700 hover:bg-orange-50 shadow-sm"
    >
      {{ c.name }}
    </button>
  </div>
</template>

<script setup>
const props = defineProps({
  categories: { type: Array, default: () => [] }
})
</script>