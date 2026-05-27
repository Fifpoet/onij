import { separateInstrumentalFromUrl, downloadInstrumentalBlob } from '@/api/uvr'
import { getOrFetchPlayUrl } from '@/player/songPlayUrlCache'

export type InstrumentalStatus = 'idle' | 'loading' | 'ready' | 'error'

type CacheEntry = {
  status: InstrumentalStatus
  blobUrl?: string
  error?: string
}

const MAX_CACHE = 12
const entryById = new Map<number, CacheEntry>()
const inflight = new Map<number, Promise<string | null>>()
const blobUrls = new Set<string>()

function trimCacheIfNeeded() {
  while (entryById.size > MAX_CACHE) {
    const k = entryById.keys().next().value
    if (k === undefined) break
    const e = entryById.get(k)
    if (e?.blobUrl) {
      URL.revokeObjectURL(e.blobUrl)
      blobUrls.delete(e.blobUrl)
    }
    entryById.delete(k)
  }
}

export function getInstrumentalStatus(songId: number): InstrumentalStatus {
  return entryById.get(songId)?.status ?? 'idle'
}

export function getCachedInstrumentalUrl(songId: number): string | null {
  const e = entryById.get(songId)
  if (e?.status === 'ready' && e.blobUrl) return e.blobUrl
  return null
}

async function extractInstrumental(songId: number): Promise<string | null> {
  const originalUrl = await getOrFetchPlayUrl(songId)
  if (!originalUrl) {
    entryById.set(songId, { status: 'error', error: '无法获取原音频' })
    return null
  }

  entryById.set(songId, { status: 'loading' })

  try {
    const sep = await separateInstrumentalFromUrl(originalUrl, { id: songId })
    const instBlob = await downloadInstrumentalBlob(sep.job_id)
    const blobUrl = URL.createObjectURL(instBlob)
    blobUrls.add(blobUrl)
    entryById.set(songId, { status: 'ready', blobUrl })
    trimCacheIfNeeded()
    return blobUrl
  } catch (e) {
    const msg = e instanceof Error ? e.message : '伴奏提取失败'
    entryById.set(songId, { status: 'error', error: msg })
    return null
  }
}

/** 后台预取伴奏 */
export function prefetchInstrumental(songId: number): void {
  if (!songId) return
  const hit = entryById.get(songId)
  if (hit?.status === 'ready' || hit?.status === 'loading' || inflight.has(songId)) return
  void getOrFetchInstrumental(songId)
}

export async function getOrFetchInstrumental(songId: number): Promise<string | null> {
  const cached = getCachedInstrumentalUrl(songId)
  if (cached) return cached

  const existing = inflight.get(songId)
  if (existing) return existing

  const p = extractInstrumental(songId).finally(() => {
    inflight.delete(songId)
  })
  inflight.set(songId, p)
  return p
}

export function clearInstrumentalCache() {
  for (const url of blobUrls) URL.revokeObjectURL(url)
  blobUrls.clear()
  entryById.clear()
  inflight.clear()
}
