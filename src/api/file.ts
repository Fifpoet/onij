// src/api/file.ts
import { post } from '../util/http'
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
  return post('/file/list', params) as Promise<GetFileListResp>
}

export const DeleteFileById = async (params: DeleteFileReq): Promise<DeleteFileResp> => {
  return post('/file/delete', params) as Promise<DeleteFileResp>
}

export const DownloadFiles = async (params: DownloadFileReq): Promise<DownloadFileResp> => {
  return post('/file/download', params) as Promise<DownloadFileResp>
}

export const GetUploadToken = async (): Promise<GetUploadTokenResp> => {
  return post('/file/upload_token', {}) as Promise<GetUploadTokenResp>
}
