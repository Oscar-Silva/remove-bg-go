import { defineStore } from "pinia";
import { ref, computed } from "vue";
import { base64ToDataUrl } from "@/lib/utils";

export type AppState = "idle" | "loading" | "processing" | "done" | "error";

export interface HistoryItem {
  id: string;
  originalImage: string;
  resultImage: string;
  timestamp: number;
}

export const useAppStore = defineStore("app", () => {
  const currentState = ref<AppState>("idle");
  const originalImage = ref<string | null>(null);
  const resultImage = ref<string | null>(null);
  const statusMessage = ref<string>("");
  const error = ref<string | null>(null);
  const history = ref<HistoryItem[]>([]);
  const progress = ref<number>(0);
  const selectedModel = ref<string>("model_fp16.onnx");
  const downloadProgress = ref<{ downloaded: number; total: number }>({
    downloaded: 0,
    total: 0,
  });

  // Computed
  const isIdle = computed(() => currentState.value === "idle");
  const isLoading = computed(() => currentState.value === "loading");
  const isProcessing = computed(() => currentState.value === "processing");
  const isDone = computed(() => currentState.value === "done");
  const isError = computed(() => currentState.value === "error");
  const hasImage = computed(() => originalImage.value !== null);
  const hasResult = computed(() => resultImage.value !== null);

  const originalImageUrl = computed(() => {
    if (!originalImage.value) return null;
    return base64ToDataUrl(originalImage.value);
  });

  const resultImageUrl = computed(() => {
    if (!resultImage.value) return null;
    return base64ToDataUrl(resultImage.value);
  });

  // Actions
  function setIdle() {
    currentState.value = "idle";
    statusMessage.value = "";
    error.value = null;
    progress.value = 0;
  }

  function setLoading(message: string = "Loading...") {
    currentState.value = "loading";
    statusMessage.value = message;
    error.value = null;
  }

  function setProcessing(message: string = "Processing...") {
    currentState.value = "processing";
    statusMessage.value = message;
    error.value = null;
  }

  function setDone() {
    currentState.value = "done";
    statusMessage.value = "Complete!";
    progress.value = 100;

    // Add to history
    if (originalImage.value && resultImage.value) {
      const item: HistoryItem = {
        id: crypto.randomUUID(),
        originalImage: originalImage.value,
        resultImage: resultImage.value,
        timestamp: Date.now(),
      };
      history.value.unshift(item);

      // Keep only last 5 items
      if (history.value.length > 5) {
        history.value = history.value.slice(0, 5);
      }
    }
  }

  function setError(message: string) {
    currentState.value = "error";
    error.value = message;
    statusMessage.value = "Error";
    progress.value = 0;
  }

  function setOriginalImage(base64: string) {
    originalImage.value = base64;
  }

  function setResultImage(base64: string) {
    resultImage.value = base64;
  }

  function setStatusMessage(message: string) {
    statusMessage.value = message;
  }

  function setProgress(value: number) {
    progress.value = value;
  }

  function setDownloadProgress(downloaded: number, total: number) {
    downloadProgress.value = { downloaded, total };
  }

  function setSelectedModel(modelId: string) {
    selectedModel.value = modelId;
  }

  function reset() {
    currentState.value = "idle";
    originalImage.value = null;
    resultImage.value = null;
    statusMessage.value = "";
    error.value = null;
    progress.value = 0;
  }

  function removeFromHistory(id: string) {
    history.value = history.value.filter((item) => item.id !== id);
  }

  return {
    // State
    currentState,
    originalImage,
    resultImage,
    statusMessage,
    error,
    history,
    progress,
    selectedModel,
    downloadProgress,

    // Computed
    isIdle,
    isLoading,
    isProcessing,
    isDone,
    isError,
    hasImage,
    hasResult,
    originalImageUrl,
    resultImageUrl,

    // Actions
    setIdle,
    setLoading,
    setProcessing,
    setDone,
    setError,
    setOriginalImage,
    setResultImage,
    setStatusMessage,
    setProgress,
    setDownloadProgress,
    setSelectedModel,
    reset,
    removeFromHistory,
  };
});
