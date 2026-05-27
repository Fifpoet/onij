/**
 * Web Crypto 在 HTTP 非安全上下文（如 http://onij.fun）下不可用。
 * 业务代码请统一使用本模块，不要直接调用 crypto.randomUUID / crypto.subtle。
 */

export function isSecureContext(): boolean {
  return typeof globalThis !== 'undefined' && globalThis.isSecureContext === true
}

function fallbackUUID(): string {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
    const r = (Math.random() * 16) | 0
    const v = c === 'x' ? r : (r & 0x3) | 0x8
    return v.toString(16)
  })
}

export function randomUUID(): string {
  if (!isSecureContext()) return fallbackUUID()
  try {
    if (typeof crypto.randomUUID === 'function') {
      return crypto.randomUUID()
    }
  } catch {
    /* ignore */
  }
  return fallbackUUID()
}

export function isSubtleDigestAvailable(): boolean {
  try {
    return (
      isSecureContext() &&
      typeof crypto !== 'undefined' &&
      crypto.subtle != null &&
      typeof crypto.subtle.digest === 'function'
    )
  } catch {
    return false
  }
}

export async function sha256Hex(data: ArrayBuffer): Promise<string> {
  if (!isSubtleDigestAvailable()) {
    throw new Error('crypto.subtle 仅在 HTTPS 或 localhost 可用')
  }
  const hashBuffer = await crypto.subtle.digest('SHA-256', data)
  return Array.from(new Uint8Array(hashBuffer))
    .map((b) => b.toString(16).padStart(2, '0'))
    .join('')
}

/** 上传登记用文件指纹（可选）；HTTP 下请直接用七牛返回的 hash */
export async function hashFileForUpload(file: File): Promise<string> {
  const fallback = `local-${file.size}-${file.lastModified}-${file.name}`
  if (!isSubtleDigestAvailable()) return fallback

  return new Promise((resolve) => {
    const reader = new FileReader()
    reader.onload = async (e) => {
      try {
        const buf = e.target?.result
        if (!(buf instanceof ArrayBuffer)) {
          resolve(fallback)
          return
        }
        resolve(await sha256Hex(buf))
      } catch {
        resolve(fallback)
      }
    }
    reader.onerror = () => resolve(fallback)
    reader.readAsArrayBuffer(file)
  })
}
