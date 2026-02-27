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
  <div class="h-full flex flex-col bg-background">
    <!-- Title Bar -->
    <TitleBar />

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
