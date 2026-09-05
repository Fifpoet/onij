import { GetUploadToken } from '@/api/file'
import { FileType } from '@/api/types/enums'
import { normalizeFileCdnUrl } from '@/util/qiniuCdn'

export { normalizeFileCdnUrl }

const DEFAULT_DOMAIN = 'http://cloud.onij.fun'
const DEFAULT_UPLOAD_URL = 'https://up-z0.qiniup.com'

async function getUploadCredentials(): Promise<{ token: string; uploadUrl: string; publicBase: string }> {
  const response = await GetUploadToken()
  const token = response.upload_token ?? response.uploadToken
  if (!token) {
    throw new Error(response.message || '获取上传凭证失败')
  }
  const publicBase = (response.domain || DEFAULT_DOMAIN).replace(/\/$/, '')
  const uploadUrl = (response.upload_url ?? response.uploadUrl ?? DEFAULT_UPLOAD_URL).replace(/\/$/, '')
  return {
    token,
    uploadUrl,
    publicBase,
  }
}

function generateUniqueFileName(originalName: string): string {
  const timestamp = Date.now()
  const randomStr = Math.random().toString(36).substring(2, 10)
  const extension = originalName.includes('.') ? originalName.split('.').pop() : ''
  const nameWithoutExt = originalName.includes('.')
    ? originalName.substring(0, originalName.lastIndexOf('.'))
    : originalName
  return `${nameWithoutExt}_${timestamp}_${randomStr}${extension ? '.' + extension : ''}`
}

function localFileFingerprint(file: File): string {
  return `local-${file.size}-${file.lastModified}-${file.name}`
}

/** 直传七牛（FormData），不经过 qiniu-js，避免 HTTP 下 crypto.subtle.digest 报错 */
function uploadViaFormData(
  file: File,
  token: string,
  key: string,
  uploadUrl: string,
  onProgress?: (percent: number) => void,
): Promise<{ key: string; hash: string }> {
  return new Promise((resolve, reject) => {
    const xhr = new XMLHttpRequest()
    const form = new FormData()
    form.append('file', file)
    form.append('token', token)
    form.append('key', key)

    xhr.upload.onprogress = (e) => {
      if (!e.lengthComputable) return
      onProgress?.(Math.min(100, Math.max(0, Math.round((e.loaded / e.total) * 100))))
    }

    xhr.onload = () => {
      if (xhr.status < 200 || xhr.status >= 300) {
        reject(new Error(xhr.responseText?.trim() || `七牛上传失败 (${xhr.status})`))
        return
      }
      try {
        const res = JSON.parse(xhr.responseText) as { key?: string; hash?: string }
        if (!res.key) {
          reject(new Error('七牛响应缺少 key'))
          return
        }
        resolve({
          key: res.key,
          hash: res.hash || localFileFingerprint(file),
        })
      } catch {
        reject(new Error('七牛响应解析失败'))
      }
    }

    xhr.onerror = () => reject(new Error('七牛上传网络错误'))
    xhr.open('POST', uploadUrl)
    xhr.send(form)
  })
}

const RESUME_THRESHOLD = 32 * 1024 * 1024
const BLOCK_SIZE = 4 * 1024 * 1024
const BLOCK_CONCURRENCY = 4

