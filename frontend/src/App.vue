<script setup lang="ts">
import { useAppStore } from '@/stores/app'
import TitleBar from '@/components/TitleBar.vue'
import DropZone from '@/components/DropZone.vue'
import ProcessingView from '@/components/ProcessingView.vue'
import ResultView from '@/components/ResultView.vue'
import Toast from '@/components/Toast.vue'

const store = useAppStore()
</script>

<template>
  <div class="h-full flex flex-col bg-background relative selection:bg-accent selection:text-white">
    <!-- Static background layers for performance (60 FPS) -->
    <div class="absolute inset-0 overflow-hidden pointer-events-none z-0">
      <div class="absolute -top-[20%] -left-[10%] w-[50%] h-[50%] rounded-full bg-accent/10 blur-[120px]" />
      <div class="absolute top-[30%] -right-[15%] w-[60%] h-[60%] rounded-full bg-purple-600/10 blur-[120px]" />
      <div class="absolute -bottom-[20%] left-[20%] w-[50%] h-[50%] rounded-full bg-blue-600/10 blur-[120px]" />
    </div>

    <!-- Title Bar -->
    <TitleBar class="z-10 relative" />

    <!-- Main Content -->
    <main class="flex-1 flex items-center justify-center p-8 overflow-auto">
      <div class="w-full max-w-2xl">
        <!-- Idle: Show DropZone -->
        <Transition name="fade-scale" mode="out-in">
          <DropZone v-if="store.isIdle" />
        </Transition>

        <!-- Loading/Processing: Show ProcessingView -->
        <Transition name="fade-scale" mode="out-in">
          <ProcessingView v-if="store.isLoading || store.isProcessing" />
        </Transition>

        <!-- Done: Show ResultView -->
        <Transition name="fade-scale" mode="out-in">
          <ResultView v-if="store.isDone" />
        </Transition>
      </div>
    </main>

    <!-- Toast Notifications -->
    <Toast />
  </div>
</template>

<style scoped>
.fade-scale-enter-active,
.fade-scale-leave-active {
  transition: all 0.3s ease;
}

.fade-scale-enter-from {
  opacity: 0;
  transform: scale(0.95);
}

.fade-scale-leave-to {
  opacity: 0;
  transform: scale(1.05);
}
</style>
