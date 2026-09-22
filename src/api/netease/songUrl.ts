/**
 * GD Studio 取链：GET /api.php（公开 API 为 music-api.gdstudio.xyz）。
 * 生产走同源 /gdstudio 反代，躲开 music.gdstudio.org 的 Cloudflare 校验。
 * 音频文件仍由浏览器直连返回的网易 CDN，不经主站。
 */

const API_ORIGIN = (
  (import.meta.env.VITE_GDSTUDIO_MUSIC_URL as string | undefined)?.replace(/\/$/, '') ||
  'https://music-api.gdstudio.xyz'
)

/** 浏览器直连网易 FLAC 常 206 断连，播放用 320k MP3 */
export const DEFAULT_PLAY_BR = 320

interface GdStudioUrlBody {
  url?: string
  br?: number
  size?: number
  from?: string
}

function isFragilePlayUrl(url: string): boolean {
  return /\.flac(\?|$)/i.test(url)
}

async function requestPlayUrl(id: number, br: number): Promise<string | null> {
  const qs = new URLSearchParams({
    types: 'url',
    source: 'netease',
    id: String(id),
    br: String(br),
  })
  const res = await fetch(`${API_ORIGIN}/api.php?${qs.toString()}`)
  if (!res.ok) return null
  const text = await res.text()
  let data: GdStudioUrlBody
  try {
    data = JSON.parse(text) as GdStudioUrlBody
  } catch {
    return null
  }
  const u = data?.url
  return typeof u === 'string' && u.length > 0 ? u : null
}

export async function fetchSongPlayUrl(id: number, br: number = DEFAULT_PLAY_BR): Promise<string | null> {
  try {
    const url = await requestPlayUrl(id, br)
    if (url && isFragilePlayUrl(url)) {
      return (await requestPlayUrl(id, 192)) ?? url
    }
    return url
  } catch {
    return null
  }
}
