<template>
  <audio ref="origRef" class="hidden" @play="onOrigPlay" @pause="onOrigPause" @ended="onEnded" @error="onOrigError" />
  <audio ref="instRef" class="hidden" @play="onInstPlay" @pause="onInstPause" @error="onInstError" />
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onBeforeUnmount, nextTick } from 'vue'
import { storeToRefs } from 'pinia'
import { useMessage } from 'naive-ui'
import { usePlayQueueStore } from '@/store/playQueue'
import { useKtvStore } from '@/store/ktv'
import { getOrFetchPlayUrl, prefetchPlayUrl } from '@/player/songPlayUrlCache'
import {
  ensureMusicMvLoaded,
  prefetchMusicMvForQueue,
} from '@/player/musicMvCache'
import {
  getCachedInstrumentalUrl,
  getOrFetchInstrumental,
  prefetchInstrumental,
} from '@/player/instrumentalCache'
import { neteasePlayerControl } from '@/player/neteasePlayerControl'

const message = useMessage()
const playQueue = usePlayQueueStore()
const ktv = useKtvStore()
const { active: ktvActive, accompanimentOn } = storeToRefs(ktv)
const origRef = ref<HTMLAudioElement | null>(null)
const instRef = ref<HTMLAudioElement | null>(null)

let loadToken = 0
let tailPrefetchDoneForSongId: number | null = null
/** 当前对外可听的轨：原唱 false / 伴奏 true */
let usingInstrumental = false
let origUrlForSong: string | null = null
let instUrlForSong: string | null = null

function activeEl(): HTMLAudioElement | null {
  return usingInstrumental ? instRef.value : origRef.value
}

function inactiveEl(): HTMLAudioElement | null {
  return usingInstrumental ? origRef.value : instRef.value
}

function isBenignPlayRejection(err: unknown): boolean {
  if (typeof err !== 'object' || err === null || !('name' in err)) return false
  const n = (err as DOMException).name
  if (n === 'NotAllowedError') return true
  if (n === 'AbortError') return true
  return false
}

function prefetchUpcomingFromQueue() {
  const q = playQueue.queue
  const n = Math.min(3, q.length)
  const ids: number[] = []
  for (let i = 0; i < n; i++) {
    if (ktvActive.value) prefetchInstrumental(q[i].id)
    else prefetchPlayUrl(q[i].id)
    ids.push(q[i].id)
  }
  prefetchMusicMvForQueue(ids)
}

function syncPlaybackFromAudio() {
  if (playQueue.playbackSeeking) return
  const el = activeEl()
  if (!el) return
  const d = el.duration
  if (!Number.isFinite(d) || d <= 0) return
  playQueue.setPlaybackProgress(el.currentTime, d)

  const sid = playQueue.nowPlaying?.id
  if (sid == null || playQueue.queue.length === 0) return
  if (tailPrefetchDoneForSongId === sid) return
  const remain = d - el.currentTime
  if (remain <= 10.5 && remain >= 0) {
    tailPrefetchDoneForSongId = sid
    prefetchUpcomingFromQueue()
  }
}

function resumePlay(el: HTMLAudioElement) {
  el.play().catch((err: unknown) => {
    if (isBenignPlayRejection(err)) {
      playQueue.setPlaying(false)
      return
    }
    message.error('播放失败')
    playQueue.setPlaying(false)
  })
}

function applyVolume() {
  const v = Math.min(1, Math.max(0, playQueue.volume / 100))
  if (origRef.value) origRef.value.volume = v
  if (instRef.value) instRef.value.volume = v
}

/** 将非当前轨 seek 到与当前轨相同进度（毫秒级，无重载） */
function syncInactiveToActive() {
  const active = activeEl()
  const inactive = inactiveEl()
  if (!active || !inactive?.src) return
  const t = active.currentTime
  if (!Number.isFinite(t)) return
  const d = inactive.duration
  if (Number.isFinite(d) && d > 0) {
    inactive.currentTime = Math.min(t, Math.max(0, d - 0.001))
  } else {
    inactive.currentTime = t
  }
}

/** 双轨瞬时切换：不重载、不重新 load，仅暂停一条、对齐另一条后播放 */
function switchAccompanimentTrack(useInst: boolean) {
  if (useInst === usingInstrumental) return

  const from = activeEl()
  const to = useInst ? instRef.value : origRef.value
  if (!from || !to?.src) {
    usingInstrumental = useInst
    return
  }

  const wasPlaying = !from.paused
  const t = from.currentTime

  from.pause()
  usingInstrumental = useInst

  const d = to.duration
  if (Number.isFinite(d) && d > 0) {
    to.currentTime = Math.min(t, Math.max(0, d - 0.001))
  } else {
    to.currentTime = t
  }

  playQueue.setPlaybackProgress(
    to.currentTime,
    Number.isFinite(to.duration) ? to.duration : 0,
  )

  if (wasPlaying) resumePlay(to)
}

function attachInstrumentalSrc(url: string) {
  const inst = instRef.value
  if (!inst || instUrlForSong === url) return
  instUrlForSong = url
  inst.src = url
  inst.load()
}

async function ensureInstrumentalLoaded(songId: number): Promise<boolean> {
  const url =
    getCachedInstrumentalUrl(songId) ?? (await getOrFetchInstrumental(songId))
  if (!url) return false
  attachInstrumentalSrc(url)
  return true
}

