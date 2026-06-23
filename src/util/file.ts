// src/util/file.ts
import { FileType } from '@/api/types/enums'
import { normalizeFileCdnUrl, resolveQiniuCdnUrl } from '@/util/qiniu'
import {
  DocumentTextOutline,
  FolderOutline,
  ImageOutline,
  MusicalNoteOutline,
  VideocamOutline,
} from '@vicons/ionicons5'

export function formatFileTime(timestamp: number): string {
  const date = new Date(timestamp * 1000)
  if (isNaN(date.getTime())) return ''
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}

export function showFileIcon(fileType: FileType): boolean {
  return fileType !== FileType.FT_Image
}

export function getFileIcon(fileType: FileType) {
  switch (fileType) {
    case FileType.FT_Folder:
      return FolderOutline
    case FileType.FT_Audio:
      return MusicalNoteOutline
    case FileType.FT_Image:
      return ImageOutline
    case FileType.FT_Video:
      return VideocamOutline
    case FileType.FT_Text:
    case FileType.FT_CSV:
    case FileType.FT_Code:
      return DocumentTextOutline
    default:
      return DocumentTextOutline
  }
}

export function getFileSize(fileType: FileType): string {
  if (fileType === FileType.FT_Folder) return ''
  const sizeMap: Partial<Record<FileType, string>> = {
    [FileType.FT_Audio]: '3.2 MB',
    [FileType.FT_Text]: '0.1 MB',
    [FileType.FT_Image]: '1.5 MB',
    [FileType.FT_Video]: '50 MB',
  }
  return sizeMap[fileType] || '0.1 MB'
}

export function getFileTypeName(fileType: FileType): string {
  const typeMap: Partial<Record<FileType, string>> = {
    [FileType.FT_Unknown]: '未知文件',
    [FileType.FT_Audio]: '音频文件',
    [FileType.FT_Text]: '文本文件',
    [FileType.FT_Image]: '图片文件',
    [FileType.FT_Video]: '视频文件',
    [FileType.FT_Folder]: '文件夹',
    [FileType.FT_Pdf]: 'PDF',
    [FileType.FT_Archive]: '压缩包',
  }
  return typeMap[fileType] || '未知文件'
}

export function isFolder(fileType: FileType): boolean {
  return fileType === FileType.FT_Folder
}

export function isImage(fileType: FileType): boolean {
  return fileType === FileType.FT_Image
}

export function isAudio(fileType: FileType): boolean {
  return fileType === FileType.FT_Audio
}

export function isText(fileType: FileType): boolean {
  return fileType === FileType.FT_Text || fileType === FileType.FT_CSV
}

/** 本地 dev：同源 /file/preview?proxy=1 经 Vite 代理 + 七牛 SDK 回源，不依赖 cloud.onij.fun DNS */
export function fileImageSrc(file: { id: number; url: string }): string {
  if (!file.url || !file.id) return ''
  if (import.meta.env.DEV) {
    return `/file/preview?file_id=${file.id}&proxy=1`
  }
  return resolveQiniuCdnUrl(file.url)
}

export function fileDownloadUrl(file: { id: number; name?: string; url: string }): string {
  if (import.meta.env.DEV && file.id) {
    const name = file.name ? `&filename=${encodeURIComponent(file.name)}` : ''
    return `/file/preview?file_id=${file.id}&download=1&proxy=1${name}`
  }
  return normalizeFileCdnUrl(file.url)
}

export function getFileExtension(filename: string): string {
  const lastDotIndex = filename.lastIndexOf('.')
  if (lastDotIndex === -1) return ''
  return filename.substring(lastDotIndex + 1).toLowerCase()
}

export function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`
}

/** 从 file.extra JSON 解析 EXIF 原始拍摄时间（Unix 秒） */
export function parseFileOriginAt(extra?: string): number {
  if (!extra) return 0
  try {
    const data = JSON.parse(extra) as { origin_at?: number }
    return data.origin_at && data.origin_at > 0 ? data.origin_at : 0
  } catch {
    return 0
  }
}

/** 展示时间：优先 EXIF 拍摄时间，否则用 created_at */
export function formatFileDisplayTime(file: { extra?: string; created_at?: number }): string {
  const originAt = parseFileOriginAt(file.extra)
  const ts = originAt > 0 ? originAt : (file.created_at ?? 0)
  return formatFileTime(ts)
}

export function getFileTypeFromString(fileType: string): FileType {
  const lowerType = fileType.toLowerCase()
  if (lowerType.startsWith('audio/') || lowerType === 'mp3') return FileType.FT_Audio
  if (lowerType.startsWith('image/') || ['png', 'jpg', 'jpeg'].includes(lowerType)) {
    return FileType.FT_Image
  }
  if (lowerType.startsWith('text/') || ['txt', 'lrc'].includes(lowerType)) return FileType.FT_Text
  const extension = getFileExtension(fileType)
  switch (extension) {
    case 'mp3':
    case 'wav':
    case 'flac':
      return FileType.FT_Audio
    case 'png':
    case 'jpg':
    case 'jpeg':
    case 'gif':
    case 'bmp':
      return FileType.FT_Image
    case 'txt':
    case 'lrc':
      return FileType.FT_Text
    case 'csv':
      return FileType.FT_CSV
    default:
      return FileType.FT_Unknown
  }
}
