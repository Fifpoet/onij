import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { ViewMusicListItem } from '@/api/view/music'

export type PlayQueueMode = 'sequential' | 'shuffle'

export const usePlayQueueStore = defineStore(
  'playQueue',
  () => {
    const queue = ref<ViewMusicListItem[]>([])
    const nowPlaying = ref<ViewMusicListItem | null>(null)
    const playMode = ref<PlayQueueMode>('sequential')
    const isPlaying = ref(false)
    /** 音量 0–100，同步到 HTMLAudioElement.volume */
    const volume = ref(85)
    /** 当前播放进度（秒），由 NeteaseAudioHost 写入 */
    const playbackCurrentSec = ref(0)
    /** 当前曲目总时长（秒） */
    const playbackDurationSec = ref(0)
    /** 正在拖动进度条时暂停用 timeupdate 覆盖 UI */
    const playbackSeeking = ref(false)

    const queueLength = computed(() => queue.value.length)
    const hasCurrent = computed(() => nowPlaying.value !== null)

    function isInQueueOrPlaying(id: number): boolean {
      if (nowPlaying.value?.id === id) return true
      return queue.value.some((s) => s.id === id)
    }

    /** 加入队列；若无在播则直接开始播放该曲 */
    function enqueue(song: ViewMusicListItem): 'added' | 'playing' | 'duplicate' {
      if (isInQueueOrPlaying(song.id)) return 'duplicate'
      const item = cloneSong(song)
      if (!nowPlaying.value) {
        nowPlaying.value = item
        resetPlaybackProgress()
        return 'playing'
      }
      queue.value.push(item)
      return 'added'
    }

    /** 按顺序批量加入队列；若当前无在播则第一首立即播放 */
    function enqueueMany(songs: ViewMusicListItem[]): 'playing' | 'added' | 'empty' {
      if (songs.length === 0) return 'empty'
      let result: 'playing' | 'added' = 'added'
      for (const song of songs) {
        const r = enqueue(song)
        if (r === 'playing') result = 'playing'
      }
      return result
    }

    /** 从待播队列取下「下一首」（随机模式只随机下标，不整体重排） */
    function takeNextFromQueue(): ViewMusicListItem | null {
      if (queue.value.length === 0) return null
      if (playMode.value === 'sequential') {
        return queue.value.shift()!
      }
      const i = Math.floor(Math.random() * queue.value.length)
      return queue.value.splice(i, 1)[0]!
    }

    /** 开始播放：若当前无在播则从队列取下一首 */
    function startIfIdle() {
      if (nowPlaying.value) return
      const next = takeNextFromQueue()
      if (next) nowPlaying.value = next
    }

    /** 当前曲自然播放结束：已在播的算「播完」，再取下一段 */
    function advanceAfterEnded() {
      nowPlaying.value = null
      isPlaying.value = false
      const next = takeNextFromQueue()
      if (next) {
        nowPlaying.value = next
      }
    }

    /** 手动下一曲：从队列取下一首并切歌；队列为空则停止在播 */
    function skipToNext() {
      const next = takeNextFromQueue()
      if (next) {
        nowPlaying.value = next
        resetPlaybackProgress()
      } else {
        nowPlaying.value = null
        isPlaying.value = false
        resetPlaybackProgress()
      }
    }

    function setPlaying(playing: boolean) {
      isPlaying.value = playing
    }

    function setPlayMode(mode: PlayQueueMode) {
      playMode.value = mode
    }

    function togglePlayMode() {
      playMode.value = playMode.value === 'sequential' ? 'shuffle' : 'sequential'
    }

    function setPlaybackProgress(currentSec: number, durationSec: number) {
      playbackCurrentSec.value = Number.isFinite(currentSec) && currentSec >= 0 ? currentSec : 0
      playbackDurationSec.value =
        Number.isFinite(durationSec) && durationSec > 0 ? durationSec : 0
    }

    function setPlaybackSeeking(seeking: boolean) {
      playbackSeeking.value = seeking
    }

    function resetPlaybackProgress() {
      playbackCurrentSec.value = 0
      playbackDurationSec.value = 0
    }

    function setVolume(v: number) {
      const n = Math.round(Math.min(100, Math.max(0, v)))
      volume.value = n
    }

    function removeFromQueue(id: number) {
      queue.value = queue.value.filter((s) => s.id !== id)
    }

    /** 将待播曲目移到队首 */
    function moveToQueueTop(id: number) {
      const idx = queue.value.findIndex((s) => s.id === id)
      if (idx <= 0) return
      const [item] = queue.value.splice(idx, 1)
      queue.value.unshift(item)
    }

    function clearQueue() {
      queue.value = []
    }

    function cloneSong(s: ViewMusicListItem): ViewMusicListItem {
      return {
        ...s,
        artists: s.artists?.map((a) => ({ ...a })) ?? [],
      }
    }

    /** 按网易云歌曲 id 更新队列与在播项的 mv_url */
    function patchMvUrl(thirdId: number, mvUrl: string) {
      const url = mvUrl.trim()
      if (nowPlaying.value?.id === thirdId) {
        nowPlaying.value = {
          ...nowPlaying.value,
          mv_url: url || undefined,
        }
      }
      queue.value = queue.value.map((s) =>
        s.id === thirdId ? { ...s, mv_url: url || undefined } : s,
      )
    }

    /**
     * 插队立即播放：目标曲从队列移除并设为在播。
     * requeuePrevious 为 true 时，当前在播（若有且不同曲）插回队首。
     */
    function playNow(song: ViewMusicListItem, options?: { requeuePrevious?: boolean }) {
      const requeuePrevious = options?.requeuePrevious !== false
      const incoming = cloneSong(song)
      const prev = nowPlaying.value

      removeFromQueue(incoming.id)

      if (requeuePrevious && prev && prev.id !== incoming.id) {
        queue.value.unshift(cloneSong(prev))
      }

      nowPlaying.value = incoming
      resetPlaybackProgress()
    }

    return {
      queue,
      nowPlaying,
      playMode,
      isPlaying,
      volume,
      playbackCurrentSec,
      playbackDurationSec,
      playbackSeeking,
      queueLength,
      hasCurrent,
      isInQueueOrPlaying,
      enqueue,
      enqueueMany,
      takeNextFromQueue,
      startIfIdle,
      advanceAfterEnded,
      setPlaying,
      setPlayMode,
      togglePlayMode,
      setPlaybackProgress,
      setPlaybackSeeking,
      resetPlaybackProgress,
      setVolume,
      removeFromQueue,
      moveToQueueTop,
      clearQueue,
      playNow,
      skipToNext,
      patchMvUrl,
    }
  },
  {
    persist: {
      paths: ['queue', 'nowPlaying', 'playMode', 'volume'],
    },
  },
)
