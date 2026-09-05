<template>
  <div
    v-if="opened"
    class="video-host"
    role="dialog"
    aria-modal="true"
    aria-label="视频播放"
    @click.self="onClose"
  >
    <div
      ref="panelRef"
      class="video-host__panel"
      @click.stop
      @mousemove="onPanelMove"
    >
      <div v-if="status === 'loading'" class="video-host__hint">加载中…</div>
      <div v-else-if="status === 'processing'" class="video-host__hint">
        <p>{{ statusHint || '视频处理中，请稍后重试' }}</p>
        <n-button quaternary type="primary" @click="loadPlay">
          <template #icon>
            <n-icon :component="RefreshOutline" />
          </template>
          刷新
        </n-button>
      </div>
      <div v-else-if="status === 'error'" class="video-host__hint">
        <p>{{ statusHint || '无法播放' }}</p>
        <n-button quaternary type="primary" @click="loadPlay">
          <template #icon>
            <n-icon :component="RefreshOutline" />
          </template>
          重试
        </n-button>
      </div>
      <video
        v-show="status === 'ready'"
        ref="videoRef"
        class="video-host__el"
        playsinline
        @timeupdate="onTime"
        @loadedmetadata="onMeta"
        @durationchange="onTime"
        @click="togglePlay"
        @volumechange="onVolumeFromEl"
      />

      <div
        class="video-host__footer"
        :class="{ 'video-host__footer--shown': !isFullscreen || toolbarShown }"
        @mouseenter="onFooterEnter"
        @mouseleave="onFooterLeave"
      >
        <div class="video-host__row">
          <div class="video-host__cluster">
            <n-tooltip>
              <template #trigger>
                <n-button text circle :focusable="false" :aria-label="playing ? '暂停' : '播放'" @click="togglePlay">
                  <n-icon :component="playing ? PauseCircleOutline : PlayCircleOutline" :size="28" />
                </n-button>
              </template>
              {{ playing ? '暂停' : '播放' }}
            </n-tooltip>
            <span class="video-host__time">{{ clock(currentTime) }} / {{ clock(duration) }}</span>
          </div>
          <div class="video-host__cluster video-host__cluster--center">
            <n-tooltip>
              <template #trigger>
                <n-button
                  text
                  circle
                  :focusable="false"
                  :disabled="!hasPrevMarker"
                  aria-label="上一个标记"
                  @click="jumpMarker(-1)"
                >
                  <span class="video-host__chev">‹</span>
                </n-button>
              </template>
              上一个标记
            </n-tooltip>
            <n-tooltip>
              <template #trigger>
                <n-button
                  text
                  circle
                  :focusable="false"
                  :disabled="!hasNextMarker"
                  aria-label="下一个标记"
                  @click="jumpMarker(1)"
                >
                  <span class="video-host__chev">›</span>
                </n-button>
              </template>
              下一个标记
            </n-tooltip>
          </div>
          <div class="video-host__cluster video-host__cluster--end">
          <n-popover
            v-model:show="markerPopoverOpen"
            trigger="manual"
            placement="top"
            :show-arrow="true"
            raw
            class="video-marker-popover"
            @clickoutside="onMarkerClickOutside"
          >
            <template #trigger>
              <n-button
                text
                circle
                :focusable="false"
                aria-label="添加标记"
                @click="openMarkerPopover"
              >
                <n-icon :component="FlagOutline" :size="22" />
              </n-button>
            </template>
            <div class="video-marker-bubble">
              <input
                ref="markerInputRef"
                v-model="markerDraft"
                type="text"
                class="video-marker-bubble__input"
                :class="{ 'video-marker-bubble__input--error': !!markerError }"
                placeholder="标记名称"
                spellcheck="false"
                autocomplete="off"
                @keydown.enter.prevent="commitMarker"
                @blur="onMarkerBlur"
              />
              <p v-if="markerError" class="video-marker-bubble__error">{{ markerError }}</p>
            </div>
          </n-popover>
          <n-tooltip>
            <template #trigger>
              <n-button
                text
                circle
                :focusable="false"
                :disabled="!nearMarker"
                aria-label="删除标记"
                @click="removeMarker"
              >
                <n-icon :component="TrashOutline" :size="22" />
              </n-button>
            </template>
            删除标记
          </n-tooltip>
          <n-icon :component="volumeIcon" :size="20" class="video-host__vol-icon" @click="toggleMute" />
          <div class="video-host__vol">
            <n-slider
              :value="volume"
              :min="0"
              :max="100"
              :step="1"
              :tooltip="false"
              :show-tooltip="false"
              aria-label="音量"
              @update:value="onVolumeChange"
            />
          </div>
          <n-tooltip>
            <template #trigger>
              <n-button
                text
                circle
                :focusable="false"
                :aria-label="isFullscreen ? '退出全屏' : '全屏'"
                @click="toggleFullscreen"
              >
                <n-icon :component="isFullscreen ? ContractOutline : ExpandOutline" :size="22" />
              </n-button>
            </template>
            {{ isFullscreen ? '退出全屏' : '全屏' }}
          </n-tooltip>
          </div>
        </div>
        <VideoMarkerBar
          :duration="duration"
          :current-time="currentTime"
          :markers="markers"
          @seek="seekTo"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, ref, watch } from 'vue'
