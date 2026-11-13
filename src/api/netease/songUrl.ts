/**
 * 取链：网易 /song/url 无 VIP 只有试听，不用。
 * 按监控页「音乐源」顺序：Meting → GD types=url → GD types=autosource。
 * 音频仍由浏览器直连返回的 CDN。
 */

import { getSongDetail } from './search'
import { fetchGdstudioServerTime, gdstudioSignPayload, gdstudioUrlEncode } from './gdstudioSign'
import { resolvePlaySourceOrigin, usePlaySourcesStore, type PlaySource } from '@/store/playSources'

/** 浏览器直连网易 FLAC 常 206 断连，播放用 320k MP3 */
export const DEFAULT_PLAY_BR = 320

export interface PlayUrlHint {
  name?: string
  artist?: string
  album?: string
  duration?: number
}

interface GdStudioUrlBody {
  url?: string
  br?: number
  size?: number
  from?: string
}

function isFragilePlayUrl(url: string): boolean {
  return /\.flac(\?|$)/i.test(url)
}

/** HTTPS 页不能播 http 音频（混合内容会被拦） */
function toPlayableCdnUrl(url: string): string {
  if (/^http:\/\//i.test(url)) {
    return `https://${url.slice('http://'.length)}`
  }
  return url
}

function usableUrl(raw: unknown): string | null {
  if (typeof raw !== 'string' || raw.length === 0) return null
  if (!/^https?:\/\//i.test(raw)) return null
  return toPlayableCdnUrl(raw)
}

const FETCH_TIMEOUT_MS = 12_000

async function fetchJson(url: string): Promise<{ ok: boolean; status: number; text: string }> {
  const ac = new AbortController()
  const timer = setTimeout(() => ac.abort(), FETCH_TIMEOUT_MS)
  try {
    const res = await fetch(url, { signal: ac.signal, cache: 'no-store' })
    return { ok: res.ok, status: res.status, text: await res.text() }
  } finally {
    clearTimeout(timer)
  }
}

async function signedGdQuery(
  origin: string,
  params: Record<string, string>,
  signPayload: string,
): Promise<GdStudioUrlBody | null> {
  const serverTime = await fetchGdstudioServerTime(origin)
  const all = { ...params, s: gdstudioSignPayload(signPayload, serverTime) }
  // 值已 urlEncode，不能再走 URLSearchParams（autosource 的 JSON 会被二次编码，签名 401）
  const qs = Object.entries(all)
    .map(([key, value]) => `${key}=${value}`)
    .join('&')
  const { ok, text } = await fetchJson(`${origin}/api.php?${qs}`)
  if (!ok) return null
  try {
    return JSON.parse(text) as GdStudioUrlBody
  } catch {
    return null
  }
}

async function requestGdStudioPlayUrl(origin: string, id: number, br: number): Promise<string | null> {
  const idStr = gdstudioUrlEncode(String(id))
  const data = await signedGdQuery(
    origin,
    {
      types: 'url',
      source: 'netease',
      id: idStr,
      br: String(br),
    },
    idStr,
  )
  return usableUrl(data?.url)
}

async function requestGdStudioAutosource(
  origin: string,
  hint: PlayUrlHint,
  br: number,
): Promise<string | null> {
  const name = hint.name?.trim()
  if (!name) return null
  const payload = gdstudioUrlEncode(
    JSON.stringify({
      name,
      artist: hint.artist?.trim() || '',
      album: hint.album?.trim() || '',
      isrc: '',
      duration: Math.round(hint.duration || 0),
    }),
  )
  const data = await signedGdQuery(
    origin,
    {
      types: 'autosource',
      source: 'netease',
      br: String(br),
      name: payload,
    },
    payload,
  )
  return usableUrl(data?.url)
}

async function requestMetingPlayUrl(origin: string, id: number, br: number): Promise<string | null> {
  const qs = new URLSearchParams({
    server: 'netease',
    type: 'url',
    id: String(id),
    br: String(br),
  })
  const ac = new AbortController()
  const timer = setTimeout(() => ac.abort(), FETCH_TIMEOUT_MS)
  let res: Response
  try {
    res = await fetch(`${origin}/?${qs.toString()}`, {
      signal: ac.signal,
      cache: 'no-store',
      redirect: 'manual',
    })
  } catch {
    return null
  } finally {
    clearTimeout(timer)
  }
  const loc = res.headers.get('Location')
  if (loc) return usableUrl(loc)
  if (!res.ok) return null
  const text = (await res.text()).trim()
  if (!text) return null
  try {
    const data = JSON.parse(text) as { url?: string }
    return usableUrl(data?.url)
  } catch {
    return usableUrl(text)
  }
}

async function hintFromNeteaseDetail(id: number): Promise<PlayUrlHint | null> {
  try {
    const detail = await getSongDetail([id])
    const song = detail?.songs?.[0]
    if (!song?.name) return null
    return {
      name: song.name,
      artist: song.ar?.[0]?.name || '',
      album: song.al?.name || '',
      duration: song.dt > 0 ? song.dt / 1000 : 0,
    }
  } catch {
    return null
  }
}

async function requestFromSource(
  source: PlaySource,
  id: number,
  br: number,
  hint?: PlayUrlHint | null,
): Promise<string | null> {
  const origin = resolvePlaySourceOrigin(source)
  if (source.kind === 'meting') return requestMetingPlayUrl(origin, id, br)
  const gdUrl = await requestGdStudioPlayUrl(origin, id, br)
  if (gdUrl) return gdUrl
  const meta = hint?.name ? hint : await hintFromNeteaseDetail(id)
  if (meta) return requestGdStudioAutosource(origin, meta, br)
  return null
}

async function requestWithBr(
  id: number,
  br: number,
  hint?: PlayUrlHint | null,
): Promise<string | null> {
  const sources = usePlaySourcesStore().enabledSources
  for (const source of sources) {
    const url = await requestFromSource(source, id, br, hint)
    if (url) return url
  }
  return null
}

export async function fetchSongPlayUrl(
  id: number,
  br: number = DEFAULT_PLAY_BR,
  hint?: PlayUrlHint,
): Promise<string | null> {
  const url = await requestWithBr(id, br, hint)
  if (url && isFragilePlayUrl(url)) {
    return (await requestWithBr(id, 192, hint)) ?? url
  }
  return url
}
