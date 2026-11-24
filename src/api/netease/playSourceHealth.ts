import {
  defaultOriginFor,
  resolvePlaySourceOrigin,
  type PlaySource,
} from '@/store/playSources'
import { requestGdApiPlayUrl, requestGdStudioAltSource, requestGdStudioPlayUrl } from './songUrl'

/** GD 公开 API 探测用能吐链的网易 id，不再用空曲 185811 */
const GDAPI_PROBE_SONG_ID = 1357375695

export function playSourceProbeUrl(source: PlaySource): string {
  const origin = resolvePlaySourceOrigin(source)
  if (source.kind === 'gdapi') {
    return `${origin}/api.php?types=url&source=netease&id=${GDAPI_PROBE_SONG_ID}&br=320`
  }
  return `${origin}/api.php`
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

async function probeGdapi(origin: string): Promise<{ ok: boolean; error?: string; body?: string }> {
  const url = await requestGdApiPlayUrl(origin, GDAPI_PROBE_SONG_ID, 320)
  if (url) return { ok: true, body: url.slice(0, 180) }
  return { ok: false, error: '无播放地址' }
}

async function probeGdstudio(origin: string): Promise<{ ok: boolean; error?: string; body?: string }> {
  const neteaseUrl = await requestGdStudioPlayUrl(origin, GDAPI_PROBE_SONG_ID, 320)
  if (neteaseUrl) return { ok: true, body: neteaseUrl.slice(0, 180) }
  const altUrl = await requestGdStudioAltSource(origin, { name: '海阔天空', artist: 'Beyond' }, 320)
  if (altUrl) return { ok: true, body: altUrl.slice(0, 180) }
  return { ok: false, error: '无播放地址' }
}

export async function probePlaySource(source: PlaySource): Promise<PlaySourceProbeResult> {
  const origin = resolvePlaySourceOrigin(source)
  const started = performance.now()
  try {
    const judged = source.kind === 'gdapi' ? await probeGdapi(origin) : await probeGdstudio(origin)
    const latencyMs = Math.round(performance.now() - started)
    return judged.ok
      ? { ok: true, latencyMs, body: judged.body }
      : { ok: false, latencyMs, error: judged.error, body: judged.body }
  } catch (err) {
    let msg = err instanceof Error ? err.message : String(err)
    if (/failed to fetch|networkerror|load failed/i.test(msg)) {
      msg = 'Failed to fetch（多半未放行或源已挂）'
    }
    return { ok: false, latencyMs: Math.round(performance.now() - started), error: msg }
  }
}
