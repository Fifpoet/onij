import { ref, computed, watch } from 'vue'
import { usePlayQueueStore } from '@/store/playQueue'
import { fetchLyricRaw } from '@/api/netease/lyric'
import { parseLrc, activeLrcIndex, type LrcLine } from '@/util/lrc'

const lines = ref<LrcLine[]>([])
const lyricError = ref('')
const loading = ref(false)
let loadToken = 0

async function loadLyricsForSong(songId: number | undefined) {
  const token = ++loadToken
  if (!songId) {
    lines.value = []
    lyricError.value = ''
    loading.value = false
    return
  }
  loading.value = true
  lyricError.value = ''
  try {
    const raw = await fetchLyricRaw(songId)
    if (token !== loadToken) return
    lines.value = parseLrc(raw)
    if (!lines.value.length) lyricError.value = '暂无歌词'
  } catch {
    if (token !== loadToken) return
    lyricError.value = '歌词加载失败'
    lines.value = []
  } finally {
    if (token === loadToken) loading.value = false
  }
}

/** 顶栏底栏 / 手机歌词页共用歌词数据 */
export function usePlayerLyrics() {
  const playQueue = usePlayQueueStore()

  const activeIdx = computed(() => activeLrcIndex(lines.value, playQueue.playbackCurrentSec))

  const currentLineText = computed(() => {
    const song = playQueue.nowPlaying
    if (!song) return '暂无播放'
    if (loading.value) return '歌词加载中…'
    if (lyricError.value) return lyricError.value
    const idx = activeIdx.value
    if (idx >= 0 && lines.value[idx]?.text) return lines.value[idx].text
    if (lines.value.length > 0) return lines.value[0].text
    return song.name
  })

  watch(
    () => playQueue.nowPlaying?.id,
    (id) => {
      void loadLyricsForSong(id)
    },
    { immediate: true },
  )

  return {
    lines,
    lyricError,
    loading,
    activeIdx,
    currentLineText,
    reloadLyrics: () => loadLyricsForSong(playQueue.nowPlaying?.id),
  }
}
