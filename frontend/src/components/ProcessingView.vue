<script setup lang="ts">
import { useAppStore } from '@/stores/app'

const store = useAppStore()

const statusMessages: Record<string, string> = {
  initializing: 'Initializing...',
  loading_model: 'Loading model...',
  loading: 'Loading image...',
  decoding: 'Decoding image...',
  preprocessing: 'Preparing image...',
  processing: 'Removing background...',
  finalizing: 'Finalizing...',
  done: 'Complete!',
  error: 'Error occurred',
}

function getStatusText(): string {
  return statusMessages[store.statusMessage.toLowerCase()] || store.statusMessage || 'Processing...'
}
</script>

<template>
  <div class="w-full max-w-xl mx-auto">
    <!-- Preview Image with Shimmer -->
    <div class="relative mb-8">
      <div
        class="w-full aspect-square rounded-2xl overflow-hidden bg-surfaceElevated relative"
      >
        <img
          v-if="store.originalImageUrl"
          :src="store.originalImageUrl"
          alt="Processing"
          class="w-full h-full object-contain"
        />

        <!-- Shimmer Overlay -->
        <div
          class="absolute inset-0 bg-gradient-to-r from-transparent via-white/10 to-transparent animate-shimmer"
        />
      </div>
    </div>

    <!-- Status Text -->
    <div class="text-center mb-6">
      <h3 class="text-lg font-medium text-textPrimary mb-2">
        {{ getStatusText() }}
      </h3>
      <p class="text-sm text-textMuted">
        Please wait while we process your image
      </p>
    </div>

    <!-- Progress Bar -->
    <div class="w-full h-1 bg-surfaceElevated rounded-full overflow-hidden">
      <div
        class="h-full bg-gradient-to-r from-accent to-accent-hover rounded-full animate-progress"
        style="width: 100%"
      />
    </div>

    <!-- Cancel Button -->
    <div class="mt-6 text-center">
      <button
        @click="store.reset()"
        class="px-4 py-2 text-sm text-textMuted hover:text-textPrimary transition-colors duration-200"
      >
        Cancel
      </button>
    </div>
  </div>
</template>

<style scoped>
@keyframes shimmer {
  0% {
    transform: translateX(-100%);
  }
  100% {
    transform: translateX(100%);
  }
}

.animate-shimmer {
  animation: shimmer 1.5s ease-in-out infinite;
}

@keyframes progress {
  0% {
    background-position: 0% 50%;
  }
  100% {
    background-position: 100% 50%;
  }
}

.animate-progress {
  background-size: 200% 100%;
  animation: progress 1.5s ease-in-out infinite;
}
</style>
