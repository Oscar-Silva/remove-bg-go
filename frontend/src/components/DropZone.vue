<script setup lang="ts">
import { ref } from 'vue'
import { Upload } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'
import { fileToBase64 } from '@/lib/utils'
import { RemoveBackground } from '@wailsjs/go/main/App'
import { EventsOn } from '@wailsjs/runtime'

const store = useAppStore()
const isDragging = ref(false)
const fileInput = ref<HTMLInputElement | null>(null)

const ACCEPTED_TYPES = ['image/png', 'image/jpeg', 'image/webp']
const MAX_SIZE = 20 * 1024 * 1024 // 20MB

// Listen for status updates from Go
EventsOn('status', (status: string) => {
  store.setStatusMessage(status)
  if (status === 'done') {
    // Result will be set after RemoveBackground returns
  } else if (status === 'error') {
    store.setError('Processing failed')
  }
})

async function handleFile(file: File) {
  if (!ACCEPTED_TYPES.includes(file.type)) {
    store.setError('Please upload a PNG, JPG, or WEBP image')
    return
  }

  if (file.size > MAX_SIZE) {
    store.setError('Image must be smaller than 20MB')
    return
  }

  try {
    const base64 = await fileToBase64(file)
    // Remove data URL prefix if present
    const imageData = base64.replace(/^data:image\/\w+;base64,/, '')
    store.setOriginalImage(base64)
    store.setLoading('Processing image...')
    
    // Call the Go backend
    const result = await RemoveBackground(imageData)
    store.setResultImage(result)
    store.setDone()
  } catch (e) {
    const errorMsg = e instanceof Error ? e.message : String(e)
    store.setError('Failed to process image: ' + errorMsg)
  }
}

function handleDrop(e: DragEvent) {
  isDragging.value = false
  const file = e.dataTransfer?.files[0]
  if (file) {
    handleFile(file)
  }
}

function handleDragOver(e: DragEvent) {
  e.preventDefault()
  isDragging.value = true
}

function handleDragLeave() {
  isDragging.value = false
}

function handleClick() {
  fileInput.value?.click()
}

function handleFileInput(e: Event) {
  const target = e.target as HTMLInputElement
  const file = target.files?.[0]
  if (file) {
    handleFile(file)
  }
}
</script>

<template>
  <div
    class="relative w-full max-w-xl mx-auto"
    @drop.prevent="handleDrop"
    @dragover="handleDragOver"
    @dragleave="handleDragLeave"
  >
    <!-- Drop Area -->
    <div
      @click="handleClick"
      :class="[
        'cursor-crosshair',
        'relative',
        'border-2',
        'border-dashed',
        'rounded-2xl',
        'p-12',
        'flex',
        'flex-col',
        'items-center',
        'justify-center',
        'gap-6',
        'transition-all',
        'duration-200',
        isDragging
          ? 'border-accent bg-accent/10 shadow-glow'
          : 'border-border hover:border-accent/50 hover:bg-surfaceElevated/50'
      ]"
    >
      <!-- Animated Icon -->
      <div
        :class="[
          'w-20 h-20',
          'rounded-2xl',
          'bg-surfaceElevated',
          'flex',
          'items-center',
          'justify-center',
          'transition-transform',
          'duration-300',
          isDragging ? 'scale-110' : 'hover:scale-105',
          isDragging ? 'animate-pulse' : 'animate-float'
        ]"
      >
        <Upload
          :class="[
            'w-10 h-10',
            'transition-colors',
            'duration-200',
            isDragging ? 'text-accent' : 'text-textMuted'
          ]"
        />
      </div>

      <!-- Text -->
      <div class="text-center">
        <p class="text-lg font-medium text-textPrimary mb-1">
          Drop your image here
        </p>
        <p class="text-sm text-textMuted">
          PNG, JPG, WEBP up to 20MB
        </p>
      </div>

      <!-- Hidden Input -->
      <input
        ref="fileInput"
        type="file"
        :accept="ACCEPTED_TYPES.join(',')"
        class="hidden"
        @change="handleFileInput"
      />
    </div>

    <!-- Border Animation -->
    <div
      v-if="!isDragging"
      class="absolute inset-0 rounded-2xl pointer-events-none"
      style="
        background: linear-gradient(90deg, transparent 50%, rgba(124, 110, 250, 0.1) 50%);
        background-size: 200% 100%;
        animation: shimmer 3s linear infinite;
      "
    />
  </div>
</template>

<style scoped>
@keyframes shimmer {
  0% {
    background-position: 200% 0;
  }
  100% {
    background-position: -200% 0;
  }
}
</style>
