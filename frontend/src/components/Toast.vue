<script setup lang="ts">
import { useAppStore } from '@/stores/app'
import { X, AlertCircle } from 'lucide-vue-next'

const store = useAppStore()
</script>

<template>
  <Transition name="toast">
    <div
      v-if="store.isError && store.error"
      class="fixed bottom-6 right-6 z-50 max-w-sm"
    >
      <div
        class="flex items-start gap-3 p-4 bg-surface border border-error/30 rounded-xl shadow-lg backdrop-blur-sm"
      >
        <!-- Icon -->
        <div class="flex-shrink-0 w-8 h-8 rounded-lg bg-error/10 flex items-center justify-center">
          <AlertCircle class="w-5 h-5 text-error" />
        </div>

        <!-- Content -->
        <div class="flex-1 min-w-0">
          <p class="text-sm font-medium text-textPrimary">
            Error
          </p>
          <p class="text-sm text-textMuted mt-0.5">
            {{ store.error }}
          </p>
        </div>

        <!-- Close Button -->
        <button
          @click="store.setIdle()"
          class="flex-shrink-0 p-1 rounded-lg text-textMuted hover:text-textPrimary hover:bg-surfaceElevated transition-colors"
        >
          <X class="w-4 h-4" />
        </button>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}

.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translateY(20px) scale(0.95);
}
</style>
