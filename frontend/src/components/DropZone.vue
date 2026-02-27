<script setup lang="ts">
import { ref } from 'vue'
import { Upload } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'
import { fileToBase64 } from '@/lib/utils'
import { RemoveBackground } from '@wailsjs/go/main/App'
import { EventsOn } from '@wailsjs/runtime'
import ModelSelector from './ModelSelector.vue'

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

EventsOn('download_progress', (data: { downloaded: number; total: number }) => {
  store.setDownloadProgress(data.downloaded, data.total)
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
    
    // Call the Go backend with selected model
    const result = await RemoveBackground(imageData, store.selectedModel)
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
        'rounded-3xl',
        'p-16',
        'flex',
        'flex-col',
        'items-center',
        'justify-center',
        'gap-6',
        'transition-all',
        'duration-300',
        'border border-white/5',
        'bg-white/[0.02] backdrop-blur-xl',
        isDragging
          ? 'scale-[1.02] shadow-[0_0_40px_-10px_rgba(108,92,231,0.5)] border-accent/50 bg-accent/[0.05]'
          : 'hover:border-accent/30 hover:bg-white/[0.04] hover:shadow-[0_8px_30px_rgb(0,0,0,0.12)]'
      ]"
    >
      <!-- Animated Icon -->
      <div
        :class="[
          'w-24 h-24',
          'rounded-2xl',
          'flex',
          'items-center',
          'justify-center',
          'transition-all',
          'duration-500',
          'z-10',
          'bg-gradient-to-tr from-surfaceElevated to-surface',
          'border border-white/5 shadow-xl',
          isDragging ? 'scale-110 shadow-[0_0_30px_rgba(108,92,231,0.3)]' : 'hover:scale-105',
          isDragging ? 'animate-pulse' : 'animate-float'
        ]"
      >
        <Upload
          :class="[
            'w-10 h-10',
            'transition-colors',
            'duration-300',
            isDragging ? 'text-accent drop-shadow-[0_0_8px_rgba(108,92,231,0.8)]' : 'text-textMuted'
          ]"
        />
      </div>

      <!-- Text -->
      <div class="text-center z-10">
        <p class="text-xl font-semibold bg-gradient-to-r from-white to-white/70 bg-clip-text text-transparent mb-2">
          Drop your image here
        </p>
        <p class="text-sm text-textMuted font-medium">
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
      class="absolute inset-0 rounded-3xl pointer-events-none p-[1px]"
      style="
        background: linear-gradient(90deg, transparent 50%, rgba(108, 92, 231, 0.4) 50%);
        background-size: 200% 100%;
        animation: shimmer 4s linear infinite;
        -webkit-mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
        mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
        -webkit-mask-composite: xor;
        mask-composite: exclude;
      "
    />
  </div>
  
  <ModelSelector v-if="!isDragging" class="mt-8" />
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
