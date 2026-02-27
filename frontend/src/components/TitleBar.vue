<script setup lang="ts">
// @ts-ignore
import { WindowMinimise, WindowToggleMaximise, Quit } from '../../wailsjs/runtime/runtime'

function minimize() {
  WindowMinimise()
}

function toggleMaximize() {
  WindowToggleMaximise()
}

function close() {
  Quit()
}
</script>

<template>
  <div
    class="h-12 bg-transparent flex items-center justify-between px-4 select-none backdrop-blur-sm border-b border-white/5"
    style="-webkit-app-region: drag; --wails-draggable: drag;"
  >
    <!-- Logo and Title -->
    <div class="flex items-center gap-2">
      <img src="../assets/images/remove-bg-go-icon.svg" class="w-6 h-6 object-contain drop-shadow-[0_0_8px_rgba(108,92,231,0.5)]" alt="Remove Background Logo" />
      <span class="text-sm font-medium text-textPrimary">Remove Background</span>
    </div>

    <!-- Window Controls -->
    <div
      class="flex items-center gap-1"
      style="-webkit-app-region: no-drag; --wails-draggable: no-drag;"
    >
      <!-- Minimize -->
      <button
        @click="minimize"
        class="w-8 h-8 rounded-lg flex items-center justify-center text-textMuted hover:text-textPrimary hover:bg-white/10 transition-all duration-200"
      >
        <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M5 12h14" />
        </svg>
      </button>

      <!-- Maximize -->
      <button
        @click="toggleMaximize"
        class="w-8 h-8 rounded-lg flex items-center justify-center text-textMuted hover:text-textPrimary hover:bg-white/10 transition-all duration-200"
      >
        <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <rect x="3" y="3" width="18" height="18" rx="2" />
        </svg>
      </button>

      <!-- Close -->
      <button
        @click="close"
        class="w-8 h-8 rounded-lg flex items-center justify-center text-textMuted hover:text-white hover:bg-error transition-all duration-150"
      >
        <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M18 6L6 18M6 6l12 12" />
        </svg>
      </button>
    </div>
  </div>
</template>
