import { FileType } from './enums'

export { FileType, FILE_TYPE_LABELS } from './enums'

export interface File {  id: number
  name: string
  format: FileType
  size: number
  parentId: number
  extra: string
  createdAt: number
  updatedAt: number
  url: string
}

/** 文件列表项（与 File 相同，遗留组件别名） */
export type FileDetail = File

// 请求和响应类型（字段与后端 JSON 一致）
export interface UploadFileReq {
  parent_id: number
  files: Array<{
    filename: string
    store_key: string
    hash: string
    format: FileType
    exif: string
    size: number
    origin_at?: number
  }>
  url?: string
}

export interface UploadFileResp {
  code: number
  message: string
  file_ids: number[]
}

export interface GetFileListReq {
  parentId?: number
  keyword?: string
  page: number
  limit: number
}

export interface GetFileListResp {
  code: number
  message: string
  files: File[]
  total: number
}

export interface DeleteFileReq {
  file_id: number
}

export interface DeleteFileResp {
  code: number
  message: string
}

export interface DownloadFileReq {
  file_ids: number[]
}

export interface DownloadFileResp {
  code: number
  message: string
  urls: string[]
}

export interface GetUploadTokenResp {
  code: number
  message: string
  upload_token: string
  uploadToken?: string
  domain: string
}