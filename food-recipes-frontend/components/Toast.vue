<template>
  <transition name="fade">
    <div v-if="show" :class="['fixed top-6 right-6 px-4 py-2 rounded shadow-lg z-50', typeClass]">
      {{ message }}
    </div>
  </transition>
</template>

<script setup>
import { ref } from 'vue'
const show = ref(false)
const message = ref('')
const typeClass = ref('bg-green-500 text-white')

function showToast(msg, type = 'success') {
  message.value = msg
  typeClass.value = type === 'success'
    ? 'bg-green-500 text-white'
    : 'bg-red-500 text-white'
  show.value = true
  setTimeout(() => { show.value = false }, 3000)
}

defineExpose({ showToast })
</script>

<style scoped>
.fade-enter-active, .fade-leave-active { transition: opacity 0.3s; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
