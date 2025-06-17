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

// 上传文件
export const UploadFile = async (params: UploadFileReq): Promise<UploadFileResp> => {
  const formData = new FormData();
  formData.append('filename', params.filename);
  formData.append('parent_id', params.parent_id.toString());
  
  if (params.folders && params.folders.length > 0) {
    params.folders.forEach((folder, index) => {
      formData.append(`folders[${index}]`, folder);
    });
  }
  
  formData.append('file', params.file);
  
  // 添加原始文件时间
  if (params.origin_at !== undefined) {
    formData.append('origin_at', params.origin_at.toString());
  }
  
  return (await post<UploadFileResp>('/file/upload', formData)).data;
};

// 获取文件列表
export const GetFileList = async (params: GetFileListReq): Promise<GetFileListResp> => {
  return (await post<GetFileListResp>('/file/list', params)).data;
};

// 下载文件
export const DownloadFile = async (params: DownloadFileReq): Promise<DownloadFileResp> => {
  return (await post<DownloadFileResp>('/file/download', params)).data;
};