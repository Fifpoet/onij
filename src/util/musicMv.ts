/** MV 网页链接规范化与校验 */

export type MvUrlValidation = {
  ok: boolean
  url: string
  error?: string
}

function extractUrlString(raw: string): string {
  let s = raw.trim()
  if (!s) return ''

  const iframeMatch = s.match(/<iframe[^>]+src=["']([^"']+)["']/i)
  if (iframeMatch?.[1]) s = iframeMatch[1]

  const srcMatch = s.match(/src=["']([^"']+)["']/i)
  if (srcMatch?.[1]) s = srcMatch[1]

  if (s.startsWith('//')) s = `https:${s}`
  if (!/^https?:\/\//i.test(s)) s = `https://${s}`

  return s
}

/** 将 B 站 embed 转为视频页（若可识别） */
function preferWebPageUrl(url: URL): string {
  if (url.hostname === 'player.bilibili.com') {
    const bvid = url.searchParams.get('bvid')
    if (bvid) return `https://www.bilibili.com/video/${bvid}`
  }
  return url.toString()
}

export function normalizeMvUrl(raw: string): string {
  const s = extractUrlString(raw)
  if (!s) return ''

  try {
    const url = new URL(s)
    if (url.protocol !== 'http:' && url.protocol !== 'https:') return ''
    return preferWebPageUrl(url)
  } catch {
    return ''
  }
}

export function validateMvUrl(raw: string): MvUrlValidation {
  const trimmed = raw.trim()
  if (!trimmed) return { ok: true, url: '' }

  const normalized = normalizeMvUrl(trimmed)
  if (!normalized) {
    return { ok: false, url: '', error: '请输入有效的 http(s) 链接' }
  }

  try {
    const url = new URL(normalized)
    if (url.protocol !== 'http:' && url.protocol !== 'https:') {
      return { ok: false, url: '', error: '仅支持 http 或 https 链接' }
    }
    if (!url.hostname) {
      return { ok: false, url: '', error: 'URL 格式不正确' }
    }
    return { ok: true, url: normalized }
  } catch {
    return { ok: false, url: '', error: 'URL 格式不正确' }
  }
}

export function openMvUrl(url: string): void {
  const { ok, url: href } = validateMvUrl(url)
  if (!ok || !href) return
  window.open(href, '_blank', 'noopener,noreferrer')
}
