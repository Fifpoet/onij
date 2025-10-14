import { FileType } from './enums'

export interface File {
  id: number
  name: string
  format: FileType
  size: number
  parentId: number
  extra: string
  createdAt: number
  updatedAt: number
  url: string
}

// 请求和响应类型
export interface UploadFileReq {
  parentId: number
  files: Array<{
    filename: string
    storeKey: string
    hash: string
    format: FileType
    exif: string
    size: number
  }>
  url?: string
}

export interface UploadFileResp {
  code: number
  message: string
  fileIds: number[]
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
  fileId: number
}

export interface DeleteFileResp {
  code: number
  message: string
}

export interface GetUploadTokenResp {
  code: number
  message: string
  uploadToken: string
  domain: string
}