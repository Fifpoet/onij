import { fetchSongPlayUrl } from '@/api/netease/songUrl'

const MAX_CACHE = 48
/** 网易 CDN 直链时效短；后台挂久了必须重新取链 */
const TTL_MS = 8 * 60 * 1000
/** 标签冻结时 inflight 可能永不结束，超时后视为作废 */
const INFLIGHT_MAX_MS = 15_000

type CacheEntry = { url: string; at: number }
type InflightEntry = { promise: Promise<string | null>; at: number }

const urlById = new Map<number, CacheEntry>()
const inflight = new Map<number, InflightEntry>()

function isFresh(entry: CacheEntry, now = Date.now()): boolean {
  return now - entry.at < TTL_MS
}

function usefulInflight(id: number, now = Date.now()): Promise<string | null> | null {
  const cur = inflight.get(id)
  if (!cur) return null
  if (now - cur.at >= INFLIGHT_MAX_MS) {
    inflight.delete(id)
    return null
  }
  return cur.promise
}

function trimCacheIfNeeded() {
  while (urlById.size > MAX_CACHE) {
    const k = urlById.keys().next().value
    if (k === undefined) break
    urlById.delete(k)
  }
}

function rememberUrl(id: number, url: string) {
  urlById.set(id, { url, at: Date.now() })
  trimCacheIfNeeded()
}

function startFetch(id: number): Promise<string | null> {
  const existing = usefulInflight(id)
  if (existing) return existing
  const p = (async () => {
    let url = await fetchSongPlayUrl(id)
    if (!url) url = await fetchSongPlayUrl(id)
    if (url) rememberUrl(id, url)
    return url
  })().finally(() => {
    const cur = inflight.get(id)
    if (cur?.promise === p) inflight.delete(id)
  })
  inflight.set(id, { promise: p, at: Date.now() })
  return p
}

export function forgetPlayUrl(id: number): void {
  urlById.delete(id)
  inflight.delete(id)
}

/**
 * 后台预取播放地址；新鲜缓存或进行中的请求则跳过。
 */
export function prefetchPlayUrl(id: number): void {
  if (!id) return
  const hit = urlById.get(id)
  if (hit && isFresh(hit)) return
  if (usefulInflight(id)) return
  void startFetch(id)
}

/**
 * 获取可播放 URL：优先未过期缓存，其次与进行中的预取合并，否则发起请求。
 */
export async function getOrFetchPlayUrl(id: number): Promise<string | null> {
  const hit = urlById.get(id)
  if (hit && isFresh(hit)) return hit.url
  if (hit) urlById.delete(id)
  return startFetch(id)
}

if (typeof document !== 'undefined') {
  document.addEventListener('visibilitychange', () => {
    if (document.visibilityState !== 'visible') return
    const now = Date.now()
    for (const [id, entry] of urlById) {
      if (!isFresh(entry, now)) urlById.delete(id)
    }
    for (const [id, entry] of inflight) {
      if (now - entry.at >= INFLIGHT_MAX_MS) inflight.delete(id)
    }
  })
}
