<template>
  <div class="ktv-pitch">
    <canvas ref="canvasRef" class="ktv-pitch__canvas" />
    <p v-if="hint" class="ktv-pitch__hint">{{ hint }}</p>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import type { SongPitch } from '@/api/uvr'
import { isNaturalC, mergePitchNotes, midiToNoteName } from '@/util/pitchNotes'

const props = defineProps<{
  pitch: SongPitch | null
  currentSec: number
  playing: boolean
  loading: boolean
  error: string
  /** 音高线与左侧标尺整体降 12 半音，伴奏不变 */
  octaveDown?: boolean
}>()

const canvasRef = ref<HTMLCanvasElement | null>(null)
let raf = 0
let stampSec = 0
let stampMs = 0

const hint = computed(() => {
  if (props.loading) return '音高提取中…'
  if (props.error) return props.error
  if (!props.pitch?.notes?.length) return '暂无音高线'
  return ''
})

const notes = computed(() => mergePitchNotes(props.pitch?.notes ?? []))

watch(
  () => props.currentSec,
  (sec) => {
    stampSec = sec
    stampMs = performance.now()
  },
  { immediate: true },
)

function displaySec() {
  if (!props.playing) return props.currentSec
  return stampSec + (performance.now() - stampMs) / 1000
}

function midiToY(midi: number, refMidi: number, height: number) {
  const span = 14
  const t = (refMidi + span / 2 - midi) / span
  return Math.min(height - 8, Math.max(8, t * height))
}

function draw() {
  const canvas = canvasRef.value
  if (!canvas) return
  const parent = canvas.parentElement
  if (!parent) return
  const dpr = window.devicePixelRatio || 1
  const w = Math.max(1, parent.clientWidth)
  const h = Math.max(1, parent.clientHeight)
  if (canvas.width !== Math.floor(w * dpr) || canvas.height !== Math.floor(h * dpr)) {
    canvas.width = Math.floor(w * dpr)
    canvas.height = Math.floor(h * dpr)
    canvas.style.width = `${w}px`
    canvas.style.height = `${h}px`
  }
  const ctx = canvas.getContext('2d')
  if (!ctx) return
  ctx.setTransform(dpr, 0, 0, dpr, 0, 0)
  ctx.clearRect(0, 0, w, h)

  const list = notes.value
  if (!list.length) return

  const now = displaySec()
  const windowSec = 8
  const labelW = 44
  const playX = Math.max(w * 0.22, labelW + 18)
  const pxPerSec = w / windowSec
  const shift = props.octaveDown ? 12 : 0
  const refMidi = Math.round(props.pitch?.ref_midi || 60) - shift
  const barH = Math.max(8, Math.min(16, h / 24))

  ctx.strokeStyle = 'rgba(255,255,255,0.08)'
  ctx.lineWidth = 1
  for (let i = -7; i <= 7; i++) {
    const y = midiToY(refMidi + i, refMidi, h)
    ctx.beginPath()
    ctx.moveTo(labelW, y)
    ctx.lineTo(w, y)
    ctx.stroke()
  }

  for (const note of list) {
    if (note.t1 < now - 1 || note.t0 > now + windowSec) continue
    const x = playX + (note.t0 - now) * pxPerSec
    const width = Math.max(4, (note.t1 - note.t0) * pxPerSec)
    const midi = note.midi - shift
    const y = midiToY(midi, refMidi, h)
    const active = now >= note.t0 && now <= note.t1
    ctx.fillStyle = active ? 'rgba(56, 189, 248, 0.95)' : 'rgba(186, 230, 253, 0.72)'
    const rr = Math.min(4, barH / 2)
    roundRect(ctx, x, y - barH / 2, width, barH, rr)
    ctx.fill()
  }

  ctx.strokeStyle = 'rgba(255,255,255,0.9)'
  ctx.lineWidth = 2
  ctx.beginPath()
  ctx.moveTo(playX, 10)
  ctx.lineTo(playX, h - 10)
  ctx.stroke()

  const fade = ctx.createLinearGradient(0, 0, labelW + 8, 0)
  fade.addColorStop(0, 'rgba(2, 6, 23, 0.88)')
  fade.addColorStop(1, 'rgba(2, 6, 23, 0)')
  ctx.fillStyle = fade
  ctx.fillRect(0, 0, labelW + 8, h)

  ctx.textAlign = 'right'
  ctx.textBaseline = 'middle'
  for (let i = -7; i <= 7; i++) {
    const midi = refMidi + i
    const y = midiToY(midi, refMidi, h)
    const name = midiToNoteName(midi)
    const isC = isNaturalC(midi)
    ctx.font = isC ? '600 12px ui-sans-serif, system-ui, sans-serif' : '500 11px ui-sans-serif, system-ui, sans-serif'
    ctx.fillStyle = isC ? 'rgba(255,255,255,0.88)' : 'rgba(226,232,240,0.48)'
    ctx.fillText(name, labelW - 6, y)
  }
}

function roundRect(
  ctx: CanvasRenderingContext2D,
  x: number,
  y: number,
  width: number,
  height: number,
  r: number,
) {
  ctx.beginPath()
  ctx.moveTo(x + r, y)
  ctx.arcTo(x + width, y, x + width, y + height, r)
  ctx.arcTo(x + width, y + height, x, y + height, r)
  ctx.arcTo(x, y + height, x, y, r)
  ctx.arcTo(x, y, x + width, y, r)
  ctx.closePath()
}

function loop() {
  draw()
  raf = requestAnimationFrame(loop)
}

onMounted(() => {
  raf = requestAnimationFrame(loop)
})

onBeforeUnmount(() => {
  cancelAnimationFrame(raf)
})
</script>

<style scoped>
.ktv-pitch {
  position: relative;
  width: min(92vw, 920px);
  height: min(52vh, 420px);
  max-height: 100%;
  border-radius: 0.85rem;
  background: rgb(2 6 23 / 0.55);
  border: 1px solid rgb(255 255 255 / 0.12);
  box-shadow: 0 20px 48px rgb(0 0 0 / 0.35);
  overflow: hidden;
}
.ktv-pitch__canvas {
  display: block;
  width: 100%;
  height: 100%;
}
.ktv-pitch__hint {
  position: absolute;
  inset: 0;
  margin: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  color: rgb(226 232 240);
  font-size: 1.05rem;
  font-weight: 600;
  text-shadow: 0 2px 12px rgb(0 0 0 / 0.6);
  pointer-events: none;
}
</style>
