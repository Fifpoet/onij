// src/api/file.ts
import { post } from '../util/http';
import { 
  UploadFileReq, 
  UploadFileResp, 
  GetFileListReq, 
  GetFileListResp,
  DownloadFileReq,
  DownloadFileResp
} from './types/file';

// 将文件转换为base64编码
const fileToBase64 = (file: File): Promise<string> => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.readAsDataURL(file)
    reader.onload = () => {
      const result = reader.result as string
      // 移除data:image/jpeg;base64,前缀，只保留base64编码部分
      const base64 = result.split(',')[1]
      resolve(base64)
    }
    reader.onerror = error => reject(error)
  })
}

export const UploadFile = async (params: UploadFileReq): Promise<UploadFileResp> => {
  // 将文件转换为base64编码
  const filesWithBase64 = await Promise.all(
    params.files.map(async (fileInfo) => ({
      filename: fileInfo.filename || '',
      file: await fileToBase64(fileInfo.file),
      origin_at: fileInfo.origin_at || Math.floor(Date.now() / 1000)
    }))
  )

  const requestData = {
    parent_id: params.parent_id,
    files: filesWithBase64
  }

  const response = await post<UploadFileResp>('/file/upload', requestData)
  return response.data
}

// 获取文件列表
export const GetFileList = async (params: GetFileListReq): Promise<GetFileListResp> => {
  return (await post<GetFileListResp>('/file/list', params)).data;
};

// 下载文件
export const DownloadFile = async (params: DownloadFileReq): Promise<DownloadFileResp> => {
  return (await post<DownloadFileResp>('/file/download', params)).data;
};