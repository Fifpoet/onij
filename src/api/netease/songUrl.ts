/**
 * 取链：网易 /song/url 无 VIP 只有试听，不用。
 * 按监控页「音乐源」顺序：GD 公开 API（无签名 types=url）→ GD 播放器
 *（types=url → autosource → kuwo/tencent 搜后 types=url）。
 * 音频仍由浏览器直连返回的 CDN。
 */

import { getSongDetail } from './search'
import {
  GdstudioSignError,
  clearGdstudioSign,
  ensureGdstudioSign,
  fetchGdstudioServerTime,
  gdstudioSignPayload,
  gdstudioUrlEncode,
} from './gdstudioSign'
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

function isInvalidSignResponse(status: number, text: string): boolean {
  return status === 401 || /invalid request/i.test(text)
}

async function signedGdJson(
  origin: string,
  params: Record<string, string>,
  signPayload: string,
  retried = false,
): Promise<unknown> {
  await ensureGdstudioSign(origin)
  const serverTime = await fetchGdstudioServerTime(origin)
  const all = { ...params, s: gdstudioSignPayload(origin, signPayload, serverTime) }
  // 值已 urlEncode，不能再走 URLSearchParams（autosource 的 JSON 会被二次编码，签名 401）
  const qs = Object.entries(all)
    .map(([key, value]) => `${key}=${value}`)
    .join('&')
  const { ok, status, text } = await fetchJson(`${origin}/api.php?${qs}`)
  if (isInvalidSignResponse(status, text) && !retried) {
    clearGdstudioSign(origin)
    try {
      await ensureGdstudioSign(origin, { force: true, remoteOnly: true })
    } catch {
      throw new GdstudioSignError()
    }
    return signedGdJson(origin, params, signPayload, true)
  }
  if (isInvalidSignResponse(status, text)) {
    throw new GdstudioSignError()
  }
  if (!ok) return null
  try {
    return JSON.parse(text) as unknown
  } catch {
    return null
  }
}

function urlFromGdBody(data: unknown): string | null {
  if (!data || typeof data !== 'object') return null
  return usableUrl((data as GdStudioUrlBody).url)
}

export async function requestGdStudioPlayUrl(
  origin: string,
  id: number,
  br: number,
  source = 'netease',
): Promise<string | null> {
  const idStr = gdstudioUrlEncode(String(id))
  const data = await signedGdJson(
    origin,
    {
      types: 'url',
      source,
      id: idStr,
      br: String(br),
    },
    idStr,
  )
  return urlFromGdBody(data)
}

/** GD 公开 API：无 crc32、不走 /time */
export async function requestGdApiPlayUrl(
  origin: string,
  id: number,
  br: number,
  source = 'netease',
): Promise<string | null> {
  const qs = new URLSearchParams({
    types: 'url',
    source,
    id: String(id),
    br: String(br),
  })
  const { ok, text } = await fetchJson(`${origin}/api.php?${qs.toString()}`)
  if (!ok) return null
  try {
    return urlFromGdBody(JSON.parse(text) as unknown)
  } catch {
    return null
  }
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
  const data = await signedGdJson(
    origin,
    {
      types: 'autosource',
      source: 'netease',
      br: String(br),
      name: payload,
    },
    payload,
  )
  const direct = urlFromGdBody(data)
  if (direct) return direct
  if (Array.isArray(data)) {
    for (const item of data) {
      const url = urlFromGdBody(item)
      if (url) return url
    }
  }
  return null
}

const GD_ALT_SOURCES = ['kuwo', 'tencent'] as const

export async function requestGdStudioAltSource(
  origin: string,
  hint: PlayUrlHint,
  br: number,
): Promise<string | null> {
  const name = hint.name?.trim()
  if (!name) return null
  const artist = hint.artist?.trim().split(',')[0]?.trim() || ''
  const wd = gdstudioUrlEncode([name, artist].filter(Boolean).join(' '))
  for (const source of GD_ALT_SOURCES) {
    const found = await signedGdJson(
      origin,
      {
        types: 'search',
        count: '5',
        source,
        pages: '1',
        name: wd,
      },
      wd,
    )
    const first = Array.isArray(found)
      ? found.find((row) => row && typeof row === 'object' && (row as { id?: unknown }).id)
      : null
    const altId = first && typeof first === 'object' ? String((first as { id?: unknown }).id ?? '') : ''
    if (!altId) continue
    const idStr = gdstudioUrlEncode(altId)
    const data = await signedGdJson(
      origin,
      { types: 'url', source, id: idStr, br: String(br) },
      idStr,
    )
    const got = urlFromGdBody(data)
    if (got) return got
  }
  return null
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
  try {
    if (source.kind === 'gdapi') return requestGdApiPlayUrl(origin, id, br)
    const gdUrl = await requestGdStudioPlayUrl(origin, id, br)
    if (gdUrl) return gdUrl
    const meta = hint?.name ? hint : await hintFromNeteaseDetail(id)
    if (!meta) return null
    return (
      (await requestGdStudioAutosource(origin, meta, br)) ||
      (await requestGdStudioAltSource(origin, meta, br))
    )
  } catch (err) {
    if (err instanceof GdstudioSignError) return null
    throw err
  }
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
