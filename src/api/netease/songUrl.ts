/**
 * GD Studio 播放直链：https://music.gdstudio.org/api.php
 * 官网播放器为 POST + /time 签名（s = md5(time[:9]|host|version|id) 后 8 位大写）。
 */

const API_ORIGIN = (
  (import.meta.env.VITE_GDSTUDIO_MUSIC_URL as string | undefined)?.replace(/\/$/, '') ||
  'https://music.gdstudio.org'
)
const API_HOST = 'music.gdstudio.org'
const PLAYER_VERSION = '2026.08.01'

/** 浏览器直连网易 FLAC 常 206 断连，播放用 320k MP3 */
export const DEFAULT_PLAY_BR = 320

function normalizeVersion(version: string): string {
  return version
    .split('.')
    .map((part) => (part.length === 1 ? part.padStart(2, '0') : part))
    .join('')
}

function md5Hex(message: string): string {
  const bytes = unescape(encodeURIComponent(message))
  const n = bytes.length
  const words: number[] = []
  for (let i = 0; i < n; i += 1) {
    words[i >> 2] |= bytes.charCodeAt(i) << ((i % 4) * 8)
  }
  const bitLen = n * 8
  words[n >> 2] |= 0x80 << ((n % 4) * 8)
  const lenIndex = (((n + 8) >> 6) + 1) * 16 - 2
  words[lenIndex] = bitLen >>> 0
  words[lenIndex + 1] = 0
  while (words.length % 16 !== 0) words.push(0)

  let a = 0x67452301
  let b = 0xefcdab89
  let c = 0x98badcfe
  let d = 0x10325476

  const add = (x: number, y: number) => (x + y) >>> 0
  const rotl = (x: number, s: number) => (x << s) | (x >>> (32 - s))
  const cmn = (q: number, a0: number, b0: number, x: number, s: number, t: number) =>
    add(rotl(add(add(a0, q), add(x, t)), s), b0)
  const ff = (a0: number, b0: number, c0: number, d0: number, x: number, s: number, t: number) =>
    cmn((b0 & c0) | (~b0 & d0), a0, b0, x, s, t)
  const gg = (a0: number, b0: number, c0: number, d0: number, x: number, s: number, t: number) =>
    cmn((b0 & d0) | (c0 & ~d0), a0, b0, x, s, t)
  const hh = (a0: number, b0: number, c0: number, d0: number, x: number, s: number, t: number) =>
    cmn(b0 ^ c0 ^ d0, a0, b0, x, s, t)
  const ii = (a0: number, b0: number, c0: number, d0: number, x: number, s: number, t: number) =>
    cmn(c0 ^ (b0 | ~d0), a0, b0, x, s, t)

  for (let i = 0; i < words.length; i += 16) {
    const oa = a
    const ob = b
    const oc = c
    const od = d
    const w = (j: number) => words[i + j] || 0

    a = ff(a, b, c, d, w(0), 7, 0xd76aa478)
    d = ff(d, a, b, c, w(1), 12, 0xe8c7b756)
    c = ff(c, d, a, b, w(2), 17, 0x242070db)
    b = ff(b, c, d, a, w(3), 22, 0xc1bdceee)
    a = ff(a, b, c, d, w(4), 7, 0xf57c0faf)
    d = ff(d, a, b, c, w(5), 12, 0x4787c62a)
    c = ff(c, d, a, b, w(6), 17, 0xa8304613)
    b = ff(b, c, d, a, w(7), 22, 0xfd469501)
    a = ff(a, b, c, d, w(8), 7, 0x698098d8)
    d = ff(d, a, b, c, w(9), 12, 0x8b44f7af)
    c = ff(c, d, a, b, w(10), 17, 0xffff5bb1)
    b = ff(b, c, d, a, w(11), 22, 0x895cd7be)
    a = ff(a, b, c, d, w(12), 7, 0x6b901122)
    d = ff(d, a, b, c, w(13), 12, 0xfd987193)
    c = ff(c, d, a, b, w(14), 17, 0xa679438e)
    b = ff(b, c, d, a, w(15), 22, 0x49b40821)

    a = gg(a, b, c, d, w(1), 5, 0xf61e2562)
    d = gg(d, a, b, c, w(6), 9, 0xc040b340)
    c = gg(c, d, a, b, w(11), 14, 0x265e5a51)
    b = gg(b, c, d, a, w(0), 20, 0xe9b6c7aa)
    a = gg(a, b, c, d, w(5), 5, 0xd62f105d)
    d = gg(d, a, b, c, w(10), 9, 0x02441453)
    c = gg(c, d, a, b, w(15), 14, 0xd8a1e681)
    b = gg(b, c, d, a, w(4), 20, 0xe7d3fbc8)
    a = gg(a, b, c, d, w(9), 5, 0x21e1cde6)
    d = gg(d, a, b, c, w(14), 9, 0xc33707d6)
    c = gg(c, d, a, b, w(3), 14, 0xf4d50d87)
    b = gg(b, c, d, a, w(8), 20, 0x455a14ed)
    a = gg(a, b, c, d, w(13), 5, 0xa9e3e905)
    d = gg(d, a, b, c, w(2), 9, 0xfcefa3f8)
    c = gg(c, d, a, b, w(7), 14, 0x676f02d9)
    b = gg(b, c, d, a, w(12), 20, 0x8d2a4d8a)

    a = hh(a, b, c, d, w(5), 4, 0xfffa3942)
    d = hh(d, a, b, c, w(8), 11, 0x8771f681)
    c = hh(c, d, a, b, w(11), 16, 0x6d9d6122)
    b = hh(b, c, d, a, w(14), 23, 0xfde5380c)
    a = hh(a, b, c, d, w(1), 4, 0xa4beea44)
    d = hh(d, a, b, c, w(4), 11, 0x4bdecfa9)
    c = hh(c, d, a, b, w(7), 16, 0xf6bb4b60)
    b = hh(b, c, d, a, w(10), 23, 0xbebfbc70)
    a = hh(a, b, c, d, w(13), 4, 0x289b7ec6)
    d = hh(d, a, b, c, w(0), 11, 0xeaa127fa)
    c = hh(c, d, a, b, w(3), 16, 0xd4ef3085)
    b = hh(b, c, d, a, w(6), 23, 0x04881d05)
    a = hh(a, b, c, d, w(9), 4, 0xd9d4d039)
    d = hh(d, a, b, c, w(12), 11, 0xe6db99e5)
    c = hh(c, d, a, b, w(15), 16, 0x1fa27cf8)
    b = hh(b, c, d, a, w(2), 23, 0xc4ac5665)

    a = ii(a, b, c, d, w(0), 6, 0xf4292244)
    d = ii(d, a, b, c, w(7), 10, 0x432aff97)
    c = ii(c, d, a, b, w(14), 15, 0xab9423a7)
    b = ii(b, c, d, a, w(5), 21, 0xfc93a039)
    a = ii(a, b, c, d, w(12), 6, 0x655b59c3)
    d = ii(d, a, b, c, w(3), 10, 0x8f0ccc92)
    c = ii(c, d, a, b, w(10), 15, 0xffeff47d)
    b = ii(b, c, d, a, w(1), 21, 0x85845dd1)
    a = ii(a, b, c, d, w(8), 6, 0x6fa87e4f)
    d = ii(d, a, b, c, w(15), 10, 0xfe2ce6e0)
    c = ii(c, d, a, b, w(6), 15, 0xa3014314)
    b = ii(b, c, d, a, w(13), 21, 0x4e0811a1)
    a = ii(a, b, c, d, w(4), 6, 0xf7537e82)
    d = ii(d, a, b, c, w(11), 10, 0xbd3af235)
    c = ii(c, d, a, b, w(2), 15, 0x2ad7d2bb)
    b = ii(b, c, d, a, w(9), 21, 0xeb86d391)

    a = add(a, oa)
    b = add(b, ob)
    c = add(c, oc)
    d = add(d, od)
  }

  const hex = (x: number) => {
    let s = ''
    for (let i = 0; i < 4; i += 1) {
      s += ((x >> (i * 8)) & 0xff).toString(16).padStart(2, '0')
    }
    return s
  }
  return hex(a) + hex(b) + hex(c) + hex(d)
}

function makeSign(payload: string, serverTime: string): string {
  const timePrefix = String(serverTime).slice(0, 9)
  const raw = `${timePrefix}|${API_HOST}|${normalizeVersion(PLAYER_VERSION)}|${payload}`
  return md5Hex(raw).slice(-8).toUpperCase()
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
function isFragilePlayUrl(url: string): boolean {
  return /\.flac(\?|$)/i.test(url)
}

async function requestPlayUrl(id: number, br: number): Promise<string | null> {
  const timeRes = await fetch(`${API_ORIGIN}/time`)
  if (!timeRes.ok) return null
  const serverTime = (await timeRes.text()).trim()
  if (!serverTime) return null

  const sid = String(id)
  const body = new URLSearchParams({
    types: 'url',
    id: sid,
    source: 'netease',
    br: String(br),
    s: makeSign(sid, serverTime),
  })
  const res = await fetch(`${API_ORIGIN}/api.php`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded; charset=UTF-8',
      'X-Requested-With': 'XMLHttpRequest',
      Accept: 'application/json, text/javascript, */*; q=0.01',
    },
    body,
  })
  if (!res.ok) return null
  const data = (await res.json()) as GdStudioUrlBody
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
