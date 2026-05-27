import { GetUploadToken } from '@/api/file'
import { FileType } from '@/api/types/enums'

const DEFAULT_DOMAIN = 'http://cloud.onij.fun'

async function getUploadCredentials(): Promise<{ token: string; uploadUrl: string; publicBase: string }> {
  const response = await GetUploadToken()
  const token = response.upload_token ?? response.uploadToken
  if (!token) {
    throw new Error(response.message || '获取上传凭证失败')
  }
  const publicBase = (response.domain || DEFAULT_DOMAIN).replace(/\/$/, '')
  return {
    token,
    uploadUrl: `${publicBase}/upload`,
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

  const res = await uploadViaFormData(file, token, key, uploadUrl, onProgress)
  return {
    key: res.key,
    hash: res.hash,
    url: `${publicBase}/${res.key}`,
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

/** 与 FileList.downloadFile 相同：iframe 打开私有下载链 */
export function downloadByPrivateUrl(url: string) {
  if (!url) return
  const iframe = document.createElement('iframe')
  iframe.style.display = 'none'
  iframe.src = url
  document.body.appendChild(iframe)
  setTimeout(() => document.body.removeChild(iframe), 1000)
}
