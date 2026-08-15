import { fetchSongPitch, type SongPitch } from '@/api/uvr'
import { getOrFetchInstrumental } from '@/player/instrumentalCache'
import { mergePitchNotes } from '@/util/pitchNotes'

type PitchStatus = 'idle' | 'loading' | 'ready' | 'error'

type Entry = {
  status: PitchStatus
  data?: SongPitch
  error?: string
}

const MAX_CACHE = 12
const entryById = new Map<number, Entry>()
const inflight = new Map<number, Promise<SongPitch | null>>()

function trimCache() {
  while (entryById.size > MAX_CACHE) {
    const k = entryById.keys().next().value
    if (k === undefined) break
    entryById.delete(k)
  }
}

export function getPitchStatus(songId: number): PitchStatus {
  return entryById.get(songId)?.status ?? 'idle'
}

export function getCachedPitch(songId: number): SongPitch | null {
  const e = entryById.get(songId)
  if (e?.status === 'ready' && e.data && (e.data.version ?? 0) >= 3) return e.data
  return null
}

async function loadPitch(songId: number): Promise<SongPitch | null> {
  entryById.set(songId, { status: 'loading' })
  try {
    await getOrFetchInstrumental(songId)
    const data = await fetchSongPitch(String(songId))
    if (!data?.notes) throw new Error('音高数据为空')
    data.notes = mergePitchNotes(data.notes)
    entryById.set(songId, { status: 'ready', data })
    trimCache()
    return data
  } catch (e) {
    const msg = e instanceof Error ? e.message : '音高提取失败'
    entryById.set(songId, { status: 'error', error: msg })
    return null
  }
}

export function prefetchPitch(songId: number): void {
  if (!songId) return
  const hit = entryById.get(songId)
  if (hit?.status === 'ready' || hit?.status === 'loading' || inflight.has(songId)) return
  void getOrFetchPitch(songId)
}

export async function getOrFetchPitch(songId: number): Promise<SongPitch | null> {
  const cached = getCachedPitch(songId)
  if (cached) return cached
  const existing = inflight.get(songId)
  if (existing) return existing
  const p = loadPitch(songId).finally(() => inflight.delete(songId))
  inflight.set(songId, p)
  return p
}
