/** 签名仍用 http://cloud.onij.fun；生产浏览器改走同源反代，避开七牛 *.qbox.me 证书 */
const QINIU_BIND_ORIGIN = 'http://cloud.onij.fun'
const QINIU_HTTPS_PROXY = 'https://onij.fun/cdn'

/** 修正误签主站域名；生产把 cloud.onij.fun 换成 https://onij.fun/cdn */
export function normalizeFileCdnUrl(url: string): string {
  if (!url) return url
  if (/^https?:\/\/onij\.fun\/cdn(\/|$)/i.test(url)) return url
  let next = url.replace(
    /^https?:\/\/onij\.fun(?=\/(?:tran|cloud|practice)\b)/i,
    QINIU_BIND_ORIGIN,
  )
  if (import.meta.env.PROD) {
    next = next.replace(/^https?:\/\/cloud\.onij\.fun(?=\/|$)/i, QINIU_HTTPS_PROXY)
  }
  return next
}
