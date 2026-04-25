<template>
  <transition name="fade">
    <div class="fixed inset-0 bg-black/40 flex items-center justify-center z-50" @click.self="$emit('close')">
      <div class="bg-white p-6 rounded-lg shadow-xl w-full max-w-md">
        <header class="mb-4">
          <slot name="title" />
        </header>
        <div class="mb-6">
          <slot name="body" />
        </div>
        <footer class="flex justify-end gap-2">
          <slot name="footer" />
        </footer>
      </div>
    </div>
  </transition>
  
</template>

<script setup>
</script>

<style scoped>
.fade-enter-active, .fade-leave-active { transition: opacity 0.2s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
