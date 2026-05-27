<template>
  <div ref="pageRef" class="ktv-page">
    <div class="ktv-cover-wrap">
      <img
        v-if="coverUrl"
        :src="coverUrl"
        :alt="song?.name ?? ''"
        class="ktv-cover"
        crossorigin="anonymous"
      >
    </div>

    <div v-if="ktv.lyricsOn && song" class="ktv-lyrics">
      <div v-if="dual.line1" class="ktv-lyric-row">
        <KtvLyricLine :text="dual.line1" :progress="dual.line1Progress" />
      </div>
      <div v-else class="ktv-lyric-row">
        <p class="ktv-lyric-hint">{{ lyricHint }}</p>
      </div>
      <div v-if="dual.line2" class="ktv-lyric-row">
        <KtvLyricLine :text="dual.line2" :progress="dual.line2Progress" />
      </div>
    </div>

    <div
      class="ktv-dock"
      @mouseenter="dockHover = true"
      @mouseleave="onDockLeave"
    >
      <div class="ktv-dock-panel" :class="{ 'ktv-dock-panel--open': dockOpen }">
        <button
          type="button"
          class="ktv-dock-handle"
          aria-label="展开播放控制"
          @click="dockPinned = !dockPinned"
        >
          <span class="ktv-dock-handle__chev" aria-hidden="true" />
        </button>
        <div class="ktv-dock-bar">
          <button
            type="button"
            class="ktv-dock-btn"
            :class="{ 'ktv-dock-btn--on': ktv.lyricsOn }"
            @click="ktv.toggleLyrics()"
          >
            词
          </button>
          <button
            type="button"
            class="ktv-dock-btn"
            :class="{
              'ktv-dock-btn--on': ktv.accompanimentOn,
              'ktv-dock-btn--loading': accompanimentLoading,
            }"
            :title="accompanimentTitle"
            @click="ktv.toggleAccompaniment()"
          >
            伴
          </button>
          <button
            type="button"
            class="ktv-dock-btn ktv-dock-btn--icon"
            :aria-label="isFullscreen ? '退出全屏' : '全屏'"
            @click="toggleFullscreen"
          >
            <n-icon :component="isFullscreen ? ContractOutline : ExpandOutline" :size="22" />
          </button>
          <button
            type="button"
            class="ktv-dock-btn ktv-dock-btn--icon"
            :disabled="!canTogglePlay"
            :aria-label="playQueue.isPlaying ? '暂停' : '播放'"
            @click="onTogglePlay"
          >
            <n-icon :component="playIcon" :size="28" />
          </button>
          <button
            type="button"
            class="ktv-dock-btn ktv-dock-btn--icon"
            :disabled="nextDisabled"
            aria-label="下一曲"
            @click="playQueue.skipToNext()"
          >
            <n-icon :component="PlaySkipForwardOutline" :size="24" />
          </button>
          <div class="ktv-dock-volume">
            <n-icon :component="volumeIcon" :size="20" class="ktv-dock-volume__icon" />
            <n-slider
              class="ktv-dock-volume__slider"
              :value="playQueue.volume"
              :min="0"
              :max="100"
              :step="1"
              :tooltip="false"
              :show-tooltip="false"
              aria-label="音量"
              @update:value="onVolumeChange"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount, watch } from 'vue'
import { NIcon, NSlider } from 'naive-ui'
import {
  ContractOutline,
  ExpandOutline,
  PauseCircleOutline,
  PlayCircleOutline,
  PlaySkipForwardOutline,
  VolumeHighOutline,
  VolumeLowOutline,
  VolumeMediumOutline,
  VolumeMuteOutline,
} from '@vicons/ionicons5'
import KtvLyricLine from '@/components/ktv/KtvLyricLine.vue'
import { usePlayQueueStore } from '@/store/playQueue'
import { useKtvStore } from '@/store/ktv'
import { usePlayerLyrics } from '@/composables/usePlayerLyrics'
import { ktvDualLineState } from '@/util/lrc'
import { neteasePlayerControl } from '@/player/neteasePlayerControl'
import {
  getInstrumentalStatus,
  prefetchInstrumental,
} from '@/player/instrumentalCache'
import { prefetchPlayUrl } from '@/player/songPlayUrlCache'

const defaultCover =
  'https://p1.music.126.net/6y-UleORITEDbvrOLV0Q8A==/5639395138885805.jpg'

const playQueue = usePlayQueueStore()
const ktv = useKtvStore()
const { lines, lyricError, loading } = usePlayerLyrics()

const dockHover = ref(false)
const dockPinned = ref(false)
const pageRef = ref<HTMLElement | null>(null)
const isFullscreen = ref(false)

const dockOpen = computed(() => dockHover.value || dockPinned.value)

const song = computed(() => playQueue.nowPlaying)
const coverUrl = computed(() => song.value?.cover_file_url || defaultCover)

const dual = computed(() =>
  ktvDualLineState(lines.value, playQueue.playbackCurrentSec),
)