async function applyKtvSourcePreference() {
  const song = playQueue.nowPlaying
  if (!song || !origUrlForSong) return

  if (ktvActive.value && accompanimentOn.value) {
    const ok = await ensureInstrumentalLoaded(song.id)
    if (ok) switchAccompanimentTrack(true)
    return
  }

  if (usingInstrumental) {
    switchAccompanimentTrack(false)
  }
}

async function loadAndPlay(songId: number) {
  const token = ++loadToken
  origUrlForSong = null
  instUrlForSong = null
  usingInstrumental = false

  for (const el of [origRef.value, instRef.value]) {
    if (!el) continue
    el.pause()
    el.removeAttribute('src')
    el.load()
  }
  playQueue.setPlaying(false)

  const originalUrl = await getOrFetchPlayUrl(songId)
  if (token !== loadToken) return

  let orig = origRef.value
  if (!orig) {
    await nextTick()
    orig = origRef.value
  }
  if (!orig) return
  if (!originalUrl) {
    message.error('无法获取播放链接')
    playQueue.advanceAfterEnded()
    return
  }

  origUrlForSong = originalUrl
  orig.src = originalUrl

  if (ktvActive.value) {
    void ensureInstrumentalLoaded(songId)
  }

  const startWithInst =
    ktvActive.value && accompanimentOn.value && !!getCachedInstrumentalUrl(songId)

  if (startWithInst && instRef.value?.src) {
    usingInstrumental = true
    const inst = instRef.value
    playQueue.resetPlaybackProgress()
    applyVolume()
    resumePlay(inst)
    return
  }

  usingInstrumental = false
  playQueue.resetPlaybackProgress()
  applyVolume()
  resumePlay(orig)
}

function onActivePlay() {
  playQueue.setPlaying(true)
}

function onActivePause() {
  playQueue.setPlaying(false)
}

function onOrigPlay() {
  if (!usingInstrumental) onActivePlay()
  else origRef.value?.pause()
}

function onOrigPause() {
  if (!usingInstrumental) onActivePause()
}

function onInstPlay() {
  if (usingInstrumental) onActivePlay()
  else instRef.value?.pause()
}

function onInstPause() {
  if (usingInstrumental) onActivePause()
}

function onEnded() {
  playQueue.advanceAfterEnded()
}

function onOrigError() {
  if (!usingInstrumental) {
    message.error('音频加载失败')
    playQueue.advanceAfterEnded()
  }
}

function onInstError() {
  if (usingInstrumental) {
    message.error('伴奏加载失败')
    switchAccompanimentTrack(false)
  }
}

watch(
  () => playQueue.volume,
  () => applyVolume(),
)

watch(
  () => playQueue.nowPlaying,
  (song) => {
    tailPrefetchDoneForSongId = null
    if (!song) {
      loadToken++
      usingInstrumental = false
      origUrlForSong = null
      instUrlForSong = null
      for (const el of [origRef.value, instRef.value]) {
        if (!el) continue
        el.pause()
        el.removeAttribute('src')
        el.load()
      }
      playQueue.setPlaying(false)
      playQueue.resetPlaybackProgress()
      return
    }
    void ensureMusicMvLoaded(song.id)
    void loadAndPlay(song.id)
  },
  { immediate: true },
)

watch(accompanimentOn, () => {
  void applyKtvSourcePreference()
})

watch(ktvActive, (on) => {
  if (on && playQueue.nowPlaying) {
    void ensureInstrumentalLoaded(playQueue.nowPlaying.id)
  }
  void applyKtvSourcePreference()
})

// 伴奏就绪后若用户已开「伴」，自动切到预加载轨
watch(
  () => {
    const id = playQueue.nowPlaying?.id
    if (!id) return ''
    return getCachedInstrumentalUrl(id) ?? ''
  },
  (url) => {
    if (!url || !ktvActive.value || !accompanimentOn.value) return
    attachInstrumentalSrc(url)
    if (!usingInstrumental) switchAccompanimentTrack(true)
  },
)

function togglePlay() {
  const el = activeEl()
  if (!el || !playQueue.nowPlaying) return
  if (el.paused) {
    el.play().catch((err: unknown) => {
      if (isBenignPlayRejection(err)) return
      message.error('播放失败')
    })
  } else {
    el.pause()
  }
}

function seekToPercent(percent: number) {
  const el = activeEl()
  if (!el) return
  const d = el.duration
  if (!Number.isFinite(d) || d <= 0) return
  const p = Math.min(100, Math.max(0, percent))
  const t = (p / 100) * d
  el.currentTime = t
  syncInactiveToActive()
}

function bindAudioListeners() {
  for (const el of [origRef.value, instRef.value]) {
    if (!el) continue
    el.addEventListener('timeupdate', syncPlaybackFromAudio)
    el.addEventListener('loadedmetadata', syncPlaybackFromAudio)
    el.addEventListener('durationchange', syncPlaybackFromAudio)
  }
}

function unbindAudioListeners() {
  for (const el of [origRef.value, instRef.value]) {
    if (!el) continue
    el.removeEventListener('timeupdate', syncPlaybackFromAudio)
    el.removeEventListener('loadedmetadata', syncPlaybackFromAudio)
    el.removeEventListener('durationchange', syncPlaybackFromAudio)
  }
}

onMounted(() => {
  bindAudioListeners()
  neteasePlayerControl.bindTogglePlay(togglePlay)
  neteasePlayerControl.bindSeekPercent(seekToPercent)
})

onBeforeUnmount(() => {
  unbindAudioListeners()
  neteasePlayerControl.bindTogglePlay(() => {})
  neteasePlayerControl.bindSeekPercent(() => {})
})
</script>