import { storeToRefs } from 'pinia'
import Hls from 'hls.js'
import { NButton, NIcon, NPopover, NSlider, NTooltip } from 'naive-ui'
import {
  ContractOutline,
  ExpandOutline,
  FlagOutline,
  PauseCircleOutline,
  PlayCircleOutline,
  RefreshOutline,
  TrashOutline,
  VolumeHighOutline,
  VolumeLowOutline,
  VolumeMediumOutline,
  VolumeMuteOutline,
} from '@vicons/ionicons5'
import { useVideoPlayerStore } from '@/store/videoPlayer'
import {
  DeleteFileVideoMarker,
  ListFileVideoMarkers,
  PlayFileVideo,
  SaveFileVideoMarkers,
  type FileVideoMarker,
} from '@/api/fileVideo'
import VideoMarkerBar from './VideoMarkerBar.vue'

const videoPlayer = useVideoPlayerStore()
const { opened, fileId, currentTime, duration } = storeToRefs(videoPlayer)

const videoRef = ref<HTMLVideoElement | null>(null)
const panelRef = ref<HTMLElement | null>(null)
const status = ref<'loading' | 'ready' | 'processing' | 'error'>('loading')
const statusHint = ref('')
const markers = ref<FileVideoMarker[]>([])
const playing = ref(false)
const volume = ref(80)
const muted = ref(false)
const unmuteVolume = ref(80)
const isFullscreen = ref(false)
const markerPopoverOpen = ref(false)
const markerInputRef = ref<HTMLInputElement | null>(null)
const markerDraft = ref('')
const markerError = ref('')
let markerTimeMs = 0
let skipMarkerBlur = false

let hls: Hls | null = null
let blobUrl = ''
let loadToken = 0

const NEAR_MS = 1000

const nearMarker = computed(() => {
  const t = currentTime.value * 1000
  let best: FileVideoMarker | null = null
  let bestDiff = NEAR_MS
  for (const m of markers.value) {
    const d = Math.abs(m.time_ms - t)
    if (d <= bestDiff) {
      bestDiff = d
      best = m
    }
  }
  return best
})

const volumeIcon = computed(() => {
  if (muted.value || volume.value <= 0) return VolumeMuteOutline
  if (volume.value < 34) return VolumeLowOutline
  if (volume.value < 67) return VolumeMediumOutline
  return VolumeHighOutline
})

const MARKER_EPS_MS = 500
const sortedMarkers = computed(() =>
  [...markers.value].sort((a, b) => a.time_ms - b.time_ms),
)
const hasPrevMarker = computed(() => {
  const t = currentTime.value * 1000
  return sortedMarkers.value.some((m) => m.time_ms < t - MARKER_EPS_MS)
})
const hasNextMarker = computed(() => {
  const t = currentTime.value * 1000
  return sortedMarkers.value.some((m) => m.time_ms > t + MARKER_EPS_MS)
})

const toolbarShown = ref(true)
const footerHover = ref(false)
let hideTimer: ReturnType<typeof setTimeout> | null = null

function clearHideTimer() {
  if (!hideTimer) return
  clearTimeout(hideTimer)
  hideTimer = null
}

function scheduleToolbarHide() {
  clearHideTimer()
  if (!isFullscreen.value || footerHover.value || markerPopoverOpen.value) return
  hideTimer = setTimeout(() => {
    if (!footerHover.value && !markerPopoverOpen.value && isFullscreen.value) {
      toolbarShown.value = false
    }
  }, 3000)
}

function revealToolbar() {
  toolbarShown.value = true
  scheduleToolbarHide()
}

function onPanelMove(e: MouseEvent) {
  if (!isFullscreen.value) return
  const el = panelRef.value
  if (!el) return
  const rect = el.getBoundingClientRect()
  if (e.clientY >= rect.bottom - 120) revealToolbar()
}