const lyricHint = computed(() => {
  if (loading.value) return '歌词加载中…'
  if (lyricError.value) return lyricError.value
  return song.value?.name ?? ''
})

const canTogglePlay = computed(
  () => !!playQueue.nowPlaying || playQueue.queue.length > 0,
)

const nextDisabled = computed(
  () => !playQueue.nowPlaying || playQueue.queue.length === 0,
)

const playIcon = computed(() =>
  playQueue.isPlaying ? PauseCircleOutline : PlayCircleOutline,
)

const volumeIcon = computed(() => {
  const v = playQueue.volume
  if (v <= 0) return VolumeMuteOutline
  if (v < 34) return VolumeLowOutline
  if (v < 67) return VolumeMediumOutline
  return VolumeHighOutline
})

const accompanimentLoading = computed(() => {
  const id = song.value?.id
  if (!id || !ktv.accompanimentOn) return false
  return getInstrumentalStatus(id) === 'loading'
})

const accompanimentTitle = computed(() => {
  const id = song.value?.id
  if (!id) return ktv.accompanimentOn ? '伴奏：开' : '伴奏：关'
  const st = getInstrumentalStatus(id)
  if (st === 'loading') return '伴奏提取中…'
  if (st === 'error') return '伴奏提取失败，点击重试'
  return ktv.accompanimentOn ? '伴奏：开（点击切原唱）' : '原唱：开（点击切伴奏）'
})

function prefetchKtvAssets() {
  const current = playQueue.nowPlaying
  if (current) {
    prefetchPlayUrl(current.id)
    prefetchInstrumental(current.id)
  }
  const n = Math.min(3, playQueue.queue.length)
  for (let i = 0; i < n; i++) {
    prefetchPlayUrl(playQueue.queue[i].id)
    prefetchInstrumental(playQueue.queue[i].id)
  }
}

function onDockLeave() {
  dockHover.value = false
}

function onTogglePlay() {
  if (!playQueue.nowPlaying && playQueue.queue.length > 0) {
    playQueue.startIfIdle()
    return
  }
  neteasePlayerControl.togglePlay()
}

function isSpaceBlockedTarget(target: EventTarget | null): boolean {
  if (!(target instanceof HTMLElement)) return false
  const tag = target.tagName
  if (tag === 'INPUT' || tag === 'TEXTAREA' || tag === 'SELECT') return true
  if (target.isContentEditable) return true
  return false
}

function onPageKeydown(e: KeyboardEvent) {
  if (e.code !== 'Space' && e.key !== ' ') return
  if (e.repeat) return
  if (isSpaceBlockedTarget(e.target)) return
  e.preventDefault()
  if (!canTogglePlay.value) return
  onTogglePlay()
}

function onVolumeChange(v: number | number[]) {
  playQueue.setVolume(Math.round(Array.isArray(v) ? v[0] : v))
}

function onFullscreenChange() {
  isFullscreen.value = document.fullscreenElement === pageRef.value
}

function toggleFullscreen() {
  const el = pageRef.value
  if (!el) return
  if (!document.fullscreenElement) {
    void el.requestFullscreen?.()
  } else {
    void document.exitFullscreen?.()
  }
}

onMounted(() => {
  ktv.enterKtv()
  prefetchKtvAssets()
  document.addEventListener('fullscreenchange', onFullscreenChange)
  document.addEventListener('keydown', onPageKeydown)
})

watch(
  () => playQueue.nowPlaying?.id,
  () => {
    if (!ktv.active) return
    prefetchKtvAssets()
  },
)

onBeforeUnmount(() => {
  ktv.leaveKtv()
  document.removeEventListener('fullscreenchange', onFullscreenChange)
  document.removeEventListener('keydown', onPageKeydown)
  if (document.fullscreenElement === pageRef.value) {
    void document.exitFullscreen?.()
  }
})
</script>

<style scoped>
.ktv-page {
  position: relative;
  width: 100%;
  height: 100%;
  overflow: hidden;
  overscroll-behavior: none;
  background: rgb(15 23 42);
}

.ktv-cover-wrap {
  position: absolute;
  inset: 0;
  z-index: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1rem 1rem clamp(5rem, 18vh, 9rem);
  box-sizing: border-box;
  pointer-events: none;
}

.ktv-cover {
  display: block;
  max-width: min(92vw, 520px);
  max-height: 100%;
  width: auto;
  height: auto;
  object-fit: contain;
  user-select: none;
  box-shadow: 0 20px 48px rgb(0 0 0 / 0.35);
}

.ktv-lyrics {
  position: absolute;
  left: 0;
  right: 0;
  bottom: clamp(3.25rem, 12vh, 5.5rem);
  z-index: 10;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.65rem;
  padding: 0 1rem;
  max-width: 100%;
  box-sizing: border-box;
  pointer-events: none;
}

.ktv-lyric-row {
  width: 100%;
  display: flex;
  justify-content: center;
}

.ktv-page:fullscreen {
  width: 100vw;
  height: 100vh;
  max-height: 100dvh;
}

