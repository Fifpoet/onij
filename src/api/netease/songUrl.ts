/**
 * 取链：同源 /gdstudio → music.gdstudio.org/api.php（播放器自己的接口，带 crc32 签名）。
 * 不要走网易 /song/url（无 VIP 只有 45s 试听）。音频文件仍由浏览器直连返回的网易 CDN。
 */

import { fetchGdstudioServerTime, gdstudioSignPayload, gdstudioUrlEncode } from './gdstudioSign'

const API_ORIGIN = (
  (import.meta.env.VITE_GDSTUDIO_MUSIC_URL as string | undefined)?.replace(/\/$/, '') ||
  'https://music.gdstudio.org'
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
  const idStr = gdstudioUrlEncode(String(id))
  const serverTime = await fetchGdstudioServerTime(API_ORIGIN)
  const sign = gdstudioSignPayload(idStr, serverTime)
  const qs = new URLSearchParams({
    types: 'url',
    source: 'netease',
    id: idStr,
    br: String(br),
    s: sign,
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

export async function fetchSongPlayUrl(id: number, br: number = DEFAULT_PLAY_BR): Promise<string | null> {
  const url = await requestGdStudioPlayUrl(id, br)
  if (url && isFragilePlayUrl(url)) {
    return (await requestGdStudioPlayUrl(id, 192)) ?? url
  }
  return url
}
