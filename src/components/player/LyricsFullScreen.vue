<template>
  <Teleport to="body">
    <Transition name="lyrics-fs">
      <div
        v-if="panel.visible"
        ref="rootRef"
        class="lyrics-fullscreen fixed inset-0 z-[220] hidden flex-col overflow-hidden overflow-x-hidden text-left text-white select-none lg:flex"
        role="dialog"
        aria-modal="true"
        aria-label="歌词"
      >
        <!-- 底层全不透明，再上叠模糊封面，避免透出主页 -->
        <div class="pointer-events-none absolute inset-0 -z-10 overflow-hidden bg-zinc-950" aria-hidden="true">
          <img
            :src="coverSrc"
            alt=""
            class="absolute inset-0 h-full w-full scale-125 object-cover blur-[96px] opacity-55 saturate-150"
          >
          <div
            class="absolute inset-0 bg-gradient-to-b from-zinc-950 via-zinc-900 to-zinc-950"
          />
        </div>

        <div
          class="flex shrink-0 items-center justify-between px-3 pt-[max(0.75rem,env(safe-area-inset-top))] pb-2 sm:px-5"
        >
          <n-button
            quaternary
            circle
            size="large"
            :focusable="false"
            class="lyrics-top-btn"
            :aria-label="isFullscreen ? '退出全屏' : '全屏'"
            @click="toggleFullscreen"
          >
            <n-icon :component="isFullscreen ? ExitOutline : ExpandOutline" :size="26" />
          </n-button>
          <n-button
            quaternary
            circle
            size="large"
            :focusable="false"
            class="lyrics-top-btn"
            aria-label="关闭歌词"
            @click="closePanel"
          >
            <n-icon :component="ChevronDownOutline" :size="28" />
          </n-button>
        </div>

        <div
          class="mx-auto flex min-h-0 w-full max-w-6xl flex-1 flex-col gap-6 px-4 pb-[max(1rem,env(safe-area-inset-bottom))] pt-1 md:flex-row md:items-stretch md:justify-center md:gap-10 md:px-8 lg:gap-14"
        >
          <!-- 左栏：竖向居中整块封面+操作；不占 lyrics 的文本对齐 -->
          <aside
            class="flex w-full shrink-0 flex-col items-center justify-center py-1 md:flex md:w-[min(100%,46%)] md:max-w-[min(520px,46vw)] md:min-h-0 md:items-center md:justify-center md:py-2"
          >
            <div
              class="lyrics-left-cluster flex w-[min(100%,380px)] flex-col items-center gap-2 sm:w-[min(100%,420px)] md:w-[min(100%,480px)] md:gap-3"
            >
              <img
                :src="coverSrc"
                alt=""
                class="aspect-square w-full rounded-2xl object-cover shadow-2xl"
              >
              <div v-if="song" class="flex w-full flex-col items-center gap-1 px-0.5 text-center">
                <h2 class="text-lg font-bold leading-snug text-white md:text-xl">
                  {{ song.name }}
                </h2>
                <p
                  v-if="artistEntries.length || song.album_name"
                  class="flex flex-wrap items-center justify-center gap-x-0.5 text-sm leading-relaxed text-white/70"
                >
                  <template v-for="(a, i) in artistEntries" :key="a.artist_id">
                    <RouterLink
                      :to="{ path: '/artist', query: { ids: String(a.artist_id) } }"
                      class="lyrics-meta-link"
                    >
                      {{ a.artist_name }}
                    </RouterLink>
                    <span v-if="i < artistEntries.length - 1" class="text-white/40">/</span>
                  </template>
                  <template v-if="artistEntries.length && song.album_name">
                    <span class="mx-1 shrink-0 text-white/45">-</span>
                  </template>
                  <RouterLink
                    v-if="song.album_name"
                    :to="{ path: '/album', query: { id: String(song.album_id) } }"
                    class="lyrics-meta-link"
                  >
                    {{ song.album_name }}
                  </RouterLink>
                </p>
              </div>
              <div class="mt-1 w-full">
                <PlayTransportStrip surface="lyrics" />
              </div>
            </div>
          </aside>

          <div
            ref="lyricScrollRef"
            class="lyrics-body min-h-0 min-w-0 flex-1 self-stretch overflow-y-auto overflow-x-hidden overscroll-contain px-1 text-left md:pl-3 md:pr-2"
            @wheel="onUserLyricsInteraction"
            @touchmove.passive="onUserLyricsInteraction"
          >
            <p v-if="lyricError" class="py-8 text-left text-lg text-white/50">{{ lyricError }}</p>
            <template v-else>
              <p
                v-for="(line, i) in lines"
                :key="`${line.time}-${i}`"
                :data-lyric-i="i"
                class="lyric-line block w-full max-w-[42rem] py-2 text-left text-white/50 transition-colors duration-200"
                :class="i === activeIdx ? 'lyric-line--active text-white' : ''"
              >
                {{ line.text || '\u00a0' }}
              </p>
            </template>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onBeforeUnmount, nextTick } from 'vue'
import { RouterLink } from 'vue-router'
import { NButton, NIcon } from 'naive-ui'
import { ExpandOutline, ExitOutline, ChevronDownOutline } from '@vicons/ionicons5'
import { useLyricsPanelStore } from '@/store/lyricsPanel'
import { usePlayQueueStore } from '@/store/playQueue'
import { fetchLyricRaw } from '@/api/netease/lyric'
import { parseLrc, activeLrcIndex, type LrcLine } from '@/util/lrc'
import PlayTransportStrip from '@/components/player/PlayTransportStrip.vue'
import type { ViewMusicListItem } from '@/api/view/music'

