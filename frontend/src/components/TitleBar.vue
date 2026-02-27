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
    class="h-10 bg-surface flex items-center justify-between px-4 select-none border-b border-border"
    style="-webkit-app-region: drag"
  >
    <!-- Logo and Title -->
    <div class="flex items-center gap-2">
      <div class="w-6 h-6 rounded-lg bg-accent flex items-center justify-center">
        <svg
          class="w-4 h-4 text-white"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
        >
          <path d="M21 12V7H5a2 2 0 0 1 0-4h14v4" />
          <path d="M3 5v14a2 2 0 0 0 2 2h16v-5" />
          <path d="M18 12a2 2 0 0 0 0 4h4v-4Z" />
        </svg>
      </div>
      <span class="text-sm font-medium text-textPrimary">Remove Background</span>
    </div>

    <!-- Window Controls -->
    <div
      class="flex items-center gap-1"
      style="-webkit-app-region: no-drag"
    >
      <!-- Minimize -->
      <button
        @click="minimize"
        class="w-8 h-8 rounded-lg flex items-center justify-center text-textMuted hover:text-textPrimary hover:bg-surfaceElevated transition-all duration-150"
      >
        <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M5 12h14" />
        </svg>
      </button>

      <!-- Maximize -->
      <button
        @click="toggleMaximize"
        class="w-8 h-8 rounded-lg flex items-center justify-center text-textMuted hover:text-textPrimary hover:bg-surfaceElevated transition-all duration-150"
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
