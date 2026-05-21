<template>
  <div
    class="mobile-player-bar fixed inset-x-0 bottom-0 z-[45] lg:hidden"
    :style="themeStyle"
    role="region"
    aria-label="正在播放"
  >
    <div class="flex h-12 items-stretch">
      <div class="h-12 w-12 shrink-0">
        <img
          :src="coverSrc"
          alt=""
          class="h-full w-full object-cover"
          :class="hasSong ? '' : 'opacity-55'"
          crossorigin="anonymous"
        >
      </div>

      <button
        type="button"
        class="mobile-tap min-w-0 flex-1 border-0 bg-transparent px-2"
        :disabled="!hasSong && playQueue.queue.length === 0"
        :aria-label="showControls ? '播放控制' : '当前歌词'"
        @click="onCenterTap"
      >
        <Transition name="bar-swap" mode="out-in">
          <div
            v-if="showControls"
            key="controls"
            class="flex h-full items-center justify-center gap-5"
          >
            <button
              type="button"
              class="mobile-tap border-0 bg-transparent p-0"
              :disabled="!canTogglePlay"
              aria-label="播放或暂停"
              @click.stop="onTogglePlay"
            >
              <n-icon :component="playIcon" :size="30" :color="theme.controlFg" />
            </button>
            <button
              type="button"
              class="mobile-tap border-0 bg-transparent p-0"
              :disabled="nextDisabled"
              aria-label="下一曲"
              @click.stop="playQueue.skipToNext()"
            >
              <n-icon :component="PlaySkipForwardOutline" :size="24" :color="theme.controlFg" />
            </button>
          </div>
          <p
            v-else
            key="lyric"
            class="m-0 w-full truncate text-left text-[13px] leading-snug"
            :style="{ color: theme.fg }"
          >
            {{ centerLabel }}
          </p>
        </Transition>
      </button>

      <button
        v-if="showQueueEntry"
        type="button"
        class="mobile-tap flex h-12 w-10 shrink-0 items-center justify-center border-0 bg-transparent p-0"
        aria-label="查看播放队列"
        @click="openQueue"
      >
        <n-icon :component="ListOutline" :size="20" :color="theme.controlFg" />
      </button>

      <button
        type="button"
        class="mobile-tap flex h-12 w-10 shrink-0 items-center justify-center border-0 bg-transparent p-0"
        aria-label="展开歌词"
        @click="sheet.open()"
      >
        <n-icon :component="ChevronUpOutline" :size="20" :color="theme.controlFg" />
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { NIcon } from 'naive-ui'
import {
  PlayCircleOutline,
  PauseCircleOutline,
  PlaySkipForwardOutline,
  ChevronUpOutline,
  ListOutline,
} from '@vicons/ionicons5'
import { usePlayQueueStore } from '@/store/playQueue'
import { useMobileLyricsSheetStore } from '@/store/mobileLyricsSheet'
import { usePlayerLyrics } from '@/composables/usePlayerLyrics'
import { useCoverTheme } from '@/composables/useCoverTheme'
import { neteasePlayerControl } from '@/player/neteasePlayerControl'

const defaultCover =
  'https://p1.music.126.net/6y-UleORITEDbvrOLV0Q8A==/5639395138885805.jpg'

const router = useRouter()
const playQueue = usePlayQueueStore()
const sheet = useMobileLyricsSheetStore()
const { currentLineText } = usePlayerLyrics()

const showQueueEntry = computed(
  () => !!playQueue.nowPlaying || playQueue.queue.length > 0,
)

function openQueue() {
  void router.push('/queue')
}

const showControls = ref(false)

const coverSrc = computed(
  () => playQueue.nowPlaying?.cover_file_url || defaultCover,
)
const { theme, themeStyle } = useCoverTheme(coverSrc)

const hasSong = computed(() => !!playQueue.nowPlaying)

const centerLabel = computed(() => {
  if (!playQueue.nowPlaying) {
    if (playQueue.queue.length > 0) return '待播放 · 点击开始'
    return '暂无播放'
  }
  return currentLineText.value
})

const canTogglePlay = computed(
  () => !!playQueue.nowPlaying || playQueue.queue.length > 0,
)

const nextDisabled = computed(
  () => !playQueue.nowPlaying || playQueue.queue.length === 0,
)

const playIcon = computed(() =>
  playQueue.nowPlaying && playQueue.isPlaying ? PauseCircleOutline : PlayCircleOutline,
)

watch(
  () => playQueue.nowPlaying?.id,
  () => {
    showControls.value = false
  },
)

function onCenterTap() {
  if (!hasSong.value && playQueue.queue.length > 0) {
    playQueue.startIfIdle()
    return
  }
  if (!hasSong.value) return
  showControls.value = !showControls.value
}

function onTogglePlay() {
  if (!playQueue.nowPlaying && playQueue.queue.length > 0) {
    playQueue.startIfIdle()
    return
  }
  neteasePlayerControl.togglePlay()
}
</script>

<style scoped>
.mobile-player-bar {
  background: var(--cover-bg);
  padding-bottom: env(safe-area-inset-bottom, 0px);
  min-height: calc(3rem + env(safe-area-inset-bottom, 0px));
}

.mobile-tap {
  -webkit-tap-highlight-color: transparent;
  outline: none !important;
  box-shadow: none !important;
}
.mobile-tap:focus,
.mobile-tap:focus-visible,
.mobile-tap:active {
  outline: none !important;
  box-shadow: none !important;
  border: none !important;
}
.mobile-tap:disabled {
  opacity: 0.45;
}

.bar-swap-enter-active,
.bar-swap-leave-active {
  transition: opacity 0.12s ease;
}
.bar-swap-enter-from,
.bar-swap-leave-to {
  opacity: 0;
}
</style>
