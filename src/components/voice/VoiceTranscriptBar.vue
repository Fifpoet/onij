<template>
  <div class="voice-bar" aria-live="polite">
    <span class="voice-bar__status">{{ statusLabel }}</span>
    <span v-if="transcript" class="voice-bar__text">{{ transcript }}</span>
    <span v-else class="voice-bar__text voice-bar__text--placeholder">等待语音…</span>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useVoicePipeline } from '@/composables/useVoicePipeline'

const { transcript, status, listening } = useVoicePipeline()

const statusLabel = computed(() => {
  if (!listening.value) return '麦克风未开'
  switch (status.value) {
    case 'recording':
      return '录音中'
    case 'transcribing':
      return '识别中'
    case 'intent':
      return '理解中'
    case 'listening':
      return '监听中'
    default:
      return status.value
  }
})
</script>

<style scoped>
.voice-bar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 80;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.6rem;
  min-height: 1.35rem;
  padding: 0.15rem 0.75rem;
  padding-top: max(0.15rem, env(safe-area-inset-top));
  pointer-events: none;
  background: transparent;
}
.voice-bar__status {
  flex-shrink: 0;
  font-size: 0.68rem;
  color: rgb(156 163 175);
}
.voice-bar__text {
  max-width: min(70vw, 36rem);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 0.72rem;
  line-height: 1.2;
  color: rgb(156 163 175);
}
.voice-bar__text--placeholder {
  color: rgb(209 213 219);
}
.dark .voice-bar__status,
.dark .voice-bar__text {
  color: rgb(156 163 175);
}
</style>
