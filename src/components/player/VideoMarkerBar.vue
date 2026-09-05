<template>
  <div class="marker-bar" @mouseleave="onLeave">
    <div
      v-if="hovering"
      class="marker-bar__flyout"
      :style="{ left: `${zoomLeft}px` }"
      @mousemove.stop
    >
      <div v-if="zoomMarkers.length" class="marker-bar__chips">
        <button
          v-for="m in zoomMarkers"
          :key="`c-${m.id}`"
          type="button"
          class="marker-bar__chip"
          :title="chipTitle(m)"
          @click.stop="emit('seek', m.time_ms / 1000)"
        >
          {{ chipLabel(m) }}
        </button>
      </div>
      <div class="marker-bar__zoom" @click.stop="onClickZoom">
        <div class="marker-bar__zoom-fill" :style="{ width: `${zoomPlayedPct}%` }" />
        <div
          v-for="m in zoomMarkers"
          :key="`z-${m.id}`"
          class="marker-bar__zoom-tick"
          :style="{ left: `${zoomTickPct(m.time_ms)}%` }"
          @click.stop="emit('seek', m.time_ms / 1000)"
        />
        <div class="marker-bar__zoom-cursor" :style="{ left: `${zoomCursorPct}%` }" />
        <div class="marker-bar__zoom-label">{{ formatClock(hoverTime) }}</div>
      </div>
    </div>
    <div
      ref="trackRef"
      class="marker-bar__track"
      @mousemove="onMove"
      @click="onClickTrack"
    >
      <div class="marker-bar__fill" :style="{ width: `${playedPct}%` }" />
      <div
        v-for="m in markers"
        :key="m.id"
        class="marker-bar__tick"
        :style="{ left: `${timePct(m.time_ms)}%` }"
        :title="chipTitle(m)"
        @click.stop="emit('seek', m.time_ms / 1000)"
      />
      <div class="marker-bar__head" :style="{ left: `${playedPct}%` }" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import type { FileVideoMarker } from '@/api/fileVideo'

const props = defineProps<{
  duration: number
  currentTime: number
  markers: FileVideoMarker[]
}>()

const emit = defineEmits<{
  seek: [timeSec: number]
}>()

const trackRef = ref<HTMLElement | null>(null)
const hovering = ref(false)
const hoverRatio = ref(0)

const dur = computed(() => (props.duration > 0 ? props.duration : 1))
const playedPct = computed(() => Math.min(100, Math.max(0, (props.currentTime / dur.value) * 100)))
const hoverTime = computed(() => hoverRatio.value * dur.value)

const ZOOM_WINDOW_SEC = 40
const ZOOM_WIDTH = 360

const zoomWindow = computed(() => {
  const half = Math.max(8, Math.min(ZOOM_WINDOW_SEC / 2, dur.value * 0.05))
  const t = hoverTime.value
  let start = t - half
  let end = t + half
  if (start < 0) {
    end = Math.min(dur.value, end - start)
    start = 0
  }
  if (end > dur.value) {
    start = Math.max(0, start - (end - dur.value))
    end = dur.value
  }
  return { start, end }
})

const zoomMarkers = computed(() =>
  props.markers.filter((m) => {
    const t = m.time_ms / 1000
    return t >= zoomWindow.value.start && t <= zoomWindow.value.end
  }),
)

const zoomLeft = computed(() => {
  const el = trackRef.value
  const w = el?.clientWidth ?? 0
  const x = hoverRatio.value * w
  return Math.min(Math.max(x - ZOOM_WIDTH / 2, 0), Math.max(0, w - ZOOM_WIDTH))
})

const zoomPlayedPct = computed(() => {
  const { start, end } = zoomWindow.value
  const span = Math.max(0.001, end - start)
  return Math.min(100, Math.max(0, ((props.currentTime - start) / span) * 100))
})

const zoomCursorPct = computed(() => {
  const { start, end } = zoomWindow.value
  const span = Math.max(0.001, end - start)
  return Math.min(100, Math.max(0, ((hoverTime.value - start) / span) * 100))
})

function timePct(timeMs: number): number {
  return Math.min(100, Math.max(0, (timeMs / 1000 / dur.value) * 100))
}

