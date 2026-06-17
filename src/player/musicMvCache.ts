import { GetMusicMvBatch } from '@/api/music'
import { usePlayQueueStore } from '@/store/playQueue'

const MAX_CACHE = 96
const mvUrlByThirdId = new Map<number, string>()
const inflightBatch = new Map<string, Promise<void>>()

function cacheKey(ids: number[]): string {
  return [...ids].sort((a, b) => a - b).join(',')
}

function trimCacheIfNeeded() {
  while (mvUrlByThirdId.size > MAX_CACHE) {
    const k = mvUrlByThirdId.keys().next().value
    if (k === undefined) break
    mvUrlByThirdId.delete(k)
  }
}

function applyMvToQueue(thirdId: number, mvUrl: string) {
  const playQueue = usePlayQueueStore()
  playQueue.patchMvUrl(thirdId, mvUrl)
}

function rememberMv(thirdId: number, mvUrl: string) {
  if (!thirdId) return
  mvUrlByThirdId.set(thirdId, mvUrl)
  trimCacheIfNeeded()
  applyMvToQueue(thirdId, mvUrl)
}

async function fetchBatch(thirdIds: number[]): Promise<void> {
  const ids = [...new Set(thirdIds.filter((id) => id > 0))]
  if (!ids.length) return

  const missing = ids.filter((id) => !mvUrlByThirdId.has(id))
  if (!missing.length) return

  const key = cacheKey(missing)
  const existing = inflightBatch.get(key)
  if (existing) {
    await existing
    return
  }

  const task = (async () => {
    try {
      const resp = await GetMusicMvBatch({ third_ids: missing })
      const items = resp.items ?? []
      const returned = new Set<number>()
      for (const item of items) {
        returned.add(item.third_id)
        rememberMv(item.third_id, item.mv_url?.trim() ?? '')
      }
      for (const id of missing) {
        if (!returned.has(id)) rememberMv(id, '')
      }
    } catch {
      /* 接口未就绪时静默失败，不影响播放 */
    }
  })().finally(() => {
    inflightBatch.delete(key)
  })

  inflightBatch.set(key, task)
  await task
}

/** 后台预取 MV 地址 */
export function prefetchMusicMv(thirdId: number): void {
  if (!thirdId || mvUrlByThirdId.has(thirdId)) return
  void fetchBatch([thirdId])
}

/** 播放前确保已拉取 MV，并写入队列项 mv_url */
export async function ensureMusicMvLoaded(thirdId: number): Promise<string | null> {
  if (!thirdId) return null
  if (mvUrlByThirdId.has(thirdId)) {
    const url = mvUrlByThirdId.get(thirdId) ?? ''
    return url || null
  }
  await fetchBatch([thirdId])
  const url = mvUrlByThirdId.get(thirdId) ?? ''
  return url || null
}

export function prefetchMusicMvForQueue(thirdIds: number[]): void {
  const ids = thirdIds.filter((id) => id > 0 && !mvUrlByThirdId.has(id))
  if (!ids.length) return
  void fetchBatch(ids.slice(0, 8))
}
