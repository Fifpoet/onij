import {
  defaultOriginFor,
  resolvePlaySourceOrigin,
  type PlaySource,
} from '@/store/playSources'
import { requestGdStudioAltSource, requestGdStudioPlayUrl } from './songUrl'

const PROBE_TIMEOUT_MS = 12_000
/** 只用来探测是否还吐直链，不依赖曲库业务 */
const PROBE_SONG_ID = '185811'

export function playSourceProbeUrl(source: PlaySource): string {
  const origin = resolvePlaySourceOrigin(source)
  if (source.kind === 'gdstudio') return `${origin}/api.php`
  return `${origin}/?server=netease&type=url&id=${PROBE_SONG_ID}`
}

export function playSourceDisplayOrigin(source: PlaySource): string {
  return resolvePlaySourceOrigin(source) || defaultOriginFor(source.kind)
}

export type PlaySourceProbeResult = {
  ok: boolean
  latencyMs: number
  body?: string
  error?: string
}

function usableUrl(raw: unknown): boolean {
  return typeof raw === 'string' && /^https?:\/\//i.test(raw)
}

function interpretMeting(status: number, text: string): { ok: boolean; error?: string } {
  if (status < 200 || status >= 400) return { ok: false, error: `HTTP ${status}` }
  const trimmed = text.trim()
  if (!trimmed) return { ok: false, error: '空响应' }
  if (/^https?:\/\//i.test(trimmed)) return { ok: true }
  try {
    const data = JSON.parse(trimmed) as { url?: string }
    if (usableUrl(data?.url)) return { ok: true }
    return { ok: false, error: '无播放地址' }
  } catch {
    if (/<html/i.test(trimmed)) return { ok: false, error: '返回了网页' }
    return { ok: false, error: '无法解析' }
  }
}

async function probeGdstudio(origin: string): Promise<{ ok: boolean; error?: string; body?: string }> {
  const neteaseUrl = await requestGdStudioPlayUrl(origin, Number(PROBE_SONG_ID), 320)
  if (neteaseUrl) return { ok: true, body: neteaseUrl.slice(0, 180) }
  const altUrl = await requestGdStudioAltSource(origin, { name: '海阔天空', artist: 'Beyond' }, 320)
  if (altUrl) return { ok: true, body: altUrl.slice(0, 180) }
  return { ok: false, error: '无播放地址' }
}

export async function probePlaySource(source: PlaySource): Promise<PlaySourceProbeResult> {
  const origin = resolvePlaySourceOrigin(source)
  const started = performance.now()
  const ac = new AbortController()
  const timer = setTimeout(() => ac.abort(), PROBE_TIMEOUT_MS)
  try {
    if (source.kind === 'gdstudio') {
      const judged = await probeGdstudio(origin)
      const latencyMs = Math.round(performance.now() - started)
      return judged.ok
        ? { ok: true, latencyMs, body: judged.body }
        : { ok: false, latencyMs, error: judged.error, body: judged.body }
    }
    const url = playSourceProbeUrl(source)
    const resp = await fetch(url, {
      method: 'GET',
      cache: 'no-store',
      redirect: 'manual',
      signal: ac.signal,
    })
    const latencyMs = Math.round(performance.now() - started)
    const loc = resp.headers.get('Location')
    if (loc && /^https?:\/\//i.test(loc)) {
      return { ok: true, latencyMs, body: loc.slice(0, 500) }
    }
    const text = (await resp.text()).slice(0, 500)
    const judged = interpretMeting(resp.status, text)
    if (!judged.ok) {
      return { ok: false, latencyMs, error: judged.error, body: text || undefined }
    }
    return { ok: true, latencyMs, body: text || undefined }
  } catch (err) {
    let msg = err instanceof Error ? err.message : String(err)
    if (/failed to fetch|networkerror|load failed/i.test(msg)) {
      msg = 'Failed to fetch（多半未放行或源已挂）'
    }
    return { ok: false, latencyMs: Math.round(performance.now() - started), error: msg }
  } finally {
    clearTimeout(timer)
  }
}
