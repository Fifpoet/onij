<template>
  <button
    v-if="playQueue.nowPlaying"
    type="button"
    class="lyrics-entry-btn shrink-0 rounded-2xl text-sm font-semibold text-gray-800 dark:text-gray-100 focus:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
    aria-label="歌词"
    @click="lyrics.open()"
  >
    词
  </button>
</template>

<script setup lang="ts">
import { usePlayQueueStore } from '@/store/playQueue'
import { useLyricsPanelStore } from '@/store/lyricsPanel'

const playQueue = usePlayQueueStore()
const lyrics = useLyricsPanelStore()
</script>

<style scoped>
.lyrics-entry-btn {
  width: 2.25rem;
  height: 2.25rem;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.92);
  border: 1px solid rgba(0, 0, 0, 0.08);
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.06);
}
.dark .lyrics-entry-btn {
  background: rgba(30, 30, 36, 0.92);
  border-color: rgba(255, 255, 255, 0.12);
}
.lyrics-entry-btn:hover {
  filter: brightness(0.97);
}
</style>
