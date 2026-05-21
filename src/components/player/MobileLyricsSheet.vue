<template>
  <Teleport to="body">
    <Transition name="mobile-sheet">
      <div
        v-if="sheet.visible"
        class="mobile-lyrics-fullscreen fixed inset-0 z-[250] flex flex-col overflow-hidden lg:hidden"
        :style="themeStyle"
        role="dialog"
        aria-modal="true"
        aria-label="歌词详情"
      >
        <div class="pointer-events-none absolute inset-0 overflow-hidden" aria-hidden="true">
          <img
            :src="coverSrc"
            alt=""
            class="absolute inset-0 h-full w-full scale-110 object-cover opacity-50 blur-3xl saturate-150"
            crossorigin="anonymous"
          >
          <div
            class="absolute inset-0"
            :style="{
              background: `linear-gradient(180deg, ${theme.bg}ee 0%, ${theme.bg} 45%, ${theme.bg} 100%)`,
            }"
          />
        </div>

        <div
          class="relative flex h-11 shrink-0 items-center justify-end px-2 pt-[max(0.25rem,env(safe-area-inset-top))]"
        >
          <button
            type="button"
            class="mobile-tap flex h-10 w-10 items-center justify-center border-0 bg-transparent p-0"
            aria-label="收起"
            @click="sheet.close()"
          >
            <n-icon :component="ChevronUpOutline" :size="24" :color="theme.controlFg" />
          </button>
        </div>

        <div class="relative flex min-h-0 flex-1 flex-col px-3 pb-[max(0.75rem,env(safe-area-inset-bottom))]">
          <Transition name="sheet-phase" mode="out-in">
            <div
              v-if="!sheet.showLyricsScroll"
              key="cover"
              class="flex min-h-0 flex-1 flex-col items-center justify-center gap-3"
            >
              <button
                type="button"
                class="mobile-tap w-[min(76vw,300px)] shrink-0 border-0 bg-transparent p-0"
                aria-label="点击查看滚动歌词"
                @click="onCoverTap"
              >
                <img
                  :src="coverSrc"
                  alt=""
                  class="aspect-square w-full object-cover"
                  crossorigin="anonymous"
                >
              </button>

              <div v-if="song" class="w-full max-w-md px-1 text-center">
                <h2 class="text-lg font-bold leading-snug" :style="{ color: theme.fg }">
                  {{ song.name }}
                </h2>
                <p
                  v-if="artistLine"
                  class="mt-0.5 truncate text-sm"
                  :style="{ color: theme.fgMuted }"
                >
                  {{ artistLine }}
                </p>
              </div>
              <p v-else class="text-sm" :style="{ color: theme.fgMuted }">暂无播放</p>

              <div class="mobile-transport w-full max-w-md">
                <PlayTransportStrip surface="lyrics" />
              </div>
            </div>

            <div
              v-else
              key="lyrics"
              ref="lyricScrollRef"
              class="mobile-lyrics-scroll min-h-0 flex-1 overflow-y-auto overscroll-contain py-2"
              role="button"
              tabindex="0"
              aria-label="点击返回封面"
              @click="onLyricsAreaTap"
              @touchstart.passive="onLyricsTouchStart"
              @touchmove.passive="onLyricsTouchMove"
              @wheel="onUserLyricsInteraction"
            >
              <p
                v-if="lyricError"
                class="py-8 text-center"
                :style="{ color: theme.fgMuted }"
              >
                {{ lyricError }}
              </p>
              <template v-else-if="lines.length">
                <p
                  v-for="(line, i) in lines"
                  :key="`${line.time}-${i}`"
                  :data-lyric-i="i"
                  class="mobile-lyric-line block py-2 text-center transition-colors duration-200"
                  :class="i === activeIdx ? 'mobile-lyric-line--active' : ''"
                  :style="i === activeIdx ? { color: theme.fg } : { color: theme.fgSubtle }"
                >
                  {{ line.text || '\u00a0' }}
                </p>
              </template>
              <p v-else class="py-8 text-center" :style="{ color: theme.fgMuted }">
                暂无歌词
              </p>
            </div>
          </Transition>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick, onBeforeUnmount } from 'vue'
import { NIcon } from 'naive-ui'
import { ChevronUpOutline } from '@vicons/ionicons5'
import { usePlayQueueStore } from '@/store/playQueue'
import { useMobileLyricsSheetStore } from '@/store/mobileLyricsSheet'
import { usePlayerLyrics } from '@/composables/usePlayerLyrics'
import { useCoverTheme } from '@/composables/useCoverTheme'
import PlayTransportStrip from '@/components/player/PlayTransportStrip.vue'
import type { ViewMusicListItem } from '@/api/view/music'

const defaultCover =
  'https://p1.music.126.net/6y-UleORITEDbvrOLV0Q8A==/5639395138885805.jpg'

const LYRIC_IDLE_MS = 10_000

const sheet = useMobileLyricsSheetStore()
const playQueue = usePlayQueueStore()
const { lines, lyricError, activeIdx } = usePlayerLyrics()

