<template>
  <div ref="viewportRef" class="ktv-lyric-line">
    <div class="ktv-lyric-shift" :style="shiftStyle">
      <div class="ktv-lyric-track">
        <span class="ktv-lyric-pending">{{ text }}</span>
        <span class="ktv-lyric-active" :style="{ width: `${progress * 100}%` }">
          <span class="ktv-lyric-active__inner" :style="activeInnerStyle">{{ text }}</span>
        </span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onBeforeUnmount, nextTick } from 'vue'

const props = defineProps<{
  text: string
  progress: number
}>()

const viewportRef = ref<HTMLElement | null>(null)
const overflowPx = ref(0)

const shiftStyle = computed(() => {
  if (overflowPx.value <= 0) return undefined
  const offset = overflowPx.value * Math.min(1, Math.max(0, props.progress))
  return { transform: `translateX(-${offset}px)` }
})

const activeInnerStyle = computed(() => {
  const p = props.progress
  if (p <= 0.001) return { width: '100%' }
  return { width: `${(100 / p).toFixed(2)}%` }
})

let ro: ResizeObserver | null = null

function measure() {
  const el = viewportRef.value
  if (!el) return
  const track = el.querySelector('.ktv-lyric-track') as HTMLElement | null
  if (!track) return
  overflowPx.value = Math.max(0, track.scrollWidth - el.clientWidth)
}

watch(
  () => props.text,
  () => {
    void nextTick(measure)
  },
)

onMounted(() => {
  void nextTick(measure)
  const el = viewportRef.value
  if (el && typeof ResizeObserver !== 'undefined') {
    ro = new ResizeObserver(() => measure())
    ro.observe(el)
  }
})

onBeforeUnmount(() => {
  ro?.disconnect()
})
</script>

<style scoped>
.ktv-lyric-line {
  width: min(96vw, 52rem);
  max-width: 100%;
  overflow: hidden;
  line-height: 1.35;
  font-size: clamp(1.75rem, 7vw, 3rem);
  font-weight: 700;
  letter-spacing: 0.02em;
}

.ktv-lyric-shift {
  display: inline-block;
  min-width: 100%;
  will-change: transform;
}

.ktv-lyric-track {
  position: relative;
  display: inline-block;
  white-space: nowrap;
}

.ktv-lyric-pending {
  display: block;
  color: rgb(255 255 255 / 0.72);
  text-shadow: 0 2px 16px rgb(0 0 0 / 0.65);
  white-space: nowrap;
}

.ktv-lyric-active {
  position: absolute;
  left: 0;
  top: 0;
  height: 100%;
  overflow: hidden;
  max-width: 100%;
  pointer-events: none;
}

.ktv-lyric-active__inner {
  display: block;
  color: rgb(125 211 252);
  text-shadow: 0 2px 16px rgb(0 0 0 / 0.55);
  white-space: nowrap;
  text-align: left;
}
</style>
