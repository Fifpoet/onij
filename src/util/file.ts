// src/util/file.ts
import { FileType } from '@/api/types/file';
import { 
  FolderOutline, 
  MusicalNoteOutline, 
  DocumentTextOutline, 
  ImageOutline
} from '@vicons/ionicons5';

/**
 * 格式化时间戳为可读的日期时间字符串
 * @param timestamp 时间戳（秒）
 * @returns 格式化后的日期时间字符串
 */
export function formatFileTime(timestamp: number): string {
  const date = new Date(timestamp * 1000);
  
  // 检查日期是否有效
  if (isNaN(date.getTime())) {
    return '';
  }
  
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  });
}

/**
 * 根据文件类型获取对应的图标组件
 * @param fileType 文件类型
 * @returns 图标组件
 */
export function getFileIcon(fileType: FileType) {
  switch (fileType) {
    case FileType.FT_Folder:
      return FolderOutline;
    case FileType.FT_Mp3:
      return MusicalNoteOutline;
    case FileType.FT_Lyrics:
      return DocumentTextOutline;
    case FileType.FT_Png:
      return ImageOutline;
    default:
      return DocumentTextOutline;
  }
}

/**
 * 获取文件大小的显示文本（暂时写死，后续接口返回后可替换）
 * @param fileType 文件类型
 * @returns 文件大小字符串
 */
export function getFileSize(fileType: FileType): string {
  if (fileType === FileType.FT_Folder) {
    return '';
  }
  
  // 暂时写死的大小，后续接口返回后可以替换
  const sizeMap = {
    [FileType.FT_Mp3]: '3.2 MB',
    [FileType.FT_Lyrics]: '0.1 MB',
    [FileType.FT_Png]: '1.5 MB',
  };
  
  return sizeMap[fileType] || '0.1 MB';
}

/**
 * 获取文件类型的显示名称
 * @param fileType 文件类型
 * @returns 文件类型名称
 */
export function getFileTypeName(fileType: FileType): string {
  const typeMap = {
    [FileType.FT_Unknown]: '未知文件',
    [FileType.FT_Mp3]: '音频文件',
    [FileType.FT_Lyrics]: '歌词文件',
    [FileType.FT_Png]: '图片文件',
    [FileType.FT_Folder]: '文件夹',
  };
  
  return typeMap[fileType] || '未知文件';
}

/**
 * 检查是否为文件夹
 * @param fileType 文件类型
 * @returns 是否为文件夹
 */
export function isFolder(fileType: FileType): boolean {
  return fileType === FileType.FT_Folder;
}

/**
 * 检查是否为图片文件
 * @param fileType 文件类型
 * @returns 是否为图片文件
 */
export function isImage(fileType: FileType): boolean {
  return fileType === FileType.FT_Png;
}

/**
 * 检查是否为音频文件
 * @param fileType 文件类型
 * @returns 是否为音频文件
 */
export function isAudio(fileType: FileType): boolean {
  return fileType === FileType.FT_Mp3;
}

/**
 * 检查是否为文本文件
 * @param fileType 文件类型
 * @returns 是否为文本文件
 */
export function isText(fileType: FileType): boolean {
  return fileType === FileType.FT_Lyrics;
}

/**
 * 获取文件扩展名
 * @param filename 文件名
 * @returns 文件扩展名
 */
export function getFileExtension(filename: string): string {
  const lastDotIndex = filename.lastIndexOf('.');
  if (lastDotIndex === -1) {
    return '';
  }
  return filename.substring(lastDotIndex + 1).toLowerCase();
}

/**
 * 格式化文件大小（字节转换为可读格式）
 * @param bytes 字节数
 * @returns 格式化后的文件大小
 */
export function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 B';
  
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}