import { useRouter } from 'vue-router'
import { getSongDetail } from '@/api/netease/search'
import type { ViewMusicListItem } from '@/api/view/music'
import { neteasePlayerControl } from '@/player/neteasePlayerControl'
import { usePlayQueueStore } from '@/store/playQueue'

export type AiClientAction = {
  name: string
  args?: Record<string, any>
}

function toViewSong(detail: {
  id: number
  name: string
  dt: number
  al: { id: number; name: string; picUrl: string }
  ar: { id: number; name: string }[]
}): ViewMusicListItem {
  return {
    id: detail.id,
    name: detail.name,
    time_long: detail.dt / 1000,
    album_id: detail.al.id,
    album_name: detail.al.name,
    cover_file_url: detail.al.picUrl,
    artists: detail.ar.map((a) => ({
      artist_id: a.id,
      artist_name: a.name,
    })),
  }
}

async function loadSong(songId: number): Promise<ViewMusicListItem | null> {
  if (!songId) return null
  const detailResp = await getSongDetail([songId])
  const detail = detailResp.songs?.[0]
  if (!detail) return null
  return toViewSong(detail)
}

/** 执行后端 SSE 下发的前端副作用 tool */
export function useAiClientTools() {
  const router = useRouter()
  const playQueue = usePlayQueueStore()

  async function runClientAction(action: AiClientAction): Promise<string> {
    const name = action?.name
    const args = action?.args ?? {}
    switch (name) {
      case 'player_toggle':
        neteasePlayerControl.togglePlay()
        return '已切换播放/暂停'
      case 'player_next':
        playQueue.skipToNext()
        return '已切下一首'
      case 'queue_add': {
        const song = await loadSong(Number(args.song_id))
        if (!song) return '未找到歌曲'
        const r = playQueue.enqueue(song)
        return r === 'duplicate' ? '已在队列中' : r === 'playing' ? '已开始播放' : '已加入队列'
      }
      case 'queue_play_now': {
        const song = await loadSong(Number(args.song_id))
        if (!song) return '未找到歌曲'
        playQueue.playNow(song)
        return `正在播放：${song.name}`
      }
      case 'open_page': {
        const path = String(args.path || '')
        if (!path.startsWith('/')) return '非法路径'
        const query = (args.query && typeof args.query === 'object' ? args.query : {}) as Record<string, any>
        const q: Record<string, string> = {}
        for (const [k, v] of Object.entries(query)) {
          if (v == null) continue
          q[k] = String(v)
        }
        await router.push({ path, query: q })
        return `已打开 ${path}`
      }
      default:
        return `未知前端动作：${name}`
    }
  }

  return { runClientAction }
}
