import { defineStore } from 'pinia'
import { ref, watch } from 'vue'
import { PostPresenceHeartbeat, type PresencePeerDTO, type PresenceTrack } from '@/api/presence'
import { randomUUID } from '@/util/crypto'
import { generatePresenceName } from '@/util/presenceName'
import { usePlayQueueStore } from '@/store/playQueue'
import type { ViewMusicListItem } from '@/api/view/music'

const HEARTBEAT_MS = 3_000
const HEARTBEAT_TIMEOUT_MS = 4_000
const QUEUE_LIMIT = 80
const PINIA_PRESENCE_KEY = 'presence'
const SESSION_CLIENT_ID_KEY = 'presence_client_id'

export type PresencePeer = {
  clientId: string
  nickname: string
  nowPlaying: PresencePeerDTO['now_playing']
  queue: PresencePeerDTO['queue']
  playbackCurrentSec: number
  isPlaying: boolean
}

function fromDTO(dto: PresencePeerDTO): PresencePeer {
  return {
    clientId: dto.client_id,
    nickname: dto.nickname,
    nowPlaying: dto.now_playing ?? null,
    queue: dto.queue ?? [],
    playbackCurrentSec: Number(dto.playback_current_sec) || 0,
    isPlaying: !!dto.is_playing,
  }
}

function slimArtists(list: ViewMusicListItem['artists'] | undefined): PresenceTrack['artists'] {
  if (!list?.length) return []
  return list.map((a) => ({
    artist_id: Number(a.artist_id) || 0,
    artist_name: String(a.artist_name || ''),
  }))
}

function slimTrack(song: ViewMusicListItem | null | undefined): PresenceTrack | null {
  if (!song) return null
  const id = Number(song.id)
  if (!Number.isFinite(id) || id <= 0) return null
  const track: PresenceTrack = {
    id,
    name: String(song.name || ''),
    artists: slimArtists(song.artists),
    time_long: Number(song.time_long) || 0,
    album_id: Number(song.album_id) || 0,
    album_name: String(song.album_name || ''),
  }
  if (song.display_artists?.length) track.display_artists = slimArtists(song.display_artists)
  if (song.cover_file_url) track.cover_file_url = song.cover_file_url
  if (song.composer_id) track.composer_id = song.composer_id
  if (song.composer_name) track.composer_name = song.composer_name
  if (song.writer_id) track.writer_id = song.writer_id
  if (song.writer_name) track.writer_name = song.writer_name
  if (song.issue_time) track.issue_time = song.issue_time
  if (song.mv_url) track.mv_url = song.mv_url
  if (song.tags?.length) track.tags = song.tags.map(String)
  return track
}

function readPersistedNickname(): string {
  try {
    const raw = localStorage.getItem(PINIA_PRESENCE_KEY)
    if (!raw) return ''
    const parsed = JSON.parse(raw) as { nickname?: string }
    return typeof parsed.nickname === 'string' ? parsed.nickname : ''
  } catch {
    return ''
  }
}

function readSessionClientId(): string {
  try {
    const v = sessionStorage.getItem(SESSION_CLIENT_ID_KEY)
    return typeof v === 'string' && v.trim() ? v : ''
  } catch {
    return ''
  }
}

function writeSessionClientId(id: string) {
  try {
    sessionStorage.setItem(SESSION_CLIENT_ID_KEY, id)
  } catch {
    /* 隐私模式受限时仍用内存里的 id */
  }
}

function startTickWorker(onTick: () => void): () => void {
  const fallback = () => {
    const id = window.setInterval(onTick, HEARTBEAT_MS)
    return () => window.clearInterval(id)
  }
  if (typeof Worker === 'undefined') return fallback()
  try {
    const worker = new Worker(new URL('./presenceTick.worker.ts', import.meta.url))
    let fallbackStop: (() => void) | null = null
    worker.onmessage = () => onTick()
    worker.onerror = () => {
      if (fallbackStop) return
      worker.terminate()
      fallbackStop = fallback()
    }
    return () => {
      worker.terminate()
      fallbackStop?.()
    }
  } catch {
    return fallback()
  }
}

export const usePresenceStore = defineStore(
  'presence',
  () => {
    const clientId = ref('')
    const nickname = ref('')
    const peers = ref<PresencePeer[]>([])

    let inflight = false
    let pending = false
    let started = false
    let lastTickAt = 0

    function ensureIdentity() {
      if (!clientId.value) {
        const sid = readSessionClientId()
        if (sid) clientId.value = sid
      }
      if (!nickname.value.trim()) {
        const saved = readPersistedNickname()
        if (saved.trim()) nickname.value = saved
      }
      if (!clientId.value) clientId.value = randomUUID()
      writeSessionClientId(clientId.value)
      if (!nickname.value.trim()) nickname.value = generatePresenceName()
    }

    function setNickname(name: string) {
      const t = name.trim()
      if (!t) return
      nickname.value = t.slice(0, 32)
      void tick()
    }

    async function tick() {
      if (inflight) {
        pending = true
        return
      }
      ensureIdentity()
      inflight = true
      try {
        const pq = usePlayQueueStore()
        const resp = await PostPresenceHeartbeat(
          {
            client_id: clientId.value,
            nickname: nickname.value,
            now_playing: slimTrack(pq.nowPlaying),
            queue: pq.queue.slice(0, QUEUE_LIMIT).map(slimTrack).filter((t): t is PresenceTrack => t != null),
            playback_current_sec: Number.isFinite(pq.playbackCurrentSec) ? pq.playbackCurrentSec : 0,
            is_playing: pq.isPlaying,
          },
          HEARTBEAT_TIMEOUT_MS,
        )
        peers.value = (resp.peers ?? []).map(fromDTO)
        lastTickAt = Date.now()
      } catch {
        /* 心跳失败保持上一份列表 */
      } finally {
        inflight = false
        if (pending) {
          pending = false
          void tick()
        }
      }
    }

    function startHeartbeat() {
      ensureIdentity()
      if (started) return
      started = true
      void tick()
      startTickWorker(() => {
        void tick()
      })
      const onVisible = () => {
        void tick()
      }
      document.addEventListener('visibilitychange', onVisible)
      window.addEventListener('focus', onVisible)
      window.addEventListener('pageshow', onVisible)
      const pq = usePlayQueueStore()
      watch(
        () => [pq.nowPlaying?.id, pq.isPlaying] as const,
        () => {
          void tick()
        },
      )
      watch(
        () => pq.playbackCurrentSec,
        () => {
          if (Date.now() - lastTickAt >= HEARTBEAT_MS) void tick()
        },
      )
    }

    function adoptPeer(peer: PresencePeer) {
      if (!peer.nowPlaying) return
      const pq = usePlayQueueStore()
      pq.replaceFromPresence({
        nowPlaying: peer.nowPlaying,
        queue: peer.queue ?? [],
        progressSec: peer.playbackCurrentSec ?? 0,
      })
    }

    function refreshPeers() {
      void tick()
    }

    return {
      clientId,
      nickname,
      peers,
      ensureIdentity,
      setNickname,
      startHeartbeat,
      adoptPeer,
      refreshPeers,
    }
  },
  {
    persist: {
      paths: ['nickname'],
    },
  },
)
