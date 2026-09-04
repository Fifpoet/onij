// src/api/file.ts
import { post } from '../util/http'
import { normalizeFileCdnUrl } from '@/util/qiniuCdn'
import type {
  DeleteFileReq,
  DeleteFileResp,
  DownloadFileReq,
  DownloadFileResp,
  GetFileListReq,
  GetFileListResp,
  GetUploadTokenResp,
  UploadFileReq,
  UploadFileResp,
} from '@/api/types'

export const UploadFile = async (params: UploadFileReq): Promise<UploadFileResp> => {
  return post('/file/upload', params) as Promise<UploadFileResp>
}

export const GetFileList = async (params: GetFileListReq): Promise<GetFileListResp> => {
  const resp = (await post('/file/list', params)) as GetFileListResp
  if (resp.files) {
    resp.files = resp.files.map((f) => (f?.url ? { ...f, url: normalizeFileCdnUrl(f.url) } : f))
  }
  return resp
}

export const DeleteFileById = async (params: DeleteFileReq): Promise<DeleteFileResp> => {
  return post('/file/delete', params) as Promise<DeleteFileResp>
}

export const DownloadFiles = async (params: DownloadFileReq): Promise<DownloadFileResp> => {
  const resp = (await post('/file/download', params)) as DownloadFileResp
  if (resp.urls) {
    resp.urls = resp.urls.map((u) => (u ? normalizeFileCdnUrl(u) : u))
  }
  return resp
}

export const GetUploadToken = async (): Promise<GetUploadTokenResp> => {
  return post('/file/upload_token', {}) as Promise<GetUploadTokenResp>
}
