import { fetchSongPlayUrl } from '@/api/netease/songUrl'

const MAX_CACHE = 48
const urlById = new Map<number, string>()
const inflight = new Map<number, Promise<string | null>>()

function trimCacheIfNeeded() {
  while (urlById.size > MAX_CACHE) {
    const k = urlById.keys().next().value
    if (k === undefined) break
    urlById.delete(k)
  }
}

function rememberUrl(id: number, url: string) {
  urlById.set(id, url)
  trimCacheIfNeeded()
}

function startFetch(id: number): Promise<string | null> {
  const existing = inflight.get(id)
  if (existing) return existing
  const p = fetchSongPlayUrl(id)
    .then((u) => {
      if (u) rememberUrl(id, u)
      return u
    })
    .finally(() => {
      inflight.delete(id)
    })
  inflight.set(id, p)
  return p
}

/**
 * 后台预取播放地址；已有缓存或已在请求中则跳过。
 */
export function prefetchPlayUrl(id: number): void {
  if (!id || urlById.has(id) || inflight.has(id)) return
  void startFetch(id)
}

/**
 * 获取可播放 URL：优先缓存，其次与进行中的预取合并，否则发起请求。
 */
export async function getOrFetchPlayUrl(id: number): Promise<string | null> {
  const hit = urlById.get(id)
  if (hit) return hit
  return startFetch(id)
}