.ktv-lyric-hint {
  margin: 0;
  width: min(96vw, 52rem);
  font-size: clamp(1.75rem, 7vw, 3rem);
  font-weight: 700;
  line-height: 1.35;
  color: rgb(255 255 255 / 0.85);
  text-shadow: 0 2px 14px rgb(0 0 0 / 0.65);
  text-align: center;
}

.ktv-dock {
  position: absolute;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 40;
  display: flex;
  justify-content: center;
  padding-bottom: max(0px, env(safe-area-inset-bottom, 0px));
  pointer-events: none;
}

.ktv-dock-panel {
  pointer-events: auto;
  display: flex;
  flex-direction: column;
  align-items: center;
  width: fit-content;
  max-width: min(100%, 42rem);
  background: rgb(15 23 42 / 0.92);
  border: 1px solid rgb(255 255 255 / 0.14);
  border-bottom: none;
  border-radius: 0.75rem 0.75rem 0 0;
  backdrop-filter: blur(14px);
  box-shadow: 0 -6px 28px rgb(0 0 0 / 0.35);
  overflow: hidden;
}

.ktv-dock-handle {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2.75rem;
  height: 1.125rem;
  padding: 0;
  border: 0;
  border-bottom: 1px solid transparent;
  background: transparent;
  cursor: pointer;
  outline: none;
  box-shadow: none;
  -webkit-tap-highlight-color: transparent;
  transition: background-color 0.15s ease;
}

.ktv-dock-handle:focus,
.ktv-dock-handle:focus-visible,
.ktv-dock-handle:active,
.ktv-dock-handle:hover {
  outline: none !important;
  box-shadow: none !important;
  border: 0 !important;
  border-bottom: 1px solid transparent !important;
  background: transparent;
}

.ktv-dock-panel--open .ktv-dock-handle {
  border-bottom-color: rgb(255 255 255 / 0.1);
}

.ktv-dock-panel--open .ktv-dock-handle:focus,
.ktv-dock-panel--open .ktv-dock-handle:focus-visible,
.ktv-dock-panel--open .ktv-dock-handle:active,
.ktv-dock-panel--open .ktv-dock-handle:hover {
  border-bottom-color: rgb(255 255 255 / 0.1) !important;
  background: transparent;
}

.ktv-dock-handle__chev {
  display: block;
  width: 0;
  height: 0;
  border-left: 6px solid transparent;
  border-right: 6px solid transparent;
  border-bottom: 7px solid rgb(255 255 255 / 0.88);
  transition: transform 0.2s ease;
}

.ktv-dock-panel--open .ktv-dock-handle__chev {
  transform: rotate(180deg);
}

.ktv-dock-bar {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-wrap: nowrap;
  gap: 0.45rem 0.55rem;
  width: max-content;
  max-height: 0;
  padding: 0 0.65rem;
  opacity: 0;
  overflow: hidden;
  transition:
    max-height 0.22s ease,
    opacity 0.18s ease,
    padding 0.22s ease;
}

.ktv-dock-panel--open .ktv-dock-bar {
  max-height: 4rem;
  padding: 0.45rem 0.65rem 0.55rem;
  opacity: 1;
}

.ktv-dock-btn {
  flex-shrink: 0;
  min-width: 2.35rem;
  height: 2.35rem;
  padding: 0 0.55rem;
  border: 0;
  border-radius: 999px;
  background: rgb(255 255 255 / 0.1);
  color: rgb(226 232 240);
  font-size: 0.875rem;
  font-weight: 600;
  cursor: pointer;
  outline: none;
  box-shadow: none;
  -webkit-tap-highlight-color: transparent;
  transition: background-color 0.12s ease;
}

.ktv-dock-btn:focus,
.ktv-dock-btn:focus-visible,
.ktv-dock-btn:active,
.ktv-dock-btn:hover:not(:disabled) {
  outline: none !important;
  box-shadow: none !important;
  border: 0 !important;
}

.ktv-dock-btn--on {
  background: rgb(56 189 248 / 0.35);
  color: #fff;
}

.ktv-dock-btn--loading {
  opacity: 0.75;
  animation: ktv-dock-pulse 1s ease-in-out infinite;
}

@keyframes ktv-dock-pulse {
  0%,
  100% {
    opacity: 0.55;
  }
  50% {
    opacity: 1;
  }
}

.ktv-dock-btn--disabled {
  opacity: 0.45;
  cursor: not-allowed;
}

.ktv-dock-btn--icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 2.5rem;
  padding: 0;
}

.ktv-dock-volume {
  display: flex;
  align-items: center;
  gap: 0.3rem;
  width: 7.5rem;
  flex-shrink: 0;
}

.ktv-dock-volume__icon {
  flex-shrink: 0;
  color: rgb(203 213 225);
}

.ktv-dock-volume__slider {
  flex: 1;
  min-width: 0;
}

.ktv-dock-volume__slider :deep(.n-slider-rail) {
  background-color: rgb(255 255 255 / 0.2) !important;
}

.ktv-dock-volume__slider :deep(.n-slider-rail__fill) {
  background-color: rgb(56 189 248) !important;
}
</style>
