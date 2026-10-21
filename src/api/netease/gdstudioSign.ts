/**
 * GD Studio 播放器 (music.gdstudio.org) 的取链签名。
 * 他们把函数叫 crc32，实际是：time[:9]|hostname|version去点|payload → 自定义 MD5 末 8 位大写。
 * 标准 MD5 会对 401；必须跑播放器自带的那份哈希。
 */
import vendorCode from '@/vendor/gdstudio-crc32.min.js?raw'

/** 与 js/player.js 的 mkPlayer.version 对齐；签名里会去点变成 20260925 */
export const GDSTUDIO_PLAYER_VERSION = '2026.09.25'

/** nginx / Vite 反代时 Host 也是这个，签名必须用同一 hostname */
export const GDSTUDIO_SIGN_HOST = 'music.gdstudio.org'

let cachedMd5: ((input: string) => string) | null = null

function loadGdstudioMd5(): (input: string) => string {
  if (cachedMd5) return cachedMd5
  const sandbox: Record<string, unknown> = {
    mkPlayer: { version: GDSTUDIO_PLAYER_VERSION },
    XMLHttpRequest: class {
      status = 0
      readyState = 0
      responseText = ''
      open() {}
      send() {
        this.status = 200
        this.readyState = 4
        this.responseText = '0'
      }
      abort() {}
    },
    setTimeout,
    clearTimeout,
    encodeURIComponent,
    unescape: (s: string) =>
      s.replace(/%([0-9A-Fa-f]{2})/g, (_, hex: string) =>
        String.fromCharCode(Number.parseInt(hex, 16)),
      ),
    decodeURIComponent,
    Number,
    String,
    Date,
    parseInt,
    Math,
    JSON,
  }
  sandbox.window = sandbox
  sandbox.self = sandbox
  sandbox.location = {
    hostname: GDSTUDIO_SIGN_HOST,
    host: GDSTUDIO_SIGN_HOST,
    href: `https://${GDSTUDIO_SIGN_HOST}/`,
  }
  ;(sandbox.window as { location: unknown }).location = sandbox.location
  const run = new Function(`${vendorCode}\nreturn this.md5;`)
  const md5 = run.call(sandbox) as unknown
  if (typeof md5 !== 'function') {
    throw new Error('gdstudio md5 unavailable')
  }
  cachedMd5 = md5 as (input: string) => string
  return cachedMd5
}

/** 与播放器 functions.js 的 urlEncode 一致 */
export function gdstudioUrlEncode(value: unknown): string {
  if (value == null) return ''
  const text = typeof value === 'string' ? value : String(value)
  return encodeURIComponent(text)
    .replace(/\(/g, '%28')
    .replace(/\)/g, '%29')
    .replace(/\*/g, '%2A')
    .replace(/'/g, '%27')
    .replace(/!/g, '%21')
}

function padPlayerVersion(version: string): string {
  return version
    .split('.')
    .map((part) => (part.length === 1 ? `0${part}` : part))
    .join('')
}

export function gdstudioSignPayload(payload: string, serverTime: string): string {
  const msg = [
    String(serverTime).trim().slice(0, 9),
    GDSTUDIO_SIGN_HOST,
    padPlayerVersion(GDSTUDIO_PLAYER_VERSION),
    payload,
  ].join('|')
  const digest = loadGdstudioMd5()(msg)
  return String(digest).slice(-8).toUpperCase()
}

const FETCH_TIMEOUT_MS = 12_000

export async function fetchGdstudioServerTime(apiOrigin: string): Promise<string> {
  const ac = new AbortController()
  const timer = setTimeout(() => ac.abort(), FETCH_TIMEOUT_MS)
  try {
    const res = await fetch(`${apiOrigin}/time`, { signal: ac.signal, cache: 'no-store' })
    if (!res.ok) throw new Error(`gdstudio time ${res.status}`)
    const text = (await res.text()).trim()
    if (!/^\d+$/.test(text)) throw new Error('gdstudio time invalid')
    return text
  } finally {
    clearTimeout(timer)
  }
}
