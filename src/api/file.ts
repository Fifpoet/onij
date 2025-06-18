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
    reader.readAsArrayBuffer(file)
    reader.onload = () => {
      const arrayBuffer = reader.result as ArrayBuffer
      const bytes = new Uint8Array(arrayBuffer)
      let binary = ''
      for (let i = 0; i < bytes.byteLength; i++) {
        binary += String.fromCharCode(bytes[i])
      }
      const base64 = btoa(binary)
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
  );

  const requestData = {
    parent_id: params.parent_id,
    files: filesWithBase64
  };

  const response = await post<UploadFileResp>('/file/upload', requestData);
  return response.data;
};

// 获取文件列表
export const GetFileList = async (params: GetFileListReq): Promise<GetFileListResp> => {
  return (await post<GetFileListResp>('/file/list', params)).data;
};

// 下载文件
export const DownloadFile = async (params: DownloadFileReq): Promise<DownloadFileResp> => {
  return (await post<DownloadFileResp>('/file/download', params)).data;
};