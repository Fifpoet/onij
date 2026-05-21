<template>
  <div
    v-if="song"
    class="flex min-w-0 max-w-[min(100%,28rem)] items-center gap-2 md:max-w-[min(100%,36rem)]"
  >
    <button
      type="button"
      class="header-player-btn shrink-0 overflow-hidden rounded-full border-0 bg-transparent p-0 shadow-none"
      aria-label="打开歌曲详情"
      @click="openSongDetail"
    >
      <img
        :src="coverUrl"
        alt=""
        class="block h-9 w-9 rounded-full object-cover"
        :class="{ 'player-disc-spin': playQueue.isPlaying }"
      >
    </button>

    <p class="min-w-0 flex-1 truncate text-left text-[13px] leading-tight md:text-sm">
      <RouterLink
        v-if="song.album_id"
        :to="{ path: '/album', query: { id: String(song.album_id) } }"
        class="font-semibold text-gray-900 hover:underline dark:text-gray-100"
        @click.stop
      >
        {{ song.name }}
      </RouterLink>
      <span v-else class="font-semibold text-gray-900 dark:text-gray-100">{{ song.name }}</span>
      <span class="text-gray-400 dark:text-gray-500"> - </span>
      <template v-for="(artist, i) in lineArtists" :key="artist.artist_id">
        <RouterLink
          :to="{ path: '/artist', query: { ids: String(artist.artist_id) } }"
          class="text-gray-600 hover:text-blue-600 hover:underline dark:text-gray-400 dark:hover:text-blue-400"
          @click.stop
        >
          {{ artist.artist_name }}
        </RouterLink>
        <span v-if="i < lineArtists.length - 1" class="text-gray-400 dark:text-gray-500">/</span>
      </template>
    </p>

    <button
      type="button"
      class="header-player-btn flex h-9 w-9 shrink-0 cursor-pointer items-center justify-center border-0 bg-transparent p-0 shadow-none"
      :aria-label="playQueue.isPlaying ? '暂停' : '播放'"
      @click="onTogglePlay"
    >
      <PlayerHeaderWave :active="playQueue.isPlaying" />
    </button>

    <button
      type="button"
      class="header-player-btn hidden h-9 w-9 shrink-0 cursor-pointer items-center justify-center border-0 bg-transparent p-0 shadow-none lg:flex"
      aria-label="查看播放队列"
      @click="openQueue"
    >
      <n-icon :component="ListOutline" :size="22" class="text-gray-600 dark:text-gray-400" />
    </button>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { NIcon } from 'naive-ui'
import { ListOutline } from '@vicons/ionicons5'
import { usePlayQueueStore } from '@/store/playQueue'
import { useLyricsPanelStore } from '@/store/lyricsPanel'
import { useMobileLyricsSheetStore } from '@/store/mobileLyricsSheet'
import { neteasePlayerControl } from '@/player/neteasePlayerControl'
import PlayerHeaderWave from '@/components/layout/PlayerHeaderWave.vue'

const defaultAlbumCover =
  'https://p1.music.126.net/6y-UleORITEDbvrOLV0Q8A==/5639395138885805.jpg'

const router = useRouter()
const playQueue = usePlayQueueStore()
const lyricsPanel = useLyricsPanelStore()
const mobileSheet = useMobileLyricsSheetStore()

const song = computed(() => playQueue.nowPlaying)

function openQueue() {
  void router.push('/queue')
}

const coverUrl = computed(
  () => song.value?.cover_file_url || defaultAlbumCover,
)

const lineArtists = computed(() => {
  const s = song.value
  if (!s) return []
  if (s.display_artists?.length) return s.display_artists
  return s.artists ?? []
})

function openSongDetail() {
  if (typeof window !== 'undefined' && window.matchMedia('(min-width: 1024px)').matches) {
    lyricsPanel.open()
  } else {
    mobileSheet.open()
  }
}

function onTogglePlay() {
  if (!playQueue.nowPlaying) return
  neteasePlayerControl.togglePlay()
}
</script>

<style scoped>
.header-player-btn {
  -webkit-tap-highlight-color: transparent;
  outline: none !important;
  border: none !important;
  box-shadow: none !important;
}
.header-player-btn:hover,
.header-player-btn:focus,
.header-player-btn:focus-visible,
.header-player-btn:active {
  outline: none !important;
  border: none !important;
  box-shadow: none !important;
  background: transparent !important;
}
.header-player-btn:focus-visible {
  outline: 2px solid #3b82f6 !important;
  outline-offset: 2px;
}

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
