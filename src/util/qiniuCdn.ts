/** 签名与浏览器访问都用 HTTPS 自定义 CDN 域；旧 /cdn 反代链改回直连 */
const QINIU_BIND_ORIGIN = 'https://cloud.onij.fun'

/** 修正误签主站域名；http / 旧反代一律改成 https://cloud.onij.fun */
export function normalizeFileCdnUrl(url: string): string {
  if (!url) return url
  let next = url.replace(
    /^https?:\/\/onij\.fun\/cdn(?=\/|$)/i,
    QINIU_BIND_ORIGIN,
  )
  next = next.replace(
    /^https?:\/\/onij\.fun(?=\/(?:tran|cloud|practice)\b)/i,
    QINIU_BIND_ORIGIN,
  )
  next = next.replace(/^https?:\/\/cloud\.onij\.fun(?=\/|$)/i, QINIU_BIND_ORIGIN)
  return next
}