const defaultAlbumCover =
  'https://p1.music.126.net/6y-UleORITEDbvrOLV0Q8A==/5639395138885805.jpg'

const LYRIC_IDLE_MS = 10_000

const panel = useLyricsPanelStore()
const playQueue = usePlayQueueStore()

const rootRef = ref<HTMLElement | null>(null)
const lyricScrollRef = ref<HTMLElement | null>(null)
const lines = ref<LrcLine[]>([])
const lyricError = ref('')
const isFullscreen = ref(false)

const followLyrics = ref(true)
let idleSnapTimer: ReturnType<typeof setTimeout> | null = null

const song = computed<ViewMusicListItem | null>(() => playQueue.nowPlaying)

const coverSrc = computed(
  () => playQueue.nowPlaying?.cover_file_url || defaultAlbumCover,
)

const artistEntries = computed(() => {
  const s = song.value
  if (!s) return []
  const list =
    s.display_artists && s.display_artists.length > 0 ? s.display_artists : s.artists
  return list ?? []
})

const activeIdx = computed(() => activeLrcIndex(lines.value, playQueue.playbackCurrentSec))

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

function onUserLyricsInteraction() {
  followLyrics.value = false
  scheduleIdleSnap()
}

function scrollActiveLineIntoView(smooth: boolean) {
  const wrap = lyricScrollRef.value
  const i = activeIdx.value
  if (!wrap || i < 0) return
  const el = wrap.querySelector(`[data-lyric-i="${i}"]`) as HTMLElement | null
  el?.scrollIntoView({ block: 'center', behavior: smooth ? 'smooth' : 'auto' })
}

async function loadLyrics() {
  const id = playQueue.nowPlaying?.id
  if (!id) {
    lines.value = []
    lyricError.value = ''
    return
  }
  lyricError.value = ''
  try {
    const raw = await fetchLyricRaw(id)
    lines.value = parseLrc(raw)
    if (!lines.value.length) lyricError.value = '暂无歌词'
  } catch {
    lyricError.value = '歌词加载失败'
    lines.value = []
  }
}

function onFullscreenChange() {
  isFullscreen.value = !!document.fullscreenElement
}

function toggleFullscreen() {
  const el = rootRef.value
  if (!el) return
  if (!document.fullscreenElement) {
    void el.requestFullscreen?.()
  } else {
    void document.exitFullscreen?.()
  }
}

function closePanel() {
  if (document.fullscreenElement) void document.exitFullscreen?.()
  clearIdleSnap()
  panel.close()
}

watch(
  () => panel.visible,
  (v) => {
    document.body.style.overflow = v ? 'hidden' : ''
    if (v) {
      followLyrics.value = true
      clearIdleSnap()
      void loadLyrics().then(() => {
        nextTick(() => scrollActiveLineIntoView(false))
      })
    } else {
      clearIdleSnap()
    }
  },
)

watch(
  () => playQueue.nowPlaying?.id,
  () => {
    if (panel.visible) {
      void loadLyrics().then(() => {
        nextTick(() => {
          if (followLyrics.value) scrollActiveLineIntoView(false)
        })
      })
    }
  },
)

watch(activeIdx, async () => {
  if (!followLyrics.value) return
  await nextTick()
  scrollActiveLineIntoView(true)
})

function onGlobalKeydown(e: KeyboardEvent) {
  if (!panel.visible) return
  if (e.key === 'Escape') closePanel()
}

onMounted(() => {
  document.addEventListener('fullscreenchange', onFullscreenChange)
  window.addEventListener('keydown', onGlobalKeydown)
})

onBeforeUnmount(() => {
  document.removeEventListener('fullscreenchange', onFullscreenChange)
  window.removeEventListener('keydown', onGlobalKeydown)
  document.body.style.overflow = ''
  clearIdleSnap()
})
</script>

<style scoped>
.lyrics-fs-enter-active,
.lyrics-fs-leave-active {
  transition: transform 0.38s cubic-bezier(0.25, 0.8, 0.25, 1);
}
.lyrics-fs-enter-from,
.lyrics-fs-leave-to {
  transform: translateY(100%);
}

.lyrics-top-btn {
  color: rgba(255, 255, 255, 0.92) !important;
}
.lyrics-top-btn:hover {
  color: #fff !important;
  background-color: rgba(255, 255, 255, 0.12) !important;
}
.lyrics-top-btn :deep(.n-button__border),
.lyrics-top-btn :deep(.n-button__state-border) {
  display: none !important;
}

.lyrics-meta-link {
  color: rgba(255, 255, 255, 0.82) !important;
  text-decoration: none !important;
}
.lyrics-meta-link:hover {
  color: #fff !important;
  text-decoration: underline !important;
}

/* 隐藏歌词区滚动条，避免横竖条；仍支持滚轮与触摸滑动 */
.lyrics-body {
  text-align: left;
  scrollbar-width: none;
  -ms-overflow-style: none;
}
.lyrics-body::-webkit-scrollbar {
  display: none;
  width: 0;
  height: 0;
}

.lyric-line {
  font-size: 1.125rem;
  line-height: 1.75rem;
}
@media (min-width: 768px) {
  .lyric-line {
    font-size: 1.25rem;
    line-height: 1.9rem;
  }
}
.lyric-line--active {
  font-weight: 700;
  font-size: 1.35rem;
  line-height: 1.85rem;
}
@media (min-width: 768px) {
  .lyric-line--active {
    font-size: 1.5rem;
    line-height: 2.05rem;
  }
}
</style>
