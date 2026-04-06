<template>
  <audio ref="audioRef" class="hidden" @play="onPlay" @pause="onPause" @ended="onEnded" @error="onError" />
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onBeforeUnmount, nextTick } from 'vue'
import { useMessage } from 'naive-ui'
import { usePlayQueueStore } from '@/store/playQueue'
import { fetchSongPlayUrl } from '@/api/netease/songUrl'
import { neteasePlayerControl } from '@/player/neteasePlayerControl'

const message = useMessage()
const playQueue = usePlayQueueStore()
const audioRef = ref<HTMLAudioElement | null>(null)

let loadToken = 0

function isBenignPlayRejection(err: unknown): boolean {
  if (typeof err !== 'object' || err === null || !('name' in err)) return false
  const n = (err as DOMException).name
  // 刷新后无手势：浏览器拒绝自动播放
  if (n === 'NotAllowedError') return true
  // 切歌时 src 被替换，上一次 play() 可能被中止
  if (n === 'AbortError') return true
  return false
}

function syncPlaybackFromAudio() {
  if (playQueue.playbackSeeking) return
  const el = audioRef.value
  if (!el) return
  const d = el.duration
  if (!Number.isFinite(d) || d <= 0) return
  playQueue.setPlaybackProgress(el.currentTime, d)
}

async function loadAndPlay(songId: number) {
  const token = ++loadToken
  const url = await fetchSongPlayUrl(songId)
  if (token !== loadToken) return
  let el = audioRef.value
  if (!el) {
    await nextTick()
    el = audioRef.value
  }
  if (!el) return
  if (!url) {
    message.error('无法获取播放链接')
    playQueue.advanceAfterEnded()
    return
  }
  playQueue.resetPlaybackProgress()
  el.volume = playQueue.volume / 100
  el.src = url
  el.play().catch((err: unknown) => {
    // 刷新后无用户手势时浏览器会拒绝自动播放，不应当作「播放失败」
    if (isBenignPlayRejection(err)) {
      playQueue.setPlaying(false)
      return
    }
    message.error('播放失败')
    playQueue.setPlaying(false)
  })
}

watch(
  () => playQueue.volume,
  (v) => {
    const el = audioRef.value
    if (el) el.volume = Math.min(1, Math.max(0, v / 100))
  },
)

watch(
  () => playQueue.nowPlaying,
  (song) => {
    if (!song) {
      loadToken++
      const el = audioRef.value
      if (el) {
        el.pause()
        el.removeAttribute('src')
      }
      playQueue.setPlaying(false)
      playQueue.resetPlaybackProgress()
      return
    }
    loadAndPlay(song.id)
  },
  { immediate: true },
)

function onPlay() {
  playQueue.setPlaying(true)
}

function onPause() {
  playQueue.setPlaying(false)
}

function onEnded() {
  playQueue.advanceAfterEnded()
}

function onError() {
  message.error('音频加载失败')
  playQueue.advanceAfterEnded()
}

function togglePlay() {
  const el = audioRef.value
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
  const el = audioRef.value
  if (!el) return
  const d = el.duration
  if (!Number.isFinite(d) || d <= 0) return
  const p = Math.min(100, Math.max(0, percent))
  el.currentTime = (p / 100) * d
}

function bindAudioListeners() {
  const el = audioRef.value
  if (!el) return
  el.addEventListener('timeupdate', syncPlaybackFromAudio)
  el.addEventListener('loadedmetadata', syncPlaybackFromAudio)
  el.addEventListener('durationchange', syncPlaybackFromAudio)
}

function unbindAudioListeners() {
  const el = audioRef.value
  if (!el) return
  el.removeEventListener('timeupdate', syncPlaybackFromAudio)
  el.removeEventListener('loadedmetadata', syncPlaybackFromAudio)
  el.removeEventListener('durationchange', syncPlaybackFromAudio)
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