function onFooterEnter() {
  footerHover.value = true
  toolbarShown.value = true
  clearHideTimer()
}

function onFooterLeave() {
  footerHover.value = false
  scheduleToolbarHide()
}

function jumpMarker(dir: -1 | 1) {
  const t = currentTime.value * 1000
  const list = sortedMarkers.value
  if (dir < 0) {
    for (let i = list.length - 1; i >= 0; i--) {
      if (list[i].time_ms < t - MARKER_EPS_MS) {
        seekTo(list[i].time_ms / 1000)
        return
      }
    }
    return
  }
  for (const m of list) {
    if (m.time_ms > t + MARKER_EPS_MS) {
      seekTo(m.time_ms / 1000)
      return
    }
  }
}

function applyVolume() {
  const el = videoRef.value
  if (!el) return
  el.muted = muted.value
  el.volume = Math.min(1, Math.max(0, volume.value / 100))
}

function clock(sec: number): string {
  if (!Number.isFinite(sec) || sec < 0) return '00:00'
  const s = Math.floor(sec)
  const h = Math.floor(s / 3600)
  const m = Math.floor((s % 3600) / 60)
  const r = s % 60
  const mm = String(m).padStart(2, '0')
  const ss = String(r).padStart(2, '0')
  if (h > 0) return `${h}:${mm}:${ss}`
  return `${mm}:${ss}`
}

function destroyMedia() {
  hls?.destroy()
  hls = null
  const el = videoRef.value
  if (el) {
    el.pause()
    el.removeAttribute('src')
    el.load()
  }
  if (blobUrl) {
    URL.revokeObjectURL(blobUrl)
    blobUrl = ''
  }
  playing.value = false
}

function attachMp4(url: string) {
  const el = videoRef.value
  if (!el) return
  el.src = url
  applyVolume()
  void el.play().then(() => {
    playing.value = true
  }).catch(() => {
    playing.value = false
  })
}

function attachHls(src: string) {
  const el = videoRef.value
  if (!el) return
  if (Hls.isSupported()) {
    hls = new Hls({ enableWorker: true })
    hls.loadSource(src)
    hls.attachMedia(el)
    hls.on(Hls.Events.MANIFEST_PARSED, () => {
      applyVolume()
      void el.play().then(() => {
        playing.value = true
      }).catch(() => {
        playing.value = false
      })
    })
    hls.on(Hls.Events.ERROR, (_e, data) => {
      if (data?.fatal) {
        status.value = 'error'
        statusHint.value = 'HLS 播放失败'
      }
    })
    return
  }
  if (el.canPlayType('application/vnd.apple.mpegurl')) {
    el.src = src
    applyVolume()
    void el.play().then(() => {
      playing.value = true
    }).catch(() => {
      playing.value = false
    })
  } else {
    status.value = 'error'
    statusHint.value = '当前浏览器不支持 HLS'
  }
}

async function loadPlay() {
  const id = fileId.value
  if (!id) return
  const token = ++loadToken
  status.value = 'loading'
  statusHint.value = ''
  destroyMedia()
  try {
    const resp = await PlayFileVideo(id)
    if (token !== loadToken) return
    if (resp.mode === 'unsupported') {
      status.value = 'error'
      statusHint.value = resp.message || '仅支持 MP4'
      return
    }
    if (resp.mode === 'processing') {
      status.value = 'processing'
      statusHint.value = resp.message || '视频处理中'
      if (resp.duration_ms > 0) videoPlayer.setProgress(0, resp.duration_ms / 1000)
      return
    }
    status.value = 'ready'
    if (resp.duration_ms > 0) videoPlayer.setProgress(0, resp.duration_ms / 1000)
    await waitVideoEl()
    if (token !== loadToken) return
    if (resp.mode === 'hls' && resp.playlist) {
      const blob = new Blob([resp.playlist], { type: 'application/vnd.apple.mpegurl' })
      blobUrl = URL.createObjectURL(blob)
      attachHls(blobUrl)
    } else if (resp.mode === 'mp4' && resp.url) {
      attachMp4(resp.url)
    } else {
      status.value = 'error'
      statusHint.value = resp.message || '缺少播放地址'
    }
  } catch (e) {
    if (token !== loadToken) return
    status.value = 'error'
    statusHint.value = e instanceof Error ? e.message : '加载失败'
  }
}

function waitVideoEl(): Promise<void> {
  return nextTick()
}

