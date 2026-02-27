<script setup lang="ts">
import { useAppStore } from '@/stores/app'
import { ref, onMounted, onUnmounted, watch } from 'vue'

const store = useAppStore()

const currentProgress = ref(0)
const timeRemaining = ref(0)
let progressInterval: ReturnType<typeof setInterval>

const EXPECTED_PROCESS_TIME_MS = 6000 // ~6 seconds for CPU processing on average

const statusMessages: Record<string, string> = {
  initializing: 'Initializing...',
  downloading_model: 'Downloading AI model...',
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

watch(() => store.statusMessage, (newStatus) => {
  if (progressInterval) clearInterval(progressInterval)
  
  const statusMap: Record<string, number> = {
    initializing: 2,
    downloading_model: 5, // We'll compute real progress in template
    loading_model: 10,
    loading: 15,
    decoding: 20,
    preprocessing: 25,
    processing: 30, // Start heavy ML phase at 30%
    finalizing: 90,
    done: 100,
  }
  
  const baseProgress = statusMap[newStatus.toLowerCase()] || 0
  
  if (newStatus.toLowerCase() === 'processing') {
    currentProgress.value = baseProgress
    timeRemaining.value = EXPECTED_PROCESS_TIME_MS / 1000
    
    const updateRate = 100
    // animate from 30 to 90
    const progressPerStep = (60 / (EXPECTED_PROCESS_TIME_MS / updateRate)) 
    
    progressInterval = setInterval(() => {
      if (currentProgress.value < 90) {
        currentProgress.value = Math.min(90, currentProgress.value + progressPerStep)
      }
      if (timeRemaining.value > 0) {
        timeRemaining.value = Math.max(0, timeRemaining.value - (updateRate / 1000))
      }
    }, updateRate)
  } else {
    currentProgress.value = baseProgress
    if (newStatus.toLowerCase() === 'done' || newStatus.toLowerCase() === 'error') {
      timeRemaining.value = 0
    }
  }
}, { immediate: true })

onUnmounted(() => {
  if (progressInterval) clearInterval(progressInterval)
})

</script>

<template>
  <div class="w-full max-w-xl mx-auto">
    <!-- Preview Image with Shimmer -->
    <div class="relative mb-8 group">
      <!-- Glow effect behind container -->
      <div class="absolute -inset-1 bg-gradient-to-r from-accent to-accent-hover rounded-3xl blur opacity-20 group-hover:opacity-40 transition duration-1000 group-hover:duration-200"></div>
      
      <div
        class="w-full aspect-square rounded-2xl overflow-hidden bg-surfaceElevated/50 backdrop-blur-md relative border border-white/5 shadow-2xl"
      >
        <img
          v-if="store.originalImageUrl"
          :src="store.originalImageUrl"
          alt="Processing"
          class="w-full h-full object-contain mix-blend-screen"
        />

        <!-- Shimmer Overlay -->
        <div
          class="absolute inset-0 bg-gradient-to-r from-transparent via-white/5 to-transparent animate-shimmer"
        />
      </div>
    </div>

    <!-- Status Text -->
    <div class="text-center mb-6">
      <h3 class="text-lg font-medium text-textPrimary mb-2">
        <span v-if="store.statusMessage === 'downloading_model'">
          {{ getStatusText() }} ({{ store.downloadProgress.total > 0 ? Math.round((store.downloadProgress.downloaded / store.downloadProgress.total) * 100) : 0 }}%)
        </span>
        <span v-else>
          {{ getStatusText() }} ({{ Math.round(currentProgress) }}%)
        </span>
      </h3>
      <p class="text-sm text-textMuted h-5">
        <span v-if="store.statusMessage === 'downloading_model' && store.downloadProgress.total > 0">
          {{ (store.downloadProgress.downloaded / 1024 / 1024).toFixed(1) }} MB / {{ (store.downloadProgress.total / 1024 / 1024).toFixed(1) }} MB
        </span>
        <span v-else-if="timeRemaining > 0">
          Estimated time remaining: ~{{ Math.ceil(timeRemaining) }}s
        </span>
        <span v-else-if="currentProgress < 100">
          Please wait...
        </span>
      </p>
    </div>

    <!-- Progress Bar -->
    <div class="w-full h-3 bg-surfaceElevated/50 backdrop-blur rounded-full overflow-hidden shadow-inner relative border border-white/5">
      <!-- Standard Progress -->
      <div
        v-if="store.statusMessage !== 'downloading_model'"
        class="h-full bg-gradient-to-r from-accent via-accent-hover to-accent rounded-full transition-all duration-300 ease-out absolute left-0 top-0 bg-[length:200%_auto] animate-gradient"
        :style="{ width: `${currentProgress}%` }"
      />
      <!-- Shimmer effect on the filled part -->
      <div
        v-if="store.statusMessage !== 'downloading_model'"
        class="h-full bg-gradient-to-r from-transparent via-white/30 to-transparent absolute top-0 animate-progress"
        :style="{ width: `${currentProgress}%`, left: 0 }"
      />
      
      <!-- Download Progress -->
      <div
        v-if="store.statusMessage === 'downloading_model'"
        class="h-full bg-gradient-to-r from-success to-emerald-400 rounded-full transition-all duration-100 ease-out absolute left-0 top-0 shadow-[0_0_15px_rgba(52,211,153,0.5)]"
        :style="{ width: `${store.downloadProgress.total > 0 ? (store.downloadProgress.downloaded / store.downloadProgress.total) * 100 : 0}%` }"
      />
    </div>

    <!-- Cancel Button -->
    <div class="mt-6 text-center">
      <button
        @click="store.reset()"
        class="px-6 py-2.5 rounded-xl border border-white/5 bg-white/[0.02] backdrop-blur text-sm text-textMuted hover:text-white hover:bg-white/[0.05] hover:border-white/10 transition-all duration-300 shadow-sm"
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
  animation: progress 2s linear infinite;
}

@keyframes gradient {
  0% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
  100% { background-position: 0% 50%; }
}

.animate-gradient {
  animation: gradient 3s ease infinite;
}
</style>
