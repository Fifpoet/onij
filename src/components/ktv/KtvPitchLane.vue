<template>
  <div class="ktv-pitch">
    <canvas ref="canvasRef" class="ktv-pitch__canvas" />
    <p v-if="hint" class="ktv-pitch__hint">{{ hint }}</p>
    <div v-if="!hint" class="ktv-pitch__hud">
      <div v-if="judgeText" :key="`${phraseScore}-${avgScore}`" class="ktv-pitch__judge" :class="judgeClass">
        {{ judgeText }}
        <span class="ktv-pitch__judge-num">{{ phraseScore }}</span>
      </div>
      <div class="ktv-pitch__avg" :class="scoreTone(avgScore)">
        <span class="ktv-pitch__hud-label">得分</span>
        <span class="ktv-pitch__hud-num">{{ avgScore ?? '—' }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import type { SongPitch } from '@/api/uvr'
import { isNaturalC, mergePitchNotes, midiToNoteName, pitchNoteKey, pitchScaleLabels, songPitchExtent } from '@/util/pitchNotes'

const BEAN_SMOOTH_SEC = 0.09
const LIT_HOLD_MS = 420

const props = defineProps<{
  pitch: SongPitch | null
  currentSec: number
  playing: boolean
  loading: boolean
  error: string
  userMidi?: number | null
  userVoiced?: boolean
  inTune?: boolean
  hitKeys?: string[]
  phraseScore?: number | null
  avgScore?: number | null
  micError?: string
}>()

type Spark = {
  x: number
  y: number
  vx: number
  vy: number
  life: number
  max: number
  r: number
}

const canvasRef = ref<HTMLCanvasElement | null>(null)
let raf = 0
let stampSec = 0
let stampMs = 0
let lastDrawMs = 0
let lastBeanMidi: number | null = null
let lastLitMs = 0
const sparks: Spark[] = []

const hint = computed(() => {
  if (props.loading) return '音高提取中…'
  if (props.error) return props.error
  if (props.micError) return props.micError
  if (!props.pitch?.notes?.length) return '暂无音高线'
  return ''
})

const judgeText = computed(() => {
  const s = props.phraseScore
  if (s == null) return ''
  if (s >= 90) return 'Perfect'
  if (s >= 60) return 'Good'
  return ''
})

const judgeClass = computed(() => {
  const s = props.phraseScore
  if (s == null) return ''
  if (s >= 90) return 'is-perfect'
  if (s >= 60) return 'is-good'
  return ''
})

function scoreTone(s: number | null | undefined) {
  if (s == null) return ''
  if (s >= 80) return 'is-good'
  if (s >= 50) return 'is-ok'
  return 'is-bad'
}

const notes = computed(() => mergePitchNotes(props.pitch?.notes ?? []))
const extent = computed(() => songPitchExtent(notes.value))
const hitSet = computed(() => new Set(props.hitKeys ?? []))

watch(
  () => props.pitch?.song_id,
  () => {
    lastBeanMidi = null
    lastLitMs = 0
  },
)

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

function midiToY(midi: number, lo: number, hi: number, height: number) {
  const span = Math.max(hi - lo, 7)
  const t = (hi - midi) / span
  return Math.min(height - 8, Math.max(8, t * height))
}

function spawnSparks(x: number, y: number, n: number) {
  for (let i = 0; i < n; i++) {
    sparks.push({
      x: x + (Math.random() - 0.5) * 10,
      y: y + (Math.random() - 0.5) * 12,
      vx: -30 - Math.random() * 90,
      vy: (Math.random() - 0.5) * 80,
      life: 1,
      max: 0.4 + Math.random() * 0.4,
      r: 1.3 + Math.random() * 2.4,
    })
  }
}

function stepSparks(dt: number) {
  for (let i = sparks.length - 1; i >= 0; i--) {
    const s = sparks[i]!
    s.x += s.vx * dt
    s.y += s.vy * dt
    s.vy += 30 * dt
    s.life -= dt / s.max
    if (s.life <= 0) sparks.splice(i, 1)
  }
}

function drawSparks(ctx: CanvasRenderingContext2D) {
  for (const s of sparks) {
    const a = Math.max(0, s.life)
    ctx.beginPath()
    ctx.arc(s.x, s.y, s.r, 0, Math.PI * 2)
    ctx.fillStyle = `rgba(253, 224, 71, ${0.15 + a * 0.85})`
    ctx.fill()
    ctx.beginPath()
    ctx.arc(s.x, s.y, s.r * 0.45, 0, Math.PI * 2)
    ctx.fillStyle = `rgba(255, 255, 255, ${a})`
    ctx.fill()
  }
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

  const nowMs = performance.now()
  const dt = lastDrawMs ? Math.min(0.05, (nowMs - lastDrawMs) / 1000) : 0.016
  lastDrawMs = nowMs

  const now = displaySec()
  const windowSec = 8
  const labelW = 44
  const playX = Math.max(w * 0.14, labelW + 16)
  const pxPerSec = w / windowSec
  const pastSec = Math.max(playX / pxPerSec, 2.4)
  const lo = extent.value.lo
  const hi = extent.value.hi
  const labels = pitchScaleLabels(lo, hi)
  const barH = Math.max(8, Math.min(16, h / 24))

  for (const midi of labels) {
    const y = midiToY(midi, lo, hi, h)
    ctx.strokeStyle = isNaturalC(midi) ? 'rgba(255,255,255,0.14)' : 'rgba(255,255,255,0.07)'
    ctx.lineWidth = 1
    ctx.beginPath()
    ctx.moveTo(labelW, y)
    ctx.lineTo(w, y)
    ctx.stroke()
  }

  for (const note of list) {
    if (note.t1 < now - pastSec || note.t0 > now + windowSec) continue
    const x = playX + (note.t0 - now) * pxPerSec
    const width = Math.max(4, (note.t1 - note.t0) * pxPerSec)
    const y = midiToY(note.midi, lo, hi, h)
    const active = now >= note.t0 && now <= note.t1
    const past = now > note.t1
    const hit = hitSet.value.has(pitchNoteKey(note))
    const rr = Math.min(4, barH / 2)

    if (active) {
      const passed = Math.max(0, Math.min(width, (now - note.t0) * pxPerSec))
      ctx.fillStyle = 'rgba(186, 230, 253, 0.55)'
      roundRect(ctx, x, y - barH / 2, width, barH, rr)
      ctx.fill()
      if (passed > 1) {
        ctx.fillStyle = props.inTune || hit
          ? 'rgba(253, 224, 71, 0.95)'
          : 'rgba(56, 189, 248, 0.95)'
        ctx.save()
        ctx.beginPath()
        ctx.rect(x, y - barH / 2 - 1, passed, barH + 2)
        ctx.clip()
        roundRect(ctx, x, y - barH / 2, width, barH, rr)
        ctx.fill()
        ctx.restore()
      }
    } else if (past && hit) {
      ctx.shadowColor = 'rgba(253, 224, 71, 0.65)'
      ctx.shadowBlur = 12
      ctx.fillStyle = 'rgba(253, 224, 71, 0.92)'
      roundRect(ctx, x, y - barH / 2, width, barH, rr)
      ctx.fill()
      ctx.shadowBlur = 0
    } else if (past) {
      ctx.fillStyle = 'rgba(148, 163, 184, 0.38)'
      roundRect(ctx, x, y - barH / 2, width, barH, rr)
      ctx.fill()
    } else {
      ctx.fillStyle = 'rgba(186, 230, 253, 0.72)'
      roundRect(ctx, x, y - barH / 2, width, barH, rr)
      ctx.fill()
    }
  }

  const activeNote = list.find((n) => now >= n.t0 && now <= n.t1)
  const mid = (extent.value.lo + extent.value.hi) / 2
  const onPitch = !!(props.inTune && activeNote)
  if (onPitch || props.userVoiced) lastLitMs = nowMs
  const lit = nowMs - lastLitMs < LIT_HOLD_MS

  let targetMidi: number
  if (props.userVoiced && props.userMidi != null && Number.isFinite(props.userMidi)) {
    targetMidi = props.userMidi
  } else {
    targetMidi = lastBeanMidi ?? mid
  }
  if (lastBeanMidi == null) lastBeanMidi = targetMidi
  else {
    const k = 1 - Math.exp(-dt / BEAN_SMOOTH_SEC)
    lastBeanMidi += (targetMidi - lastBeanMidi) * k
  }
  const beanY = midiToY(lastBeanMidi, lo, hi, h)

  if (onPitch) spawnSparks(playX, beanY, 2)
  stepSparks(dt)
  drawSparks(ctx)

  ctx.strokeStyle = 'rgba(255,255,255,0.9)'
  ctx.lineWidth = 2
  ctx.beginPath()
  ctx.moveTo(playX, 10)
  ctx.lineTo(playX, h - 10)
  ctx.stroke()

  const beanR = onPitch ? 9 : 7
  if (onPitch) {
    ctx.beginPath()
    ctx.arc(playX, beanY, 16, 0, Math.PI * 2)
    ctx.fillStyle = 'rgba(253, 224, 71, 0.22)'
    ctx.fill()
  }
  ctx.beginPath()
  ctx.arc(playX, beanY, beanR, 0, Math.PI * 2)
  ctx.fillStyle = onPitch
    ? 'rgba(253, 224, 71, 0.98)'
    : lit
      ? 'rgba(250, 204, 21, 0.92)'
      : 'rgba(250, 204, 21, 0.42)'
  ctx.fill()
  ctx.lineWidth = onPitch ? 2.5 : 2
  ctx.strokeStyle = onPitch
    ? 'rgba(255,255,255,1)'
    : lit
      ? 'rgba(255,255,255,0.95)'
      : 'rgba(255,255,255,0.45)'
  ctx.stroke()

  const fade = ctx.createLinearGradient(0, 0, labelW + 8, 0)
  fade.addColorStop(0, 'rgba(2, 6, 23, 0.88)')
  fade.addColorStop(1, 'rgba(2, 6, 23, 0)')
  ctx.fillStyle = fade
  ctx.fillRect(0, 0, labelW + 8, h)

  ctx.textAlign = 'right'
  ctx.textBaseline = 'middle'
  for (const midi of labels) {
    const y = midiToY(midi, lo, hi, h)
    const isC = isNaturalC(midi)
    ctx.font = isC ? '600 12px ui-sans-serif, system-ui, sans-serif' : '500 11px ui-sans-serif, system-ui, sans-serif'
    ctx.fillStyle = isC ? 'rgba(255,255,255,0.88)' : 'rgba(226,232,240,0.48)'
    ctx.fillText(midiToNoteName(midi), labelW - 6, y)
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
.ktv-pitch__hud {
  position: absolute;
  inset: 0.55rem 0.7rem auto 3.2rem;
  height: 2.2rem;
  pointer-events: none;
}
.ktv-pitch__judge,
.ktv-pitch__avg {
  position: absolute;
  top: 0;
  display: flex;
  align-items: baseline;
  gap: 0.32rem;
  padding: 0.12rem 0.55rem;
  border-radius: 0.45rem;
  background: rgb(15 23 42 / 0.55);
}
.ktv-pitch__judge {
  left: 50%;
  transform: translateX(-50%);
  font-size: 1.35rem;
  font-weight: 800;
  letter-spacing: 0.02em;
  animation: ktv-judge-in 0.28s ease-out;
}
.ktv-pitch__judge-num {
  font-size: 1rem;
  font-weight: 700;
  opacity: 0.85;
}
.ktv-pitch__judge.is-perfect { color: rgb(253 224 71); }
.ktv-pitch__judge.is-good { color: rgb(74 222 128); }
.ktv-pitch__avg {
  right: 0;
  color: rgb(226 232 240);
}
.ktv-pitch__hud-label {
  font-size: 0.72rem;
  font-weight: 600;
  opacity: 0.72;
}
.ktv-pitch__hud-num {
  font-size: 1.2rem;
  font-weight: 700;
  min-width: 1.6rem;
  text-align: right;
}
.ktv-pitch__hud .is-good { color: rgb(74 222 128); }
.ktv-pitch__hud .is-ok { color: rgb(250 204 21); }
.ktv-pitch__hud .is-bad { color: rgb(248 113 113); }
@keyframes ktv-judge-in {
  from { opacity: 0; transform: translateX(-50%) scale(0.86); }
  to { opacity: 1; transform: translateX(-50%) scale(1); }
}
</style>
