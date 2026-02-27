<script setup lang="ts">
import { ref, computed } from 'vue'
import { Download, Plus, Copy, Check } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'

const store = useAppStore()
const sliderPosition = ref(50)
const isDragging = ref(false)
const copied = ref(false)

const containerRef = ref<HTMLElement | null>(null)

function startDrag(e: MouseEvent | TouchEvent) {
  isDragging.value = true
  updatePosition(e)
}

function onDrag(e: MouseEvent | TouchEvent) {
  if (isDragging.value) {
    updatePosition(e)
  }
}

function stopDrag() {
  isDragging.value = false
}

function updatePosition(e: MouseEvent | TouchEvent) {
  if (!containerRef.value) return

  const rect = containerRef.value.getBoundingClientRect()
  const clientX = 'touches' in e ? e.touches[0].clientX : e.clientX
  const x = clientX - rect.left
  const percentage = (x / rect.width) * 100
  sliderPosition.value = Math.max(0, Math.min(100, percentage))
}

async function downloadImage() {
  if (!store.resultImageUrl) return

  const link = document.createElement('a')
  link.href = store.resultImageUrl
  link.download = `removed-bg-${Date.now()}.png`
  link.click()
}

async function copyToClipboard() {
  if (!store.resultImageUrl) return

  try {
    const response = await fetch(store.resultImageUrl)
    const blob = await response.blob()
    await navigator.clipboard.write([
      new ClipboardItem({ [blob.type]: blob })
    ])
    copied.value = true
    setTimeout(() => {
      copied.value = false
    }, 2000)
  } catch (e) {
    console.error('Failed to copy:', e)
  }
}

function removeAnother() {
  store.reset()
}
</script>

<template>
  <div class="w-full">
    <!-- Comparison View -->
    <div
      ref="containerRef"
      class="relative w-full aspect-square rounded-2xl overflow-hidden bg-surface cursor-grab active:cursor-grabbing select-none"
      @mousedown="startDrag"
      @mousemove="onDrag"
      @mouseup="stopDrag"
      @mouseleave="stopDrag"
      @touchstart="startDrag"
      @touchmove="onDrag"
      @touchend="stopDrag"
    >
      <!-- Checkerboard Background -->
      <div
        class="absolute inset-0"
        style="
          background-image:
            linear-gradient(45deg, #1a1a2e 25%, transparent 25%),
            linear-gradient(-45deg, #1a1a2e 25%, transparent 25%),
            linear-gradient(45deg, transparent 75%, #1a1a2e 75%),
            linear-gradient(-45deg, transparent 75%, #1a1a2e 75%);
          background-size: 20px 20px;
          background-position: 0 0, 0 10px, 10px -10px, -10px 0px;
        "
      />

      <!-- Result Image (Background) -->
      <img
        v-if="store.resultImageUrl"
        :src="store.resultImageUrl"
        alt="Result"
        class="absolute inset-0 w-full h-full object-contain"
      />

      <!-- Original Image (Clipped) -->
      <div
        class="absolute inset-0 overflow-hidden"
        :style="{ clipPath: `inset(0 ${100 - sliderPosition}% 0 0)` }"
      >
        <img
          v-if="store.originalImageUrl"
          :src="store.originalImageUrl"
          alt="Original"
          class="absolute inset-0 w-full h-full object-contain"
        />
      </div>

      <!-- Slider Handle -->
      <div
        class="absolute top-0 bottom-0 w-1 bg-white cursor-ew-resize shadow-lg"
        :style="{ left: `${sliderPosition}%` }"
      >
        <!-- Handle Button -->
        <div
          class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-10 h-10 rounded-full bg-white shadow-xl flex items-center justify-center"
        >
          <div class="flex gap-0.5">
            <div class="w-0.5 h-3 bg-gray-400 rounded-full" />
            <div class="w-0.5 h-3 bg-gray-400 rounded-full" />
          </div>
        </div>
      </div>

      <!-- Labels -->
      <div class="absolute top-4 left-4 px-3 py-1.5 rounded-lg bg-black/50 text-xs text-white font-medium backdrop-blur-sm">
        Original
      </div>
      <div class="absolute top-4 right-4 px-3 py-1.5 rounded-lg bg-accent/80 text-xs text-white font-medium backdrop-blur-sm">
        Result
      </div>
    </div>

    <!-- Action Buttons -->
    <div class="flex items-center justify-center gap-4 mt-8">
      <!-- Download -->
      <button
        @click="downloadImage"
        class="flex items-center gap-2 px-6 py-3 bg-accent hover:bg-accent-hover text-white font-medium rounded-xl transition-all duration-200 hover:scale-[0.97] active:scale-[0.95]"
      >
        <Download class="w-5 h-5" />
        Download PNG
      </button>

      <!-- Remove Another -->
      <button
        @click="removeAnother"
        class="flex items-center gap-2 px-6 py-3 bg-surfaceElevated hover:bg-surface text-textMuted hover:text-textPrimary font-medium rounded-xl border border-border transition-all duration-200 hover:scale-[0.97] active:scale-[0.95]"
      >
        <Plus class="w-5 h-5" />
        Remove another
      </button>

      <!-- Copy to Clipboard -->
      <button
        @click="copyToClipboard"
        class="flex items-center gap-2 px-4 py-3 bg-surfaceElevated hover:bg-surface text-textMuted hover:text-textPrimary font-medium rounded-xl border border-border transition-all duration-200 hover:scale-[0.97] active:scale-[0.95]"
        :title="copied ? 'Copied!' : 'Copy to clipboard'"
      >
        <Check v-if="copied" class="w-5 h-5 text-success" />
        <Copy v-else class="w-5 h-5" />
      </button>
    </div>
  </div>
</template>