async function loadMarkers() {
  const id = fileId.value
  if (!id) return
  try {
    const resp = await ListFileVideoMarkers(id)
    markers.value = resp.markers ?? []
  } catch {
    markers.value = []
  }
}

function onTime() {
  const el = videoRef.value
  if (!el) return
  playing.value = !el.paused
  const d = Number.isFinite(el.duration) && el.duration > 0 ? el.duration : duration.value
  videoPlayer.setProgress(el.currentTime, d)
}

function onMeta() {
  applyVolume()
  onTime()
}

function onVolumeChange(v: number | number[]) {
  const n = Math.round(Array.isArray(v) ? v[0] : v)
  volume.value = n
  muted.value = n <= 0
  if (n > 0) unmuteVolume.value = n
  applyVolume()
}

function onVolumeFromEl() {
  const el = videoRef.value
  if (!el) return
  muted.value = el.muted || el.volume <= 0
  if (!el.muted) volume.value = Math.round(el.volume * 100)
}

function toggleMute() {
  if (muted.value || volume.value <= 0) {
    muted.value = false
    volume.value = unmuteVolume.value > 0 ? unmuteVolume.value : 80
  } else {
    unmuteVolume.value = volume.value
    muted.value = true
  }
  applyVolume()
}

async function toggleFullscreen() {
  const el = panelRef.value
  if (!el) return
  try {
    if (document.fullscreenElement) {
      await document.exitFullscreen()
    } else {
      await el.requestFullscreen()
    }
  } catch {
    /* ignore */
  }
}

function onFullscreenChange() {
  isFullscreen.value = document.fullscreenElement === panelRef.value
  if (isFullscreen.value) revealToolbar()
  else {
    toolbarShown.value = true
    clearHideTimer()
  }
}

function togglePlay() {
  const el = videoRef.value
  if (!el || status.value !== 'ready') return
  if (el.paused) {
    void el.play().then(() => {
      playing.value = true
    }).catch(() => {})
  } else {
    el.pause()
    playing.value = false
  }
}

function seekTo(timeSec: number) {
  const el = videoRef.value
  const t = Math.max(0, timeSec)
  if (el && Number.isFinite(el.duration) && el.duration > 0) {
    el.currentTime = Math.min(t, el.duration)
  }
  videoPlayer.setProgress(t, duration.value)
}

async function openMarkerPopover() {
  markerTimeMs = Math.round(currentTime.value * 1000)
  markerDraft.value = clock(currentTime.value)
  markerError.value = ''
  markerPopoverOpen.value = true
  await nextTick()
  markerInputRef.value?.focus()
  markerInputRef.value?.select()
}

function commitMarker(): boolean {
  const label = markerDraft.value.trim()
  if (!label) {
    markerError.value = '请输入标记名称'
    return false
  }
  markerError.value = ''
  skipMarkerBlur = true
  markerPopoverOpen.value = false
  void saveNewMarker(label)
  return true
}

function onMarkerBlur() {
  if (skipMarkerBlur) {
    skipMarkerBlur = false
    return
  }
  if (!commitMarker()) markerPopoverOpen.value = true
}

function onMarkerClickOutside() {
  skipMarkerBlur = true
  if (!commitMarker()) markerPopoverOpen.value = true
}

async function saveNewMarker(label: string) {
  const id = fileId.value
  if (!id) return
  try {
    const resp = await SaveFileVideoMarkers(id, [{ time_ms: markerTimeMs, label }])
    markers.value = resp.markers ?? []
  } catch (e) {
    window.alert(e instanceof Error ? e.message : '保存失败')
  }
}

async function removeMarker() {
  const m = nearMarker.value
  if (!m) return
  try {
    await DeleteFileVideoMarker(m.id)
    markers.value = markers.value.filter((x) => x.id !== m.id)
  } catch (e) {
    window.alert(e instanceof Error ? e.message : '删除失败')
  }
}

function onClose() {
  loadToken++
  if (document.fullscreenElement) {
    void document.exitFullscreen()
  }
  destroyMedia()
  videoPlayer.close()
}

function onKey(e: KeyboardEvent) {
  if (!opened.value) return
  if (e.key === 'Escape') {
    if (markerPopoverOpen.value) {
      markerPopoverOpen.value = false
      return
    }
    if (document.fullscreenElement) return
    onClose()
  }
  if (e.key === ' ' && !(e.target instanceof HTMLInputElement)) {
    e.preventDefault()
    togglePlay()
  }
}