const lyricScrollRef = ref<HTMLElement | null>(null)
const followLyrics = ref(true)
let idleSnapTimer: ReturnType<typeof setTimeout> | null = null
let lyricsTouchMoved = false

const coverSrc = computed(
  () => playQueue.nowPlaying?.cover_file_url || defaultCover,
)
const { theme, themeStyle } = useCoverTheme(coverSrc)

const song = computed<ViewMusicListItem | null>(() => playQueue.nowPlaying)

const artistLine = computed(() => {
  const s = song.value
  if (!s) return ''
  const list =
    s.display_artists && s.display_artists.length > 0 ? s.display_artists : s.artists
  const names = list?.map((a) => a.artist_name).filter(Boolean) ?? []
  const main = names.join(' / ')
  if (main && s.album_name) return `${main} · ${s.album_name}`
  return main || s.album_name || ''
})

function clearIdleSnap() {
  if (idleSnapTimer) {
    clearTimeout(idleSnapTimer)
    idleSnapTimer = null
  }
}

function scheduleIdleSnap() {
  clearIdleSnap()
  idleSnapTimer = setTimeout(() => {
    idleSnapTimer = null
    followLyrics.value = true
    scrollActiveLineIntoView(true)
  }, LYRIC_IDLE_MS)
}

function onLyricsTouchStart() {
  lyricsTouchMoved = false
}

function onLyricsTouchMove() {
  lyricsTouchMoved = true
  onUserLyricsInteraction()
}

function onUserLyricsInteraction() {
  followLyrics.value = false
  scheduleIdleSnap()
}

function onLyricsAreaTap() {
  if (lyricsTouchMoved) {
    lyricsTouchMoved = false
    return
  }
  sheet.backToCover()
}

function scrollActiveLineIntoView(smooth: boolean) {
  const wrap = lyricScrollRef.value
  const i = activeIdx.value
  if (!wrap || i < 0) return
  const el = wrap.querySelector(`[data-lyric-i="${i}"]`) as HTMLElement | null
  el?.scrollIntoView({ block: 'center', behavior: smooth ? 'smooth' : 'auto' })
}

function onCoverTap() {
  if (!song.value) return
  sheet.enterLyricsScroll()
}

watch(
  () => sheet.visible,
  (v) => {
    document.body.style.overflow = v ? 'hidden' : ''
    if (v) {
      followLyrics.value = true
      clearIdleSnap()
    } else {
      clearIdleSnap()
    }
  },
)

watch(
  () => sheet.showLyricsScroll,
  async (scroll) => {
    if (scroll && sheet.visible) {
      followLyrics.value = true
      clearIdleSnap()
      await nextTick()
      scrollActiveLineIntoView(false)
    }
  },
)

watch(activeIdx, async () => {
  if (!sheet.visible || !sheet.showLyricsScroll || !followLyrics.value) return
  await nextTick()
  scrollActiveLineIntoView(true)
})

watch(
  () => playQueue.nowPlaying?.id,
  () => {
    if (sheet.visible) sheet.backToCover()
  },
)

onBeforeUnmount(() => {
  document.body.style.overflow = ''
  clearIdleSnap()
})
</script>

<style scoped>
.mobile-lyrics-fullscreen {
  background: var(--cover-bg);
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

.mobile-sheet-enter-active,
.mobile-sheet-leave-active {
  transition: opacity 0.26s ease, transform 0.32s cubic-bezier(0.25, 0.8, 0.25, 1);
}
.mobile-sheet-enter-from,
.mobile-sheet-leave-to {
  opacity: 0;
  transform: translateY(100%);
}

.sheet-phase-enter-active,
.sheet-phase-leave-active {
  transition: opacity 0.18s ease;
}
.sheet-phase-enter-from,
.sheet-phase-leave-to {
  opacity: 0;
}

.mobile-lyrics-scroll {
  scrollbar-width: none;
  -ms-overflow-style: none;
}
.mobile-lyrics-scroll::-webkit-scrollbar {
  display: none;
}

.mobile-lyric-line {
  font-size: 1.05rem;
  line-height: 1.6rem;
}
.mobile-lyric-line--active {
  font-weight: 700;
  font-size: 1.18rem;
  line-height: 1.7rem;
}

.mobile-transport :deep(.play-transport--lyrics span) {
  color: var(--cover-fg-muted) !important;
}
.mobile-transport :deep(.play-transport--lyrics .n-icon) {
  color: var(--cover-control-fg) !important;
}
.mobile-transport :deep(.play-transport--lyrics .n-button) {
  outline: none !important;
  box-shadow: none !important;
}
.mobile-transport :deep(.play-transport--lyrics .n-button__border),
.mobile-transport :deep(.play-transport--lyrics .n-button__state-border) {
  display: none !important;
}
.mobile-transport :deep(.play-transport--lyrics .n-slider-rail) {
  background-color: color-mix(in srgb, var(--cover-fg) 18%, transparent) !important;
}
.mobile-transport :deep(.play-transport--lyrics .n-slider-rail__fill) {
  background-color: var(--cover-control-fg) !important;
}
</style>
