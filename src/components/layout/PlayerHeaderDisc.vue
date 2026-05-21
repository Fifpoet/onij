<template>
  <button
    v-if="playQueue.nowPlaying"
    type="button"
    class="shrink-0 rounded-full focus:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 overflow-hidden"
    :aria-label="'打开播放队列'"
    @click="router.push('/queue')"
  >
    <img
      :src="coverUrl"
      alt=""
      class="w-9 h-9 object-cover block rounded-full"
      :class="{ 'player-disc-spin': playQueue.isPlaying }"
    >
  </button>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { usePlayQueueStore } from '@/store/playQueue'

const router = useRouter()

const defaultAlbumCover = 'https://p1.music.126.net/6y-UleORITEDbvrOLV0Q8A==/5639395138885805.jpg'

const playQueue = usePlayQueueStore()

const coverUrl = computed(
  () => playQueue.nowPlaying?.cover_file_url || defaultAlbumCover,
)
</script>

<style scoped>
.player-disc-spin {
  animation: player-disc-rotate 8s linear infinite;
}

@keyframes player-disc-rotate {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
</style>