watch(
  () => [opened.value, fileId.value] as const,
  ([on, id]) => {
    if (!on || !id) {
      loadToken++
      destroyMedia()
      markers.value = []
      status.value = 'loading'
      markerPopoverOpen.value = false
      return
    }
    void loadPlay()
    void loadMarkers()
  },
  { immediate: true, flush: 'post' },
)

watch(opened, (on) => {
  if (on) {
    window.addEventListener('keydown', onKey)
    document.addEventListener('fullscreenchange', onFullscreenChange)
  } else {
    window.removeEventListener('keydown', onKey)
    document.removeEventListener('fullscreenchange', onFullscreenChange)
    clearHideTimer()
  }
})

watch(markerPopoverOpen, (open) => {
  if (open) {
    toolbarShown.value = true
    clearHideTimer()
  } else {
    scheduleToolbarHide()
  }
})

onBeforeUnmount(() => {
  window.removeEventListener('keydown', onKey)
  document.removeEventListener('fullscreenchange', onFullscreenChange)
  clearHideTimer()
  destroyMedia()
})
</script>

<style scoped>
.video-host {
  position: fixed;
  inset: 0;
  z-index: 90;
  background: rgb(0 0 0 / 0.86);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1rem;
}

.video-host__panel {
  position: relative;
  width: min(1100px, 100%);
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.video-host__panel:fullscreen {
  width: 100%;
  height: 100%;
  padding: 0;
  background: #000;
  box-sizing: border-box;
  overflow: hidden;
}

.video-host__el {
  width: 100%;
  max-height: min(70vh, 720px);
  background: #000;
  border-radius: 0.5rem;
}

.video-host__panel:fullscreen .video-host__el {
  flex: 1;
  max-height: none;
  min-height: 0;
  height: 100%;
  border-radius: 0;
}

.video-host__hint {
  min-height: 12rem;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 0.75rem;
  color: #fff;
  font-size: 0.95rem;
}

.video-host__footer {
  color: #fff;
}

.video-host__panel:fullscreen .video-host__footer {
  position: absolute;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 4;
  padding: 0.75rem 1.25rem 1rem;
  background: linear-gradient(transparent, rgb(0 0 0 / 0.82));
  transform: translateY(100%);
  transition: transform 0.22s ease;
}

.video-host__panel:fullscreen .video-host__footer--shown {
  transform: translateY(0);
}

.video-host__row {
  display: grid;
  grid-template-columns: 1fr auto 1fr;
  align-items: center;
  gap: 0.35rem;
  color: #fff;
}

.video-host__cluster {
  display: flex;
  align-items: center;
  gap: 0.35rem;
  min-width: 0;
}

.video-host__cluster--center {
  justify-content: center;
}

.video-host__cluster--end {
  justify-content: flex-end;
}

.video-host__row :deep(.n-button),
.video-host__row :deep(.n-icon) {
  color: #fff;
}

.video-host__row :deep(.n-button.n-button--disabled) {
  opacity: 0.35;
}

.video-host__time {
  font-size: 0.8rem;
  color: rgb(209 213 219);
}

.video-host__chev {
  font-size: 1.75rem;
  line-height: 1;
  font-weight: 600;
}

.video-host__vol-icon {
  cursor: pointer;
  flex-shrink: 0;
}

.video-host__vol {
  width: 4.5rem;
  flex-shrink: 0;
}

.video-marker-bubble {
  padding: 0.125rem;
}

.video-marker-bubble__input {
  box-sizing: border-box;
  display: block;
  width: min(72vw, 16rem);
  height: 2rem;
  padding: 0 0.875rem;
  border: 1px solid rgb(229 231 235);
  border-radius: 9999px;
  background: #fff;
  box-shadow:
    0 4px 14px rgb(0 0 0 / 0.08),
    0 1px 3px rgb(0 0 0 / 0.06);
  font-size: 0.8125rem;
  line-height: 1.25rem;
  color: rgb(31 41 55);
  outline: none;
}

.video-marker-bubble__input::placeholder {
  color: rgb(156 163 175);
}

.video-marker-bubble__input:focus {
  border-color: rgb(251 114 153);
  box-shadow:
    0 4px 14px rgb(251 114 153 / 0.15),
    0 0 0 2px rgb(251 114 153 / 0.2);
}

.video-marker-bubble__input--error {
  border-color: rgb(248 113 113);
}

.video-marker-bubble__error {
  margin: 0.375rem 0.875rem 0.125rem;
  font-size: 0.6875rem;
  line-height: 1.3;
  color: rgb(239 68 68);
}
</style>
