<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAppStore } from '@/stores/app'
import { GetModels } from '@wailsjs/go/main/App'
import type { main } from '@wailsjs/go/models'

const store = useAppStore()
const availableModels = ref<main.ModelConfig[]>([])
const loading = ref(true)

onMounted(async () => {
  try {
    const models = await GetModels()
    availableModels.value = models
    
    // Auto-select default if none selected
    if (!store.selectedModel) {
      const defaultModel = models.find(m => m.isDefault)
      if (defaultModel) {
        store.setSelectedModel(defaultModel.id)
      }
    }
  } catch (e) {
    console.error("Failed to load models:", e)
  } finally {
    loading.value = false
  }
})

function selectModel(id: string) {
  store.setSelectedModel(id)
}
</script>

<template>
  <div class="w-full max-w-xl mx-auto flex flex-col gap-4">
    <div class="flex items-center justify-between px-2">
      <h3 class="text-textPrimary font-semibold text-sm tracking-wide bg-gradient-to-r from-accent to-accent-hover bg-clip-text text-transparent">AI Engine</h3>
      <div v-if="loading" class="text-textMuted text-xs flex items-center gap-2">
        <svg class="w-3 h-3 animate-spin" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 2v4m0 12v4M4.93 4.93l2.83 2.83m8.48 8.48l2.83 2.83M2 12h4m12 0h4M4.93 19.07l2.83-2.83m8.48-8.48l2.83-2.83"/></svg>
        Waking up...
      </div>
    </div>
    
    <div class="grid grid-cols-1 md:grid-cols-2 gap-3" v-if="!loading">
      <button
        v-for="model in availableModels"
        :key="model.id"
        @click="selectModel(model.id)"
        :class="[
          'relative text-left p-4 rounded-2xl border transition-all duration-300 group overflow-hidden',
          store.selectedModel === model.id
            ? 'border-accent bg-accent/10 shadow-[0_0_20px_rgba(108,92,231,0.2)] transform scale-[1.02]'
            : 'border-white/5 bg-white/[0.02] backdrop-blur hover:border-accent/40 hover:bg-white/[0.04] hover:shadow-lg'
        ]"
      >
        <!-- Background Glow -->
        <div 
          class="absolute inset-0 bg-gradient-to-br from-accent/20 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-500 pointer-events-none"
          v-if="store.selectedModel !== model.id"
        />

        <!-- Selection Indicator -->
        <div 
          class="absolute top-4 right-4 w-2 h-2 rounded-full transition-all duration-300 shadow-[0_0_10px_currentColor]"
          :class="store.selectedModel === model.id ? 'bg-accent text-accent' : 'bg-surfaceElevated border border-border opacity-0 group-hover:opacity-50'"
        />

        <div class="flex flex-col gap-1 relative z-10">
          <span class="text-textPrimary font-semibold text-sm flex items-center gap-2">
            {{ model.name }}
            <span v-if="model.isDefault" class="px-2 py-0.5 rounded-full bg-accent/20 text-[9px] font-bold text-accent border border-accent/30 tracking-wider">RECOMMENDED</span>
          </span>
          
          <div class="flex items-center gap-3 mt-1 text-[11px] text-textMuted font-medium">
            <span class="flex items-center gap-1.5" title="Size">
              <svg class="w-3.5 h-3.5 opacity-70" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 15a4 4 0 004 4h9a5 5 0 10-.1-9.999 5.002 5.002 0 10-9.78 2.096A4.001 4.001 0 003 15z" /></svg>
              {{ model.sizeMB }} MB
            </span>
            <span class="flex items-center gap-1.5" title="RAM">
              <svg class="w-3.5 h-3.5 opacity-70" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z" /></svg>
              {{ model.ram }}
            </span>
          </div>

          <div class="flex items-center justify-between mt-3 pt-3 border-t border-white/5">
            <span class="text-[11px] font-medium transition-colors" :class="store.selectedModel === model.id ? 'text-accent' : 'text-textMuted group-hover:text-textPrimary'">
              <span class="opacity-50">Speed:</span> {{ model.speed }}
            </span>
            <span class="text-[11px] font-medium transition-colors" :class="store.selectedModel === model.id ? 'text-success drop-shadow-[0_0_8px_rgba(52,211,153,0.5)]' : 'text-textMuted group-hover:text-textPrimary'">
              <span class="opacity-50">Quality:</span> {{ model.quality }}
            </span>
          </div>
        </div>
      </button>
    </div>
  </div>
</template>