function urlSafeBase64(str: string): string {
  const bytes = new TextEncoder().encode(str)
  let binary = ''
  for (let i = 0; i < bytes.length; i++) binary += String.fromCharCode(bytes[i])
  return btoa(binary).replace(/\+/g, '-').replace(/\//g, '_')
}

function normalizeUpHost(host: string, fallback: string): string {
  const raw = (host || fallback).replace(/\/$/, '')
  if (/^https?:\/\//i.test(raw)) return raw
  return `https://${raw}`
}

type ResumeRet = { ctx?: string; host?: string; key?: string; hash?: string }

function xhrBinary(
  url: string,
  token: string,
  body: Blob,
  contentType: string,
  onBytes?: (loaded: number) => void,
): Promise<ResumeRet> {
  return new Promise((resolve, reject) => {
    const xhr = new XMLHttpRequest()
    xhr.open('POST', url)
    xhr.setRequestHeader('Authorization', `UpToken ${token}`)
    xhr.setRequestHeader('Content-Type', contentType)
    xhr.upload.onprogress = (e) => {
      if (e.lengthComputable) onBytes?.(e.loaded)
    }
    xhr.onload = () => {
      onBytes?.(body.size)
      if (xhr.status < 200 || xhr.status >= 300) {
        reject(new Error(xhr.responseText?.trim() || `七牛分片失败 (${xhr.status})`))
        return
      }
      try {
        resolve(JSON.parse(xhr.responseText) as ResumeRet)
      } catch {
        reject(new Error('七牛分片响应解析失败'))
      }
    }
    xhr.onerror = () => reject(new Error('七牛分片网络错误'))
    xhr.send(body)
  })
}

async function mapPool<T, R>(
  items: T[],
  limit: number,
  fn: (item: T, index: number) => Promise<R>,
): Promise<R[]> {
  const out = new Array<R>(items.length)
  let next = 0
  async function worker() {
    while (next < items.length) {
      const idx = next
      next += 1
      out[idx] = await fn(items[idx], idx)
    }
  }
  const n = Math.min(Math.max(1, limit), items.length)
  await Promise.all(Array.from({ length: n }, () => worker()))
  return out
}

/** ≥32MB 按 4MB 块并行 mkblk，避免 1MB 串行把带宽打成几 MB/s */
async function uploadViaResume(
  file: File,
  token: string,
  key: string,
  uploadUrl: string,
  onProgress?: (percent: number) => void,
): Promise<{ key: string; hash: string }> {
  const host = normalizeUpHost(uploadUrl, uploadUrl)
  const blocks: Array<{ offset: number; size: number }> = []
  for (let offset = 0; offset < file.size; offset += BLOCK_SIZE) {
    blocks.push({ offset, size: Math.min(BLOCK_SIZE, file.size - offset) })
  }
  const loaded = new Array(blocks.length).fill(0)
  const bump = () => {
    const sum = loaded.reduce((a, b) => a + b, 0)
    onProgress?.(Math.min(99, Math.max(0, Math.round((sum / file.size) * 100))))
  }

  const ctxs = await mapPool(blocks, BLOCK_CONCURRENCY, async (block, idx) => {
    const blob = file.slice(block.offset, block.offset + block.size)
    const ret = await xhrBinary(
      `${host}/mkblk/${block.size}`,
      token,
      blob,
      'application/octet-stream',
      (n) => {
        loaded[idx] = n
        bump()
      },
    )
    if (!ret.ctx) throw new Error('七牛分片缺少 ctx')
    return ret.ctx
  })

  const mkfileUrl = `${host}/mkfile/${file.size}/key/${urlSafeBase64(key)}`
  const done = await xhrBinary(mkfileUrl, token, new Blob([ctxs.join(',')]), 'text/plain')
  onProgress?.(100)
  return {
    key: done.key || key,
    hash: done.hash || localFileFingerprint(file),
  }
}

export type QiniuUploadProgress = (percent: number) => void

export async function uploadToQiniu(
  file: File,
  folderPath: string[] = [],
  onProgress?: QiniuUploadProgress,
): Promise<{ key: string; hash: string; url: string }> {
  const { token, uploadUrl, publicBase } = await getUploadCredentials()
  const uniqueFileName = generateUniqueFileName(file.name)
  const key =
    folderPath.length > 0 ? `${folderPath.join('/')}/${uniqueFileName}` : uniqueFileName

  const res =
    file.size >= RESUME_THRESHOLD
      ? await uploadViaResume(file, token, key, uploadUrl, onProgress)
      : await uploadViaFormData(file, token, key, uploadUrl, onProgress)
  return {
    key: res.key,
    hash: res.hash,
    url: normalizeFileCdnUrl(`${publicBase}/${res.key}`),
  }
}

/** 根据浏览器 File 推断 proto FileType */
export function getFileTypeFromFile(file: File): FileType {
  const mime = file.type.toLowerCase()
  const ext = file.name.split('.').pop()?.toLowerCase() || ''

  if (mime.startsWith('audio/') || ['mp3', 'wav', 'flac', 'aac', 'm4a', 'ogg', 'wma'].includes(ext)) {
    return FileType.FT_Audio
  }
  if (
    mime.startsWith('image/') ||
    ['png', 'jpg', 'jpeg', 'gif', 'bmp', 'webp', 'svg', 'ico', 'heic'].includes(ext)
  ) {
    return FileType.FT_Image
  }
  if (
    mime.startsWith('video/') ||
    ['mp4', 'mkv', 'avi', 'mov', 'webm', 'flv', 'wmv'].includes(ext)
  ) {
    return FileType.FT_Video
  }
  if (mime === 'application/pdf' || ext === 'pdf') {
    return FileType.FT_Pdf
  }
  if (
    ['doc', 'docx'].includes(ext) ||
    mime.includes('wordprocessingml') ||
    mime === 'application/msword'
  ) {
    return FileType.FT_Word
  }
  if (['xls', 'xlsx'].includes(ext) || mime.includes('spreadsheetml') || mime === 'application/vnd.ms-excel') {
    return FileType.FT_Excel
  }
  if (
    ['ppt', 'pptx'].includes(ext) ||
    mime.includes('presentationml') ||
    mime === 'application/vnd.ms-powerpoint'
  ) {
    return FileType.FT_PowerPoint
  }
  if (
    mime.startsWith('text/') ||
    ['txt', 'lrc', 'md', 'json', 'xml', 'csv'].includes(ext)
  ) {
    if (ext === 'csv' || mime === 'text/csv') return FileType.FT_CSV
    return FileType.FT_Text
  }
  if (['zip', 'rar', '7z', 'tar', 'gz', 'bz2'].includes(ext) || mime.includes('zip') || mime.includes('compressed')) {
    return FileType.FT_Archive
  }
  if (
    ['js', 'ts', 'jsx', 'tsx', 'vue', 'go', 'py', 'java', 'c', 'cpp', 'h', 'rs', 'sh'].includes(ext)
  ) {
    return FileType.FT_Code
  }
  if (['exe', 'msi', 'dmg', 'app', 'deb', 'rpm'].includes(ext) || mime === 'application/octet-stream') {
    if (['exe', 'msi', 'dmg', 'app'].includes(ext)) return FileType.FT_Executable
  }

  return FileType.FT_Unknown
}

export function batchUploadToQiniu(
  files: File[],
  folderPath: string[] = [],
): Promise<Array<{ key: string; hash: string; url: string; file: File }>> {
  return Promise.all(
    files.map((file) =>
      uploadToQiniu(file, folderPath).then((result) => ({
        ...result,
        file,
      })),
    ),
  )
}

export function resolveQiniuCdnUrl(url: string): string {
  const fixed = normalizeFileCdnUrl(url)
  if (!fixed || !import.meta.env.DEV) return fixed
  const origin = (import.meta.env.VITE_QINIU_CDN_ORIGIN as string | undefined)?.replace(/\/$/, '')
  if (!origin || /onij\.fun/i.test(origin)) return fixed
  return fixed.replace(/^https?:\/\/cloud\.onij\.fun(?=\/|$)/i, origin)
}

/** 经后端 /file/preview 拉取并触发本地下载 */
export async function downloadFileBlob(url: string, filename: string) {
  const resp = await fetch(url)
  if (!resp.ok) throw new Error(`download failed: ${resp.status}`)
  const blob = await resp.blob()
  const objectUrl = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = objectUrl
  a.download = filename
  a.style.display = 'none'
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(objectUrl)
}

/** 七牛私有链接触发下载（中转站等仍走签名链） */
export function downloadByPrivateUrl(url: string, filename?: string) {
  if (!url) return

  let href = resolveQiniuCdnUrl(url)
  const name = filename?.trim()
  if (name) {
    const sep = href.includes('?') ? '&' : '?'
    href = `${href}${sep}attname=${encodeURIComponent(name)}`
  }

  const a = document.createElement('a')
  a.href = href
  a.style.display = 'none'
  a.rel = 'noopener noreferrer'
  if (name) a.download = name
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
}
