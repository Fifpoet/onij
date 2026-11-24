/**
 * GD Studio 播放器 (music.gdstudio.org) 的取链签名。
 * 他们把函数叫 crc32，实际是：time[:9]|hostname|version去点|payload → 自定义 MD5 末 8 位大写。
 * 标准 MD5 会对 401；必须跑播放器自带的那份哈希。
 * 监控页「同步签名」从源 origin 拉 js/player.js + js/crc32.min.js，校验后写入服务端；
 * 播放优先用服务端那份。打包 vendor 只作两边都没有时的离线兜底。
 */
import vendorCode from '@/vendor/gdstudio-crc32.min.js?raw'

/** 离线兜底，对齐打包时的 js/player.js mkPlayer.version */
export const GDSTUDIO_PLAYER_VERSION = '2026.11.22'

/** nginx / Vite 反代时 Host 也是这个，签名必须用同一 hostname */
export const GDSTUDIO_SIGN_HOST = 'music.gdstudio.org'

const FETCH_TIMEOUT_MS = 12_000
const FALLBACK_ORIGIN = '__vendor__'

export class GdstudioSignError extends Error {
  constructor(message = '签名失效，请到监控页点「同步签名」') {
    super(message)
    this.name = 'GdstudioSignError'
  }
}

export type GdstudioSignSource = 'remote' | 'server' | 'vendor'

type SignBundle = {
  version: string
  md5: (input: string) => string
  source: GdstudioSignSource
}

const cache = new Map<string, SignBundle>()
const inflight = new Map<string, Promise<SignBundle>>()
let serverBundle: SignBundle | null = null

function normalizeOrigin(origin: string): string {
  return origin.replace(/\/$/, '')
}

function pickMd5(sandbox: Record<string, unknown>, direct: unknown): ((input: string) => string) | null {
  const candidates = [
    direct,
    sandbox.md5,
    typeof window !== 'undefined' ? (window as unknown as { md5?: unknown }).md5 : undefined,
    (globalThis as { md5?: unknown }).md5,
  ]
  for (const fn of candidates) {
    if (typeof fn === 'function') return fn as (input: string) => string
  }
  return null
}

