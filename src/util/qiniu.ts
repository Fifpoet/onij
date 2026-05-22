import * as qiniu from 'qiniu-js'
import { FileType } from '@/api/types/enums'
import { GetUploadToken } from '@/api/file'

// 七牛云配置
const QINIU_CONFIG = {
  bucket: 'onij',
  domain: 'http://cloud.onij.fun',
}

// 获取上传token
async function getUploadToken(): Promise<string> {
  const response = await GetUploadToken()
  const token = response.upload_token ?? response.uploadToken
  if (!token) {
    throw new Error(response.message || 'Failed to get upload token')
  }
  return token
}

// 计算文件hash
async function calculateFileHash(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = async (e) => {
      try {
        const arrayBuffer = e.target?.result as ArrayBuffer
        const hashBuffer = await crypto.subtle.digest('SHA-256', arrayBuffer)
        const hashArray = Array.from(new Uint8Array(hashBuffer))
        const hashHex = hashArray.map((b) => b.toString(16).padStart(2, '0')).join('')
        resolve(hashHex)
      } catch (error) {
        reject(error)
      }
    }
    reader.onerror = reject
    reader.readAsArrayBuffer(file)
  })
}

// 生成唯一的文件名
function generateUniqueFileName(originalName: string): string {
  const timestamp = Date.now()
  const randomStr = Math.random().toString(36).substring(2, 10)
  const extension = originalName.includes('.') ? originalName.split('.').pop() : ''
  const nameWithoutExt = originalName.includes('.')
    ? originalName.substring(0, originalName.lastIndexOf('.'))
    : originalName

  return `${nameWithoutExt}_${timestamp}_${randomStr}${extension ? '.' + extension : ''}`
}

export type QiniuUploadProgress = (percent: number) => void

// 上传文件到七牛云
export async function uploadToQiniu(
  file: File,
  folderPath: string[] = [],
  onProgress?: QiniuUploadProgress,
): Promise<{ key: string; hash: string; url: string }> {
  return new Promise(async (resolve, reject) => {
    try {
      const uniqueFileName = generateUniqueFileName(file.name)
      const key =
        folderPath.length > 0 ? `${folderPath.join('/')}/${uniqueFileName}` : uniqueFileName

      const token = await getUploadToken()

      const config = {
        useCdnDomain: false,
        region: qiniu.region.z0,
      }

      // @ts-ignore
      const observable = qiniu.upload(file, key, token, config)

      observable.subscribe({
        next: (res) => {
          const pct = Math.min(100, Math.max(0, res.total.percent ?? 0))
          onProgress?.(pct)
        },
        error: (err) => {
          reject(err)
        },
        complete: (res) => {
          calculateFileHash(file)
            .then((hash) => {
              resolve({
                key: res.key,
                hash,
                url: `${QINIU_CONFIG.domain}/${res.key}`,
              })
            })
            .catch(reject)
        },
      })
    } catch (error) {
      reject(error)
    }
  })
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
