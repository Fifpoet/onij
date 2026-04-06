/**
 * GD Studio 音乐文件地址（与 id、br 替换）
 * 文档形态：https://music-api.gdstudio.xyz/api.php?types=url&id=287671&br=999
 */

const DEFAULT_BASE = 'https://music-api.gdstudio.xyz'
/** 与示例接口一致，可按需改为 320 等 */
export const DEFAULT_PLAY_BR = 999

function buildRequestUrl(id: number, br: number): string {
  const q = new URLSearchParams({
    types: 'url',
    id: String(id),
    br: String(br),
  })
  const path = `/api.php?${q.toString()}`
  if (import.meta.env.DEV) {
    return `/gdstudio-music${path}`
  }
  const base = (import.meta.env.VITE_GDSTUDIO_MUSIC_URL as string | undefined)?.replace(/\/$/, '') || DEFAULT_BASE
  return `${base}${path}`
}

interface GdStudioUrlBody {
  url?: string
  br?: number
  size?: number
  from?: string
}

/**
 * 获取单曲可播放 URL
 */
export async function fetchSongPlayUrl(id: number, br: number = DEFAULT_PLAY_BR): Promise<string | null> {
  try {
    const res = await fetch(buildRequestUrl(id, br))
    if (!res.ok) return null
    const body = (await res.json()) as GdStudioUrlBody
    const u = body?.url
    return typeof u === 'string' && u.length > 0 ? u : null
  } catch {
    return null
  }
}
