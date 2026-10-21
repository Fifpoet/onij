/**
 * 取链：优先 GD Studio（同源 /gdstudio → music-api.gdstudio.xyz）。
 * 该站现被 Cloudflare 人机校验挡住 nginx/curl，失败则改走已有网易代理 /song/url/v1。
 * 音频文件仍由浏览器直连返回的网易 CDN，不经主站。
 */

import { getNetease } from '@/util/http'

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

interface NeteaseSongUrlBody {
  data?: Array<{ url?: string | null }>
}

function isFragilePlayUrl(url: string): boolean {
  return /\.flac(\?|$)/i.test(url)
}

/** HTTPS 页不能播 http://*.126.net（混合内容会被拦） */
function toPlayableCdnUrl(url: string): string {
  if (/^http:\/\//i.test(url) && /(?:^|\.)126\.net(?:\/|:|$)/i.test(url)) {
    return `https://${url.slice('http://'.length)}`
  }
  return url
}

function usableUrl(raw: unknown): string | null {
  if (typeof raw !== 'string' || raw.length === 0) return null
  return toPlayableCdnUrl(raw)
}

const FETCH_TIMEOUT_MS = 12_000

async function requestGdStudioPlayUrl(id: number, br: number): Promise<string | null> {
  const qs = new URLSearchParams({
    types: 'url',
    source: 'netease',
    id: String(id),
    br: String(br),
  })
  const ac = new AbortController()
  const timer = setTimeout(() => ac.abort(), FETCH_TIMEOUT_MS)
  let res: Response
  try {
    res = await fetch(`${API_ORIGIN}/api.php?${qs.toString()}`, {
      signal: ac.signal,
      cache: 'no-store',
    })
  } finally {
    clearTimeout(timer)
  }
  if (!res.ok) return null
  const text = await res.text()
  let data: GdStudioUrlBody
  try {
    data = JSON.parse(text) as GdStudioUrlBody
  } catch {
    return null
  }
  return usableUrl(data?.url)
}

async function requestNeteasePlayUrl(id: number): Promise<string | null> {
  const tryOnce = async (path: string, params: Record<string, string | number>) => {
    const data = await getNetease<NeteaseSongUrlBody>(path, params)
    return usableUrl(data?.data?.[0]?.url)
  }
  try {
    const hi = await tryOnce('/song/url/v1', { id, level: 'exhigh' })
    if (hi) return hi
  } catch {
    /* 继续试旧接口 */
  }
  try {
    return await tryOnce('/song/url', { id })
  } catch {
    return null
  }
}

export async function fetchSongPlayUrl(id: number, br: number = DEFAULT_PLAY_BR): Promise<string | null> {
  try {
    const url = await requestGdStudioPlayUrl(id, br)
    if (url && isFragilePlayUrl(url)) {
      return (await requestGdStudioPlayUrl(id, 192)) ?? url
    }
    if (url) return url
  } catch {
    /* GD Studio 被 Cloudflare 拦时改走网易代理 */
  }
  return requestNeteasePlayUrl(id)
}