function zoomTickPct(timeMs: number): number {
  const { start, end } = zoomWindow.value
  const span = Math.max(0.001, end - start)
  return Math.min(100, Math.max(0, ((timeMs / 1000 - start) / span) * 100))
}

function ratioFromEvent(e: MouseEvent): number {
  const el = trackRef.value
  if (!el) return 0
  const rect = el.getBoundingClientRect()
  if (rect.width <= 0) return 0
  return Math.min(1, Math.max(0, (e.clientX - rect.left) / rect.width))
}

function onMove(e: MouseEvent) {
  hovering.value = true
  hoverRatio.value = ratioFromEvent(e)
}

function onLeave() {
  hovering.value = false
}

function onClickTrack(e: MouseEvent) {
  emit('seek', ratioFromEvent(e) * dur.value)
}

function onClickZoom(e: MouseEvent) {
  const el = e.currentTarget as HTMLElement | null
  if (!el) return
  const rect = el.getBoundingClientRect()
  if (rect.width <= 0) return
  const local = Math.min(1, Math.max(0, (e.clientX - rect.left) / rect.width))
  const { start, end } = zoomWindow.value
  emit('seek', start + local * (end - start))
}

function formatClock(sec: number): string {
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

function chipLabel(m: FileVideoMarker): string {
  const t = formatClock(m.time_ms / 1000)
  const name = m.label?.trim()
  if (!name || name === t) return t
  return name
}

function chipTitle(m: FileVideoMarker): string {
  return formatClock(m.time_ms / 1000)
}

defineExpose({ formatClock })
</script>

<style scoped>
.marker-bar {
  position: relative;
  width: 100%;
  padding: 0.5rem 0 0.25rem;
}

.marker-bar__track {
  position: relative;
  height: 18px;
  background: rgb(55 65 81);
  border-radius: 4px;
  cursor: pointer;
}

.marker-bar__fill {
  position: absolute;
  inset: 7px auto 7px 0;
  height: 4px;
  background: rgb(59 130 246);
  border-radius: 2px;
  pointer-events: none;
}

.marker-bar__tick {
  position: absolute;
  top: 0;
  width: 3px;
  height: 100%;
  background: rgb(34 197 94);
  transform: translateX(-50%);
  z-index: 2;
}

.marker-bar__head {
  position: absolute;
  top: 2px;
  width: 8px;
  height: 14px;
  background: #fff;
  border-radius: 2px;
  transform: translateX(-50%);
  pointer-events: none;
  z-index: 3;
}

.marker-bar__flyout {
  position: absolute;
  bottom: 100%;
  width: 360px;
  z-index: 5;
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding-bottom: 4px;
}

.marker-bar__chips {
  display: flex;
  flex-wrap: nowrap;
  gap: 4px;
  overflow-x: auto;
  max-width: 100%;
  padding-bottom: 2px;
}

.marker-bar__chip {
  flex: 0 0 auto;
  max-width: 7.5rem;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  border: none;
  border-radius: 9999px;
  padding: 0.2rem 0.55rem;
  background: rgb(31 41 55);
  color: #fff;
  font-size: 0.7rem;
  cursor: pointer;
}

.marker-bar__chip:hover {
  background: rgb(59 130 246);
}

.marker-bar__zoom {
  position: relative;
  height: 40px;
  background: rgb(17 24 39);
  border: 1px solid rgb(75 85 99);
  border-radius: 6px;
  box-shadow: 0 8px 24px rgb(0 0 0 / 0.4);
  cursor: pointer;
}

.marker-bar__zoom-fill {
  position: absolute;
  inset: 16px auto 16px 0;
  height: 8px;
  background: rgb(59 130 246 / 0.5);
  pointer-events: none;
}

.marker-bar__zoom-tick {
  position: absolute;
  top: 4px;
  width: 4px;
  height: 32px;
  background: rgb(34 197 94);
  transform: translateX(-50%);
}

.marker-bar__zoom-cursor {
  position: absolute;
  top: 0;
  width: 2px;
  height: 100%;
  background: #fff;
  transform: translateX(-50%);
  pointer-events: none;
}

.marker-bar__zoom-label {
  position: absolute;
  right: 6px;
  top: 2px;
  font-size: 10px;
  color: rgb(209 213 219);
  pointer-events: none;
}
</style>
