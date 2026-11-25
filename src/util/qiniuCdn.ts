/** 签名与浏览都用七牛绑定的 HTTPS CDN 域，不要再绕主站 /cdn */
const QINIU_BIND_ORIGIN = 'https://cloud.onij.fun'

/** 修正误签主站 / 旧 /cdn 反代，统一成 https://cloud.onij.fun */
export function normalizeFileCdnUrl(url: string): string {
  if (!url) return url
  let next = url.replace(/^https?:\/\/onij\.fun\/cdn(?=\/|$)/i, QINIU_BIND_ORIGIN)
  next = next.replace(
    /^https?:\/\/onij\.fun(?=\/(?:tran|cloud|practice)\b)/i,
    QINIU_BIND_ORIGIN,
  )
  return next
}