function loadMd5FromCode(code: string, version: string): (input: string) => string {
  const sandbox: Record<string, unknown> = {
    mkPlayer: { version },
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
  sandbox.globalThis = sandbox
  sandbox.location = {
    hostname: GDSTUDIO_SIGN_HOST,
    host: GDSTUDIO_SIGN_HOST,
    href: `https://${GDSTUDIO_SIGN_HOST}/`,
  }
  ;(sandbox.window as { location: unknown }).location = sandbox.location
  const prevGlobalMd5 = (globalThis as { md5?: unknown }).md5
  const run = new Function(`${code}
return (typeof this.md5 === 'function' && this.md5)
  || (typeof md5 === 'function' && md5)
  || (typeof window !== 'undefined' && typeof window.md5 === 'function' && window.md5)
  || (typeof globalThis !== 'undefined' && typeof globalThis.md5 === 'function' && globalThis.md5);`)
  try {
    const md5 = pickMd5(sandbox, run.call(sandbox))
    if (!md5) throw new Error('gdstudio md5 unavailable')
    return md5
  } finally {
    if ((globalThis as { md5?: unknown }).md5 !== prevGlobalMd5) {
      if (prevGlobalMd5 === undefined) delete (globalThis as { md5?: unknown }).md5
      else (globalThis as { md5?: unknown }).md5 = prevGlobalMd5
    }
  }
}

function fallbackBundle(): SignBundle {
  const cached = cache.get(FALLBACK_ORIGIN)
  if (cached) return cached
  const bundle: SignBundle = {
    version: GDSTUDIO_PLAYER_VERSION,
    md5: loadMd5FromCode(vendorCode, GDSTUDIO_PLAYER_VERSION),
    source: 'vendor',
  }
  cache.set(FALLBACK_ORIGIN, bundle)
  return bundle
}

function extractPlayerVersion(playerJs: string): string | null {
  const match = playerJs.match(/version\s*:\s*["'](\d+(?:\.\d+)*)["']/)
  return match?.[1] || null
}

async function fetchText(url: string): Promise<string> {
  const ac = new AbortController()
  const timer = setTimeout(() => ac.abort(), FETCH_TIMEOUT_MS)
  try {
    const res = await fetch(url, { signal: ac.signal, cache: 'no-store' })
    if (!res.ok) throw new Error(`gdstudio fetch ${res.status}`)
    return await res.text()
  } finally {
    clearTimeout(timer)
  }
}

export type RemoteGdstudioSign = {
  version: string
  crc32: string
}

async function fetchRemoteScripts(origin: string): Promise<RemoteGdstudioSign> {
  const base = normalizeOrigin(origin)
  const [playerJs, crc32Js] = await Promise.all([
    fetchText(`${base}/js/player.js`),
    fetchText(`${base}/js/crc32.min.js`),
  ])
  const version = extractPlayerVersion(playerJs)
  if (!version) throw new Error('gdstudio version missing')
  if (!crc32Js.trim()) throw new Error('gdstudio crc32 empty')
  loadMd5FromCode(crc32Js, version)
  return { version, crc32: crc32Js }
}

function bundleFromRemote(remote: RemoteGdstudioSign, source: GdstudioSignSource): SignBundle {
  return {
    version: remote.version,
    md5: loadMd5FromCode(remote.crc32, remote.version),
    source,
  }
}

/** 启动时灌入监控页同步过的服务端解码，全站共用 */
export function applyServerGdstudioSign(version: string, crc32: string) {
  const remote = { version: version.trim(), crc32 }
  if (!remote.version || !crc32.trim()) throw new Error('gdstudio sign empty')
  serverBundle = bundleFromRemote(remote, 'server')
  cache.clear()
  inflight.clear()
}

export function gdstudioSignVersion(): string {
  return serverBundle?.version || cache.values().next().value?.version || ''
}

export type EnsureGdstudioSignOptions = {
  force?: boolean
  remoteOnly?: boolean
}

export async function ensureGdstudioSign(
  origin: string,
  opts?: EnsureGdstudioSignOptions,
): Promise<SignBundle> {
  const key = normalizeOrigin(origin)
  if (!opts?.force && !opts?.remoteOnly) {
    const hit = cache.get(key)
    if (hit) return hit
    const pending = inflight.get(key)
    if (pending) return pending
    if (serverBundle) {
      cache.set(key, serverBundle)
      return serverBundle
    }
  }

  const task = (async () => {
    try {
      const remote = await fetchRemoteScripts(key)
      const bundle = bundleFromRemote(remote, 'remote')
      cache.set(key, bundle)
      return bundle
    } catch (err) {
      if (opts?.remoteOnly) {
        throw err instanceof Error ? err : new Error(String(err))
      }
      if (serverBundle) {
        cache.set(key, serverBundle)
        return serverBundle
      }
      return fallbackBundle()
    }
  })()

  inflight.set(key, task)
  try {
    return await task
  } finally {
    if (inflight.get(key) === task) inflight.delete(key)
  }
}

/** 从源站拉最新解码，失败则抛错（不退回 vendor） */
export async function fetchAndApplyRemoteGdstudioSign(origin: string): Promise<RemoteGdstudioSign> {
  const key = normalizeOrigin(origin)
  cache.delete(key)
  const remote = await fetchRemoteScripts(key)
  const bundle = bundleFromRemote(remote, 'remote')
  cache.set(key, bundle)
  serverBundle = { ...bundle, source: 'server' }
  return remote
}

export function resyncGdstudioSign(origin: string): Promise<RemoteGdstudioSign> {
  return fetchAndApplyRemoteGdstudioSign(origin)
}

export function clearGdstudioSign(origin: string) {
  cache.delete(normalizeOrigin(origin))
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

export function gdstudioSignPayload(origin: string, payload: string, serverTime: string): string {
  const bundle = cache.get(normalizeOrigin(origin)) ?? serverBundle ?? fallbackBundle()
  const msg = [
    String(serverTime).trim().slice(0, 9),
    GDSTUDIO_SIGN_HOST,
    padPlayerVersion(bundle.version),
    payload,
  ].join('|')
  return String(bundle.md5(msg)).slice(-8).toUpperCase()
}

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
